package main

import (
	"flag"
	"fmt"
	"os"
)

// Version of dsp
var Version = ""

var (
	flagVersion          = false
	flagDockerHost       = "http://0.0.0.0:2375"
	flagDockerAPIVersion = ""
)

func parseCli() {
	flag.BoolVar(&flagVersion, "v", false, "version")
	flag.StringVar(&flagDockerHost, "docker.host", flagDockerHost, "Docker host")
	flag.StringVar(&flagDockerAPIVersion, "docker.apiversion", flagDockerAPIVersion, "Docker API version")

	flag.Parse()
}
func main() {
	parseCli()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

}
