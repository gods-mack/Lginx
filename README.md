# Lginx
its a in-house simple Load-Balancer, ReverseProxy and more.

## Done
 - Round Robin Algorithm for balancing reqs
 - Reverse Proxy for multiple backend servers
 - Auto HealthCheck (backend server's)
 - Mutex (Down/Up for backends)


## Pending
 - config file support
 - http in-mem local cache
 - more load-balancing algorithms



## How to run
 - git clone repo
 - `go run main.go`
 - add configs in `config.json`
    - "backend_hosts"  -> add all your backend server IP
    - "default_proxy_server" -> "add your desired proxy server,
        by default it is `8081`"
    
e.g -> `config.json`
```
{
	"default_lginx_port":"8081",
	"proxy_pass":"",
	"proxy_server_name":"lginx",
	"backend_hosts":[
		"http://127.0.0.1:8001",
		"http://127.0.0.1:8002"
	] 
}
     
