## Prepare for running tests - setup load balancer
* Run any number of crdt-app nodes forming a cluster
* `python load_balancer/prepare_LB.py single_node_host` with default to 'http://localhost:8080' 
* `./load_balancer/start_nginx.sh`


## Running tests:
`locust -f main_load_test.py --host=http://localhost:9999 -c 300 -r 100 -n 500 --no-web` 

