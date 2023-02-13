
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	//"html"
	//"math/rand"
	"time"
	//"strings"
	//"os"
	//"io/ioutil"
	//"encoding/json"
	"sync"
	"sync/atomic"
	"net/http/httputil"
	"net"
)

/*

httputil.ReverseProxy
ReverseProxy is an HTTP Handler that takes an incoming request 
and sends it to another server, proxying the rebponse back to the client.
*/



// func get_backend_IP() string {
// 	rand.Seed(time.Now().Unix())
// 	hosts := []string{
// 		"127.0.0.1:8001",
// 		"192.168.1.27:8001",
// 	}
// 	n := rand.Int() % len(hosts)
// 	return "http://" + hosts[n]
// }

// func str_to_json(body string)  {
// 	...
// }

// func proxy_handler(w http.RebponseWriter, r *http.Request) {

// 	node_ip := get_backend_IP()
// 	fmt.Printf("Htting %s\n",node_ip)
	
// 	endpoint := node_ip + r.URL.Path
// 	if r.Method == http.MethodGet {
// 		req, err := http.NewRequest(r.Method, endpoint, nil )
// 		if err != nil {
// 			fmt.Printf("client: could not create request: %s\n", err)
// 		}
	
// 	}
// 	if r.Method == http.MethodPost {
// 		req, err:= http.NewRequest(r.Method, endpoint, req_body )
// 		if err != nil {
// 			fmt.Printf("client: could not create request: %s\n", err)
// 		}
	
// 		fmt.Println("body\n")
// 		fmt.Println(r.Body)
// 		req_body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			fmt.Println("couldnt get body")
// 		}
// 		//req_body = bytes.NewBuffer(body)

// 	}
// 	//req, err:= http.NewRequest(r.Method, endpoint, req_body )
// 	req.Header.Set("Content-Type", "application/json")
	




// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Printf("error making http request %s", err)
// 		os.Exit(1)
// 	}

// 	fmt.Printf("rebponse status: %d\n", res.StatusCode)
// 	rebp_content, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Printf("could not parse %s", err)
// 	}
// 	//fmt.Printf("%s",rebp_content)
	
// 	var data map[string]interface{}
//     json.Unmarshal(rebp_content, &data)
//     //fmt.Printf("Results: %v\n", data)
//    	//return data
//    	jsonStr, err := json.Marshal(data)

	
// 	//fmt.Fprintf(w, "%q", jsonStr)
// 	w.Write(jsonStr)
// 	//fmt.Fprintf(w, "%q", string(jsonStr))


// }


type BackendHost struct {
	IP string
	IsAlive bool
	ReverseProxy *httputil.ReverseProxy
	mutex sync.RWMutex
}

type BackendPool struct {
	backends []*BackendHost
	current uint64
}
 
func (bp *BackendPool) RegisterBackend(b *BackendHost) {
	bp.backends = append(bp.backends, b)

}

func (bp *BackendPool) GetBackend() *BackendHost{
	indx := int(atomic.AddUint64(&bp.current, uint64(1)) % uint64(len(bp.backends)))
	
	for i:=indx; i < len(bp.backends)+indx; i++ {
		if GetIsAlive(bp.backends[i%len(bp.backends)]) {
			if i != indx {
				atomic.StoreUint64(&bp.current, uint64(i%len(bp.backends)))
			}
			return bp.backends[i%len(bp.backends)]
		}
	}
	return nil
}


func proxy_handler1(w http.ResponseWriter, r *http.Request)  {

	available_backend := backendPool.GetBackend()
	log.Print(r.Method, " ", r.URL.Path, " ", 
			available_backend.IP)

	if available_backend != nil {
		available_backend.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)

}


func make_backend_alive(b *BackendHost, is_alive bool) {
	b.mutex.Lock()
	b.IsAlive = is_alive
	b.mutex.Unlock()
}

func GetIsAlive(b *BackendHost) (is_alive bool) {
	b.mutex.RLock()
	is_alive = b.IsAlive
	b.mutex.RUnlock()
	return
}

func isBackendAlive(ip string) bool {
	timeout := 2 * time.Second
	u, err := url.Parse(ip)
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}

// HealthCheck pings the backends and update the status
func (s *BackendPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.IP)
		make_backend_alive(b, alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.IP, status)
	}
}

// HealthCheck pings the backends and update the status every minute
func healthCheck() {
	t := time.NewTicker(time.Minute * 1) 
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			backendPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}


func (bp *BackendPool) all_hosts_status() {
	for _ , b := range bp.backends {
		fmt.Println(b.IP)
		fmt.Println(b.IsAlive)
	}
}
func (bp *BackendPool) UpdateBackendStatus(target_ip string, status string) {


	for _ , b := range bp.backends {
		if b.IP == target_ip {
			if status == "down"{
				make_backend_alive(b, false)
				break
			} else if status == "up" {
				make_backend_alive(b, true)
				break
			}
		}
	}
}


func (bp *BackendPool) proxy_error_handler(w http.ResponseWriter, r *http.Request, err error) {

	target_url := "http://" + r.URL.Hostname() + ":" + r.URL.Port()
	log.Print("ERR: ", target_url, " is DOWN")
	bp.UpdateBackendStatus(target_url, "down")
	bp.all_hosts_status()
}

// func local_http_handler(w http.RebponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		log.Print("GET: ")
// 		fmt.Fprintf(w, "GET - beta, %q", html.EscapeString(r.URL.Path))
// 		node_ip := "http://"
// 		node_ip +=  get_backend_IP()
// 		node_ip += r.URL.Path
// 		fmt.Println(node_ip)
// 		rebp, err := http.Get(node_ip)
// 		if err != nil {
// 			fmt.Printf("error making http request: %s\n", err)
// 			os.Exit(1)
// 		} 
// 		fmt.Printf("rebponse status: %d", rebp.StatusCode)

// 	case http.MethodPost:
// 		fmt.Fprintf(w, "POST - beta, %q", html.EscapeString(r.URL.Path))
// 	default:
// 		http.Error(w, "Invalid request method.", 405)
//	}
// }



var backendPool BackendPool


func main() {
	fmt.Println("=================================================")
	fmt.Println(" Welcome to our in-house Load Balancer (Lginx) ")
    fmt.Println("=================================================")

    proxy_port := 80
    hosts := []string{
    		"http://127.0.0.1:8001",
    		"http://192.168.1.27:8001",
    		"http://192.168.1.27:8002",
    		"http://192.168.1.27:8003",
    		"http://192.168.1.27:8004",
    	}


   

    for _, ip := range hosts {
    	bkend_ip, err := url.Parse(ip)
    	if err != nil {
    		log.Fatal(err)
    	}
    	proxy := httputil.NewSingleHostReverseProxy(bkend_ip)
    	proxy.ErrorHandler = backendPool.proxy_error_handler
    	backendPool.RegisterBackend(
    		&BackendHost{
    			IP: 			ip,
    			IsAlive: 		true,
    			ReverseProxy: 	proxy,
    		})

    }


    http.HandleFunc("/", proxy_handler1)
    fmt.Printf("Lginx started at port %d\n", proxy_port)
    go healthCheck()

    log.Fatal(http.ListenAndServe(":80", nil))

	
	
	
}
