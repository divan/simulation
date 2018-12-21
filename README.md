# P2P messaging simulation toolkit
---
This repository holds different simulators for exploring and researching p2p networks and messaging.

Original intent of these simulators is to provide stats and resulting traces/logs for further analysis and visualization.

## Design
```

+------------------+   +----------------+   +-------------+   +------------------+                 
| Choose           |   |  Run nodes in  |   |             |   |                  |                 
| network topology |----  simulated     |   |             |   |                  |                 
+------------------+   |  environment   |   | Propagate   |   | Collect network  |                 
                       |   - in-memory  |---- message(s)  |---- events &         |                 
+------------------+   |   - exec       |   |             |   | generate stats   |                 
|  Choose          |----   - docker     |   |             |   |                  |                 
|  Simulator       |   |                |   |             |   |                  |                 
+------------------+   +----------------+   +-------------+   +------------------+                 
```

### Simulators support

| Simulator  | Description | State |
|---|---|---|
| **WhisperV6** | Master branch if go-ethereum Whisper implementation  | Done |
| **Gossip**  | Naive gossip p2p propagation  | Done |
| PSS | Swarm's PSS messaging | TBD |

### Network environments support

| Node type  | Description | State |
|---|---|---|
| **In-Memory** | Done | Single node in-memory network  | Done |
| Exec  | Single node native binary network with localhost connection | TBD |
| Docker | Docker-based network | TBD |

## Usage
As a backend for the visualization frontend:

```
go get github.com/divan/simulation/cmd/propagation_server

propagation_server
```

As a commandline tool:

```
go get github.com/divan/simulation/cmd/propagation_simulator
// copy network.json to current directory
propagation_simulator --help
```

## License
MIT
