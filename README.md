# golang-discovery-crdt

Very simple app to play with Golang and Locust performance testing for the first time.

Each app consist of rest service which holds in-memory CRDT of type Last-Write-Wins Element Set.

Nodes can connect with each other and form a cluster to share information and stay eventually consistent while updates to LWWES are being posted.

Locust scripts checks both update operations performance and cluster consistency convergence.

## Running app:
* build app:
`export GOPATH=$PWD && go build org/miejski/crdt-app`

* run first node:
`./crdt-app --port=8080`
 
* run another node to join the first one:
`./crdt-app --port=8081 --join=http://localhost:8080`

* post some element set:
`curl -X POST http://localhost:8080/status/update -d '{"Value":{"Value": 2}, "Operation": "ADD"}'`

* check if rest of the cluster received the value:
    `curl http://localhost:8081/status` or `curl http://localhost:8081/status/readable`

## running tests:


* Test dependencies:
    
    `go get github.com/stretchr/testify`

* Run all unit tests:

    `go test -v -race org/miejski/... -short `
    
* Run all integration tests:   
       
    `go test -v -race org/miejski/... -run Integration`

* Run performance tests:
    * Create virtualenv: `virtualenv  .venv `
    * `source .venv/bin/activate`
    * `pip install -r requirements.txt`
    * Run any number of crdt-app instances to create a cluster
    * `python load_balancer/prepare_LB.py single_node_host` (default to http://localhost:8080)
    * `./load_balancer/start_nginx.sh`
    * `locust -f main_load_test.py --host=http://localhost:9999 -c 300 -r 100 -n 500 --no-web`
        