package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func SCP() {
	workspace := ReadWorkspace()

	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey("ubuntu", workspace.PrivateKey, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(workspace.Host+":22", &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	file := "cj.sh"

	// Open a file
	f, _ := os.Open(file)

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
	err = client.CopyFromFile(context.Background(), *f, file, "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
		return
	}

	fmt.Println("File copied successfully")
}
