package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func SSH() {
	workspace := ReadWorkspace()

	var privateKey []byte

	privateKey, err := ioutil.ReadFile(workspace.PrivateKey)
	if err != nil {
		fmt.Println("Error while reading private key ", err)
		return
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		fmt.Println("Error while parsing private key ", err)
		return
	}

	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", workspace.Host+":22", config)

	if err != nil {
		fmt.Println("Error while dialing ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Error while creating session ", err)
		return
	}

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b

	if err := session.Run("chmod +x cj.sh && sed -i -e 's/\r$//' cj.sh && ./cj.sh"); err != nil {
		fmt.Println("Error while running command ", err)
		return
	}

	fmt.Println(b.String())

	fmt.Println("connection finished")
}
