package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shawntoffel/docker-secrets-provisioner/pkg/docker"
	"github.com/shawntoffel/docker-secrets-provisioner/pkg/provider/azurekv"
	"github.com/shawntoffel/docker-secrets-provisioner/pkg/provisioner"
)

// Version of dsp
var Version = ""

var (
	flagVersion    = false
	flagSourceID   = ""
	flagTargetName = ""
)

func parseCli() {
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.StringVar(&flagSourceID, "source-id", flagSourceID, "The source secret id")
	flag.StringVar(&flagTargetName, "target-name", flagTargetName, "The target secret name")

	flag.Parse()
}
func main() {
	parseCli()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	provider := azurekv.NewProviderFromEnv()
	dockerClient := docker.NewClientFromEnv()

	p := provisioner.New(provider, dockerClient)

	id, err := p.Provision(flagSourceID, flagTargetName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("created docker secret with ID: " + id)
}
