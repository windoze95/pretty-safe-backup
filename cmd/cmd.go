package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/orange-lightsaber/pretty-safe-backup/run"
)

func Exec(version string) {
	v := flag.Bool("v", false, "Print version.")
	loadProfileDefault := "/path/to/profile"
	loadProfile := flag.String("L", loadProfileDefault, "Load a profile.")
	flag.Parse()
	if *v {
		fmt.Printf("psb v%s\n", version)
		os.Exit(0)
	}
	if *loadProfile != loadProfileDefault {
		rc, err := run.DecodeRunConfig(*loadProfile)
		if err != nil {
			log.Fatal(err.Error())
		}
		newRunConfigFile, err := rc.WriteRunConfig()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Run config created: %s\n", newRunConfigFile)
		os.Exit(0)
	}
	run.Daemonize()
}
