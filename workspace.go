package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Workspace struct {
	Name       string    `json:"name"`
	Ami        string    `json:"ami"`
	Type       string    `json:"type"`
	InstanceID string    `json:"instance_id"`
	Host       string    `json:"host"`
	CreatedAt  time.Time `json:"created_at"`
	PrivateKey string    `json:"private_key"`
}

func CreateWorkspace(workspace Workspace) {
	file, err := os.Create("cj.json")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	wJson, err := json.Marshal(workspace)
	if err != nil {
		fmt.Println(err)
	}

	_, err = file.WriteString(string(wJson))
	if err != nil {
		fmt.Println(err)
	}
}
