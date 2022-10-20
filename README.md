# Peer-to-Peer Distributed Data Store using the Kademlia protocol

This is the lab work for the LTU course "D7024E - Mobile and Distributed Computing Systems"

## Usage

### Start demo network (script)

The script [`restart.sh`](./restart.sh) will automatically remove all artifacts from previous network. Then, build the required docker image and start the network. To attach a terminal connection to the network, i.e. a container running ubuntu connected to the Kademlia network, add the `-a` option.

```sh
./restart.sh -a
```

Note: The script must be executed in the project root.

### Start demo network (manually)

Use Docker compose to start up a network of 50 nodes. They will automatically connect and start communicating with each other. First, build the docker image

```sh
docker build -t docker-go .
```

Then, use the following to start the network.

```sh
docker compose up -d
```

The network will start in detached mode and to connect to a Kademlia node, follow the instructions in the next section [Connect to a container](#connect-to-an-active-node). To stop the network use

```sh
docker compose down
```

### Connect to an active node

Use

```sh
docker attach kademlia-node-{number}
```

to attach a running Kademlia node to the current terminal. Replace `{number}` with a numerical value, e.g. `5` to connect to node number five.

## Testing

Use [`test.sh`](./test.sh) to run all unit tests in the project. Once all tests have passed, a coverage report `./coverage.html` will be generated alongside with a `coverage.out` artifact.
