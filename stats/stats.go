package stats

import (
	"fmt"
	"time"

	"github.com/divan/simulation/propagation"
)

// Stats represents stats data for given simulation log.
type Stats struct {
	NodeHits            map[int]int
	NodeCoverage        Coverage
	LinkCoverage        Coverage
	NodeHistogram       *Histogram
	LinkHistogram       *Histogram
	TimeToNodeHistogram *Histogram
	Time                time.Duration
}

// PrintVerbose prints detailed terminal-friendly stats to
// the console.
func (s *Stats) PrintVerbose() {
	fmt.Println("Stats:")
	fmt.Println("Time elapsed:", s.Time)
	fmt.Println("Nodes coverage:", s.NodeCoverage)
	fmt.Println("Links coverage:", s.LinkCoverage)
	fmt.Println("Nodes histogram:", s.NodeHistogram)
	fmt.Println("Links histogram:", s.LinkHistogram)
	fmt.Println("TimeToNode histogram:", s.TimeToNodeHistogram)
}

// Analyze analyzes given propagation log and returns filled Stats object.
func Analyze(plog *propagation.Log, nodeCount, linkCount int) *Stats {
	t := analyzeTiming(plog)
	nodeHits, nodeHistogram := analyzeNodeHits(plog)
	nodeCoverage := analyzeNodeCoverage(nodeHits, nodeCount)
	linkCoverage, linkHistogram := analyzeLinkCoverage(plog, linkCount)
	timeToNodeHistogram := analyzeTimeToNode(plog)

	return &Stats{
		NodeHits:            nodeHits,
		NodeCoverage:        nodeCoverage,
		LinkCoverage:        linkCoverage,
		NodeHistogram:       nodeHistogram,
		LinkHistogram:       linkHistogram,
		TimeToNodeHistogram: timeToNodeHistogram,
		Time:                t,
	}
}

func analyzeNodeHits(plog *propagation.Log) (map[int]int, *Histogram) {
	nodeHits := make(map[int]int)

	x := make([]float64, 0, len(plog.Timestamps))
	for _, nodes := range plog.Nodes {
		for _, j := range nodes {
			nodeHits[j]++
		}
		count := len(nodes)
		x = append(x, float64(count))
	}

	return nodeHits, NewHistogram(x, 20)
}

func analyzeNodeCoverage(nodeHits map[int]int, total int) Coverage {
	actual := len(nodeHits)
	return NewCoverage(actual, total)

}

func analyzeLinkCoverage(plog *propagation.Log, total int) (Coverage, *Histogram) {
	linkHits := make(map[int]struct{})

	x := make([]float64, 0, len(plog.Timestamps))
	for _, links := range plog.Links {
		for _, j := range links {
			linkHits[j] = struct{}{}
		}

		count := len(links)
		x = append(x, float64(count))
	}

	actual := len(linkHits)
	return NewCoverage(actual, total), NewHistogram(x, 20)
}

// analyzeTiming returns the amount of time the simulation took.
func analyzeTiming(plog *propagation.Log) time.Duration {
	// log contains timestamps in milliseconds, so the
	// max value will be our number
	var max int
	for _, ts := range plog.Timestamps {
		if ts > max {
			max = ts
		}
	}
	return time.Duration(max) * time.Millisecond
}

func analyzeTimeToNode(plog *propagation.Log) *Histogram {
	var hits = make(map[int]int)
	for i, ts := range plog.Timestamps {
		nodes := plog.Nodes[i]
		for _, j := range nodes {
			if _, ok := hits[j]; !ok {
				hits[j] = ts
			}
		}
	}

	x := make([]float64, 0, len(plog.Nodes))
	for _, ts := range hits {
		x = append(x, float64(ts))
	}
	return NewHistogram(x, 20)
}
