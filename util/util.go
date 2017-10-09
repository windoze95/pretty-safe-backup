package util

import (
	"os"
	"os/exec"
)

func ClearClient() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
