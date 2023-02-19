
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
	"os"
	"io/ioutil"
	//"encoding/json"
	"sync"
	"sync/atomic"
	"net/http/httputil"
	"net"
	"lginx/cache"
	"github.com/bitly/go-simplejson"
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


var CACHE_OBJ cache.Cache
func init() {
	fmt.Println("Main setuping cache.")
	var hmap = make(map[string]string)
   	//f := new(parent.Father)
	CACHE_OBJ = (cache.Cache{Capacity: 8, Storage: hmap, Current_size: 0})
	CACHE_OBJ.Put("Manish", "sjf")

}

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
	req_log := r.Method +  "_" + r.URL.Path +  "_" + available_backend.IP
	fmt.Println(r.URL.Query())

	log.Print(req_log)

	if available_backend != nil {
		w.Header().Add("url", req_log)
		available_backend.ReverseProxy.ServeHTTP(NewCustomWriter(w), r)
		//fmt.Println(w.Body)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)

}
type customWriter struct {
    http.ResponseWriter
}

func NewCustomWriter(w http.ResponseWriter) *customWriter {
    return &customWriter{w}
}

func (c *customWriter) Header() http.Header {
    return c.ResponseWriter.Header()
}

func (c *customWriter) Write(data []byte) (int, error) {
    fmt.Println((c.Header())) //get response here
    //CACHE_OBJ.put()
    return c.ResponseWriter.Write(data)
}

func (c *customWriter) WriteHeader(i int) {
    c.ResponseWriter.WriteHeader(i)
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

    configFile, err := os.Open("config.json")
    if err != nil {
    	log.Fatal("config file read err", err)
    }
    //defer configFile.Close()
    byteValue, _ := ioutil.ReadAll(configFile)

    // var configMap map[string]interface{}
    // json.Unmarshal([]byte(byteValue), &configMap)
    // //fmt.Println("configMap", configMap["backend_hosts"])
    
    // hosts := configMap["backend_hosts"].(string)
    // //hosts  := configMap["backend_hosts"]
    // proxy_port := configMap["default_lginx_port"].(int)

    js, err := simplejson.NewJson(byteValue)
    if err != nil {
    	fmt.Print("err")
    }
    fmt.Println("print")
    //fmt.Println(hosts, proxy_port)


    //proxy_port := 81
   	//hosts := []string{
   	//	"http://127.0.0.1:8001",
   	//}
   	hosts := js.Get("backend_hosts").MustStringArray()
   	lginxPort := js.Get("default_lginx_port").MustString()
   	fmt.Println(hosts)
   
    for _, ip := range hosts {
    	bkend_ip, err := url.Parse(ip)
    	if err != nil {
    		log.Fatal(err)
    	}
    	fmt.Println(bkend_ip)
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
    fmt.Printf("Lginx started at port %s\n", lginxPort)
    go healthCheck()

    log.Fatal(http.ListenAndServe(":"+lginxPort, nil))

	
	
	
}
