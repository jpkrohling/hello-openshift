package main

import (
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	sp, _ := opentracing.StartSpanFromContext(r.Context(), "request")
	defer sp.Finish()

	fmt.Fprintln(w, "Hello")
	fmt.Println("Request processed.")
}

func main() {
	config := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	closer, err := config.InitGlobalTracer("hello-openshift", jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()

	http.HandleFunc("/", helloHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
