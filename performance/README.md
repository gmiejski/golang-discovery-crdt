## Prepare for running tests - setup load balancer
* `python load_balancer/prepare_LB.py` 
* `./load_balancer/start_nginx.sh`


## Running tests:
`locust -f main_load_test.py --host=http://localhost:9999 -c 300 -r 100 -n 500 --no-web` 

