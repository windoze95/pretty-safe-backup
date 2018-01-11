package rsync

import (
	"bytes"
	"fmt"
	"os/exec"
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
		iArgs = append(iArgs, "--include", i)
	}
	return
}

func (r *Rsync) generateExcludeArgs() (eArgs []string) {
	for _, e := range r.Excludes {
		eArgs = append(eArgs, "--exclude", e)
	}
	return
}

func (r *Rsync) loadArgs() {
	if r.Flags != "" {
		r.args = append(r.args, "-"+r.Flags)
	}
	r.args = append(r.args, r.generateIncludeArgs()...)
	r.args = append(r.args, r.generateExcludeArgs()...)
	if r.Host != "" {
		r.args = append(r.args, "-e", fmt.Sprintf("/usr/bin/ssh -i %s -p %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -l %s", r.Key, r.Port, r.User))
		r.args = append(r.args, r.From)
		r.args = append(r.args, fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.To))
		return
	}
	r.args = append(r.args, r.From, r.To)
}

func (r *Rsync) Run() (string, error) {
	r.loadArgs()
	rsyncCmd := exec.Command("rsync", r.args...)
	var outputBuffer bytes.Buffer
	rsyncCmd.Stdout = &outputBuffer
	err := rsyncCmd.Run()
	if err != nil {
		err = fmt.Errorf("Rsync error:", err.Error())
	}
	return outputBuffer.String(), err
}
