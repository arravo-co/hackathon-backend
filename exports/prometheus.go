package exports

import "github.com/prometheus/client_golang/prometheus"

var MyFirstCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ping_request_count",
	Help: "No of requests handled by Hello handler",
})
