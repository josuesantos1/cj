package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Profile       string `hcl:"profile"`
	Region        string `hcl:"region"`
	Path          string `hcl:"path,optional"`
	Shell         string `hcl:"shell,optional"`
	KeyName       string `hcl:"keyName"`
	PrivateKey    string `hcl:"privateKey"`
	SecurityGroup string `hcl:"securityGroup"`
}

func Configure() Config {
	var config Config
	err := hclsimple.DecodeFile("cj.hcl", nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	fmt.Printf("%#v\n", config)
	return config
}
