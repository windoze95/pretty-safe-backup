package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Username string
	Host     string
	Port     string
	KeyPath  string
	Cmd      string
}

func (s *SSH) Run() (res string, err error) {
	privateKey, err := ioutil.ReadFile(s.KeyPath)
	if err != nil {
		err = fmt.Errorf("Unable to read private key: %s", err.Error())
		return
	}
	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		err = fmt.Errorf("Error parsing private key: %s", err.Error())
		return
	}
	conf := &ssh.ClientConfig{
		User:            s.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	// Connect
	authority := s.Host + ":" + s.Port
	client, err := ssh.Dial("tcp", authority, conf)
	if err != nil {
		err = fmt.Errorf("Failed to connect: %s", err.Error())
		return
	}
	// Create a session
	session, err := client.NewSession()
	if err != nil {
		err = fmt.Errorf("Failed to create session: %s", err.Error())
		return
	}
	defer session.Close()
	var outputBuffer bytes.Buffer
	session.Stdout = &outputBuffer
	// Run the command
	err = session.Run(s.Cmd)
	if err != nil {
		err = fmt.Errorf("SSH command '%s' failed: %s", s.Cmd, err.Error())
	}
	res = outputBuffer.String()
	return
}
