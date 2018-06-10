package rsync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Rsync struct {
	From     string
	To       string
	Flags    string
	User     string
	Port     string
	Host     string
	Key      string
	Includes []string
	Excludes []string
	args     []string
}

func (r *Rsync) generateIncludeArgs() (iArgs []string) {
	for _, i := range r.Includes {
		iArgs = append(iArgs, "--include", fmt.Sprintf("'%s'", i))
	}
	return
}

func (r *Rsync) generateExcludeArgs() (eArgs []string) {
	for _, e := range r.Excludes {
		eArgs = append(eArgs, "--exclude", fmt.Sprintf("'%s'", e))
	}
	return
}

func (r *Rsync) generateScript() (data string) {
	r.loadArgs()
	scriptHead := "#!/bin/sh\n# This file is maintained by pretty-safe-backup, any manual input will be overwritten before it's ran"
	return fmt.Sprintf("%s\nrsync %s", scriptHead, strings.Join(r.args, " "))
}

func (r *Rsync) writeScript(file string) (err error) {
	err = ioutil.WriteFile(file, []byte(r.generateScript()), 0644)
	if err != nil {
		err = fmt.Errorf("Error writing Rsync script: %s", err.Error())
	}
	return
}

func (r *Rsync) loadArgs() {
	if r.Flags != "" {
		r.args = append(r.args, "-"+r.Flags)
	}
	r.args = append(r.args, "--rsync-path='sudo", "rsync'")
	r.args = append(r.args, "--delete")
	r.args = append(r.args, r.generateIncludeArgs()...)
	r.args = append(r.args, r.generateExcludeArgs()...)
	if r.Host != "" {
		r.args = append(r.args, "-e", fmt.Sprintf("'/usr/bin/ssh -i %s -p %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -l %s'", r.Key, r.Port, r.User))
		r.args = append(r.args, r.From+"/")
		r.args = append(r.args, fmt.Sprintf("%s@%s:'%s'", r.User, r.Host, r.To))
	} else {
		r.args = append(r.args, r.From+"/", r.To)
	}
}

func (r *Rsync) Run(configDir string, name string) (res string, err error) {
	file := filepath.Join(configDir, name+".sh")
	if _, e := os.Stat(file); !os.IsNotExist(e) {
		var scriptBuf []byte
		scriptBuf, err = ioutil.ReadFile(file)
		if err != nil {
			err = fmt.Errorf("Error reading Rsync script: %s", err.Error())
			return
		}
		if string(scriptBuf) != r.generateScript() {
			err = r.writeScript(file)
			if err != nil {
				return
			}
		}
	} else {
		err = r.writeScript(file)
		if err != nil {
			return
		}
	}
	rsyncCmd := exec.Command("/bin/sh", file)
	var outputBuffer bytes.Buffer
	rsyncCmd.Stdout = &outputBuffer
	err = rsyncCmd.Run()
	if err != nil {
		// Remove script if Rsync has an error
		os.Remove(file)
		err = fmt.Errorf("Rsync error: %s", err.Error())
	}
	return outputBuffer.String(), err
}
