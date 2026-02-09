package pprof

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// http://<port>/debug/pprof/
func StartPprof(port string) {
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatalf("pprof error: %v", err)
		}
	}()
}
