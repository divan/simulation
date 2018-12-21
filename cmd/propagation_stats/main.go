package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/divan/graphx/formats"
	"github.com/divan/simulation/propagation"
	"github.com/divan/simulation/stats"
)

func main() {
	var (
		network  = flag.String("n", "network.json", "Input filename for network graph data")
		plogFile = flag.String("p", "propagation.json", "Input filename for propagation log data")
	)
	flag.Parse()

	data, err := formats.FromD3JSON(*network)
	if err != nil {
		log.Fatal("Opening network file failed: ", err)
	}
	log.Printf("Loaded network graph from %s file", *network)

	fd, err := os.Open(*plogFile)
	if err != nil {
		log.Fatal("Opening propagation file failed: ", err)
	}
	defer fd.Close()
	log.Printf("Loaded propagation log from %s file", *plogFile)

	plog := &propagation.Log{}
	err = json.NewDecoder(fd).Decode(&plog)
	if err != nil {
		log.Fatalf("Parsing propagation log failed: %v", err)
	}

	ss := stats.Analyze(plog, data.NumNodes(), data.NumLinks())
	ss.PrintVerbose()
}
