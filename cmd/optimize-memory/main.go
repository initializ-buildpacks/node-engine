package main

import (
	"log"
	"os"

	"github.com/initializ-buildpacks/node-engine/cmd/optimize-memory/internal"
	"github.com/initializ-buildpacks/node-engine/cmd/util"
)

func main() {
	err := internal.Run(util.LoadEnvironmentMap(os.Environ()), os.NewFile(3, "/dev/fd/3"), "/")
	if err != nil {
		log.Fatal(err)
	}
}
