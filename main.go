package main

import (
	"fmt"
	"log"

	"flag"

	"github.com/APITeamLimited/agent/agent"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	pProfPort := flag.Int("pprof-port", 0, "Enable pprof on the given port")

	flag.Parse()

	// If pprof is enabled, start the profiling server
	if *pProfPort != 0 {
		fmt.Printf("Starting pprof server on port %d\n", *pProfPort)
		go func() {
			log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", *pProfPort), nil))
		}()
	}

	agent.Run()
}
