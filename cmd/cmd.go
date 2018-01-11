package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/orange-lightsaber/pretty-safe-backup/run"
)

func Exec(version string) {
	v := flag.Bool("v", false, "Print version.")
	flag.Parse()
	if *v {
		fmt.Printf("psb v%s\n", version)
		os.Exit(0)
	}
	run.Daemonize()
}
