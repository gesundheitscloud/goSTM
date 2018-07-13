package main

import (

	//"fmt"

	"io"

	"github.com/kevinburke/ssh_config"
)

func readSSHConfig(cfgFile io.Reader) (*ssh_config.Config, error) {
	cfg, err2 := ssh_config.Decode(cfgFile)
	/*for _, host := range cfg.Hosts {
		fmt.Println("patterns:", host.Patterns)
		for _, node := range host.Nodes {
			fmt.Println(node.String())
		}
	}*/

	// Write the cfg back to disk:
	//fmt.Println(cfg.String())
	return cfg, err2
}
