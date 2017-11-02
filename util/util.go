// add methods to string??
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

func IsEmptyString(str string) bool {
	return str == ""
}
