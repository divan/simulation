package main

import (
	"flag"
	"log"
	"os"

	"github.com/divan/graphx/formats"
	"github.com/divan/simulation/stats"
	gethlog "github.com/ethereum/go-ethereum/log"
)

func main() {
	var (
		input        = flag.String("i", "network.json", "Input filename for pregenerated data to be used with simulation")
		output       = flag.String("o", "propagation.json", "Output filename for p2p sending data")
		gethlogLevel = flag.String("loglevel", "crit", "Geth log level for whisper simulator (crti, error, warn, info, debug, trace)")
		ttl          = flag.Int("ttl", 10, "TTL for generated messages")
		size         = flag.Int("msgSize", 400, "Payload size for generated messages")
		algorithm    = flag.String("algorithm", "whisperv6", "Propagation algorithm to use (whisperv6, gossip)")
	)
	flag.Parse()

	setGethLogLevel(*gethlogLevel)

	data, err := formats.FromD3JSON(*input)
	if err != nil {
		log.Fatal("Opening input file failed: ", err)
	}
	log.Printf("Loaded network graph from %s file", *input)

	algo := "whisperv6"
	if *algorithm == "gossip" {
		algo = "gossip"
	} // TODO: add proper validation for algorithm
	log.Printf("Using %s propagation algorithm", algo)

	sim := NewSimulation(algo, data)
	log.Printf("Starting message sending simulation for graph with %d nodes...", len(data.Nodes()))
	sim.Start(*ttl, *size)
	defer sim.Stop()
	sim.WriteOutputToFile(*output)

	// stats
	ss := stats.Analyze(sim.plog, data.NumNodes(), data.NumLinks())
	ss.PrintVerbose()

	log.Printf("Written propagation data into %s", *output)
}

func setGethLogLevel(level string) {
	lvl, err := gethlog.LvlFromString(level)
	if err != nil {
		lvl = gethlog.LvlCrit
	}
	gethlog.Root().SetHandler(gethlog.LvlFilterHandler(lvl, gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true))))
}
