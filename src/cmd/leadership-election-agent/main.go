package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"code.cloudfoundry.org/go-envstruct"
	"code.cloudfoundry.org/go-loggregator/metrics"
	"code.cloudfoundry.org/leadership-election/app/agent"
)

func main() {
	log.Printf("Starting Leadership Election...")
	defer log.Printf("Closing Leadership Election...")

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	envstruct.WriteReport(&cfg)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	m := metrics.NewRegistry(
		logger,
		metrics.WithTLSServer(
			int(cfg.MetricsServer.Port),
			cfg.MetricsServer.CertFile,
			cfg.MetricsServer.KeyFile,
			cfg.MetricsServer.CAFile,
		),
	)

	a := agent.New(
		cfg.NodeIndex,
		cfg.NodeAddrs,
		agent.WithLogger(logger),
		agent.WithMetrics(m),
		agent.WithPort(int(cfg.Port)),
	)

	a.Start(
		cfg.CAFile,
		cfg.CertFile,
		cfg.KeyFile,
	)

	// health endpoints (pprof and expvar)
	log.Printf("Health: %s", http.ListenAndServe(fmt.Sprintf("localhost:%d", cfg.HealthPort), nil))
}
