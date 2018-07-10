package main

import (
	"os"
	//"fmt"

	"github.com/kevinburke/ssh_config"
)

func readSSHConfig(cfgFile string) (*ssh_config.Config, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		panic(err)
	}
	cfg, err2 := ssh_config.Decode(f)
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