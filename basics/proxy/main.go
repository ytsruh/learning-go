package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type loadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func newLoadBalancer(port string, servers []Server) *loadBalancer {
	return &loadBalancer{
		port:            port,
		servers:         servers,
		roundRobinCount: 0,
	}
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (lb *loadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *loadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	log.Printf("Forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

func RunReverseProxy() {
	servers := []Server{
		newSimpleServer("https://www.ytsruh.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("https://www.duckduckgo.com"),
	}

	lb := newLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)
	log.Println("Serving request at port 8000")
	http.ListenAndServe(":"+lb.port, nil)
}
