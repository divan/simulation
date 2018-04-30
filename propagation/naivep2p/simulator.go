package naivep2p

import (
	"sync"
	"time"

	"github.com/divan/graph-experiments/graph"
	"github.com/status-im/simulator/propagation"
)

// Simulator is responsible for running propagation simulation.
type Simulator struct {
	data            *graph.Graph
	delay           time.Duration
	peers           map[int][]int
	nodesCh         []chan Message
	reportCh        chan LogEntry
	peersToSendTo   int // number of peers to propagate message
	wg              *sync.WaitGroup
	simulationStart time.Time
}

// Message represents the message propagated in the simulation.
type Message struct {
	Content string
	TTL     int
}

// NewSimulator initializes new simulator for the given graph data.
func NewSimulator(data *graph.Graph, N int, delay time.Duration) *Simulator {
	nodeCount := len(data.Nodes())
	sim := &Simulator{
		data:          data,
		delay:         delay,
		peers:         PrecalculatePeers(data),
		peersToSendTo: N,
		reportCh:      make(chan LogEntry),
		nodesCh:       make([]chan Message, nodeCount), // one channel per node
		wg:            new(sync.WaitGroup),
	}
	sim.wg.Add(nodeCount)
	for i := 0; i < nodeCount; i++ {
		ch := sim.startNode(i)
		sim.nodesCh[i] = ch // this channel will be used to talk to node by index
	}
	return sim
}

// Stop stops simulator and frees all resources if any. Implements propagation.PropagationSimulator.
func (s *Simulator) Stop() error {
	return nil
}

// SendMessage sends single message and tracks propagation. Implements propagation.PropagationSimulator.
func (s *Simulator) SendMessage(startNodeIdx, ttl int) *propagation.Log {
	message := Message{
		Content: "dummy",
		TTL:     ttl,
	}
	s.simulationStart = time.Now()
	s.propagateMessage(startNodeIdx, message)

	done := make(chan bool)
	go func() {
		s.wg.Wait()
		done <- true
	}()

	var ret []*LogEntry
	for {
		select {
		case val := <-s.reportCh:
			ret = append(ret, &val)
		case <-done:
			return s.logEntries2PropagationLog(ret)
		}
	}
}

func (s *Simulator) startNode(i int) chan Message {
	ch := make(chan Message)
	go s.runNode(i, ch)
	return ch
}

// runNode does actual node processing part
func (s *Simulator) runNode(i int, ch chan Message) {
	defer s.wg.Done()
	t := time.NewTimer(10 * time.Second)

	cache := make(map[string]bool)
	for {
		select {
		case message := <-ch:
			if cache[message.Content] {
				continue
			}
			cache[message.Content] = true
			message.TTL--
			if message.TTL == 0 {
				return
			}
			s.propagateMessage(i, message)
		case <-t.C:
			return
		}
	}
}

// propagateMessage simulates message sending from node to its peers.
func (s *Simulator) propagateMessage(from int, message Message) {
	time.Sleep(s.delay)
	peers := s.peers[from]
	for i := range peers {
		go s.sendMessage(from, peers[i], message)
	}
}

// sendMessage simulates message sending for given from and to indexes.
func (s *Simulator) sendMessage(from, to int, message Message) {
	s.nodesCh[to] <- message
	s.reportCh <- NewLogEntry(s.simulationStart, from, to)
}