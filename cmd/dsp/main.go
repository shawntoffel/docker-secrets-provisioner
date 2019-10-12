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

var dockerHostEnvVar = "DOCKER_HOST"

var (
	flagVersion          = false
	flagDockerHost       = ""
	flagDockerAPIVersion = ""
	flagSourceName       = ""
	flagSourceVersion    = ""
	flagTargetName       = ""
)

func parseCli() {
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.StringVar(&flagDockerHost, "docker.host", flagDockerHost, "The Docker host. "+dockerHostEnvVar)
	flag.StringVar(&flagDockerAPIVersion, "docker.apiversion", flagDockerAPIVersion, "Docker API version")
	flag.StringVar(&flagSourceName, "source.name", flagSourceName, "The source secret name")
	flag.StringVar(&flagSourceVersion, "source.version", flagSourceVersion, "The source secret version")
	flag.StringVar(&flagTargetName, "target.name", flagTargetName, "The target secret name")

	if flagDockerHost == "" {
		flagDockerHost = os.Getenv(dockerHostEnvVar)
	}

	flag.Parse()
}
func main() {
	parseCli()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	provider := azurekv.NewProvider()
	dockerClient := docker.NewClient(flagDockerHost, flagDockerAPIVersion)

	p := provisioner.New(provider, dockerClient)

	id, err := p.Provision(flagSourceName, flagSourceVersion, flagTargetName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("created docker secret with ID: " + id)
}
