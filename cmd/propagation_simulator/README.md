# Propagation simulator
---

This simulator implements command for running different simulation implementations. Currently supported:
 - whisperv6
 - naive gossip propagation

# Installation

```
go get github.com/divan/simulation/cmd/propagation_simulator
```

# Usage

Just run:
```
propagation_simulator
```

This tool is looking for the `network.json` file as an input. You may override this name with `-i filename.json` command line flag. This should be valid JSON file with graph structure described here (link TBD). See examples/ directory.

```
propagation_simulator -i graph.json
```

Output statistics will be printed to the stdout, and final propagation data will be writtein into `propagation.json` file. (TODO: describe file format and further steps)

See `propagation_simulator --help` for more options.
