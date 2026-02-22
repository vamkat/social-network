package pprof

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	tele "social-network/shared/go/telemetry"
)

// http://<port>/debug/pprof/
func StartPprof(port string) {
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			tele.Fatalf("pprof error: %v", err)
		} else {
			tele.Info(context.Background(), "prof server running: @1", "port", port)
		}
	}()
}
