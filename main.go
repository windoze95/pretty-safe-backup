package main

import "github.com/orange-lightsaber/pretty-safe-backup/cmd"

var (
	// VERSION is set during build
	VERSION = "0.1.0"
)

func main() {
	cmd.Exec(VERSION)
}
