package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Endpoint ..
type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// SSHtunnel ..
type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

// Start a ssh tunnel
func (tunnel *SSHtunnel) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		go func() {
			select {
			case <-ctx.Done():
				listener.Close()
				return
			}
		}()
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	fmt.Printf("Dial to %s", tunnel.Server.String())
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}

		defer serverConn.Close()
		defer remoteConn.Close()
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// SSHAgent ..
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// PublicKeyFile ..
func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func start(ctx context.Context, sshcfg *ssh_config.Config, sshHost string) {
	/*
	   Host my-tunnel
	     Hostname 52.233.225.199
	     User gesund
	     IdentityFile ~/.ssh/id_rsa
	     DynamicForward 8080
	     ControlMaster auto
	     ControlPath ~/.ssh/sockets/%r@%h:%p
	     UserKnownHostsFile /dev/null
	     StrictHostKeyChecking no
	     LocalForward 127.0.0.1:9906 127.0.0.1:3306
	*/
	localForward, _ := sshcfg.Get(sshHost, "LocalForward")

	if len(localForward) <= 0 {
		panic("Missing LocalForward ssh configuration.")
	}

	localForwards := strings.Split(localForward, " ")

	if len(localForwards) < 2 {
		panic("Missing LocalForward local and remote definition.")
	}

	localForward, remoteForward := localForwards[0], localForwards[1]
	fmt.Println(localForward, remoteForward)

	host, strPort, err := net.SplitHostPort(localForward)
	fmt.Println(host, strPort)
	port, err := strconv.Atoi(strPort)
	if err != nil {
		panic(err)
	}
	// Read hosts and users from ssh config
	localEndpoint := &Endpoint{
		Host: host,
		Port: port, // "localForwardPort"
	}

	remoteHostname, _ := sshcfg.Get(sshHost, "Hostname")
	fmt.Println(remoteHostname)
	serverEndpoint := &Endpoint{
		Host: remoteHostname,
		Port: 22, //default ssh port
	}

	host, strPort, err = net.SplitHostPort(remoteForward)

	port, err = strconv.Atoi(strPort)
	if err != nil {
		panic(err)
	}
	//TODO localForwardCfg must be parsed into
	//remote localfowrad [optional ip]:port dest_ip:port
	remoteEndpoint := &Endpoint{
		Host: host, // get from localForwardCfg
		Port: port, // localForwardCfg
	}

	username, _ := sshcfg.Get(sshHost, "User")
	fmt.Println(username)

	identityFile, _ := sshcfg.Get(sshHost, "IdentityFile")

	fmt.Println(identityFile)
	// Check in case of paths like "/something/~/something/"
	if strings.HasPrefix(identityFile, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		identityFile = filepath.Join(dir, identityFile[2:])
	}

	fmt.Println(identityFile)
	publicKey, err := PublicKeyFile(identityFile)
	if err != nil {
		panic(err)
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		/*Auth: []ssh.AuthMethod{
			SSHAgent(),
		},*/
	}

	fmt.Println()

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	go tunnel.Start(ctx)
}
