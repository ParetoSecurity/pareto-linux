package main

import (
	"github.com/caarlos0/log"
	"go.uber.org/automaxprocs/maxprocs"

	"paretosecurity.com/auditor/cmd"
	"paretosecurity.com/auditor/shared"
)

func init() {

	// automatically set GOMAXPROCS to match available CPUs.
	// GOMAXPROCS will be used as the default value for the --parallelism flag.
	if _, err := maxprocs.Set(); err != nil {
		log.WithError(err).Fatal("failed to set GOMAXPROCS")
	}

}
func main() {
	if err := shared.LoadConfig(); err != nil {
		log.WithError(err).Warn("failed to load config")
	}
	cmd.Execute()
}
