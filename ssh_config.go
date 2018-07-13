package main

import (

	//"fmt"

	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinburke/ssh_config"
)

func getDefaultConfig() (*ssh_config.Config, error) {
	fmt.Println("Read config from default")
	freader, ferr := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "ssh_config_example"))
	if ferr != nil {
		panic(ferr)
	}
	return readSSHConfig(bufio.NewReader(freader))
}

func getConfigFromURL() (*ssh_config.Config, error) {
	resp, herr := http.Get(goSTMUI.ConfigField.Text())
	if herr != nil {
		panic(herr)
	}
	defer resp.Body.Close()
	cfgBytes, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		panic(rerr)
	}
	cfg := string(cfgBytes)
	fmt.Println(cfg)
	return readSSHConfig(strings.NewReader(cfg))
}

func getConfigFromLocalPath() (*ssh_config.Config, error) {
	freader, ferr := os.Open(goSTMUI.ConfigField.Text())
	if ferr != nil {
		panic(ferr)
	}
	return readSSHConfig(bufio.NewReader(freader))
}

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
