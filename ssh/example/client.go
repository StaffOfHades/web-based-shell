package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	hostPrivateKeySigner ssh.Signer
)

func init() {
	keyPath := "./user_key"
	if os.Getenv("USER_KEY") != "" {
		keyPath = os.Getenv("USER_KEY")
	}

	hostPrivateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}

	hostPrivateKeySigner, err = ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		panic(err)
	}
}

func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Println(conn.RemoteAddr(), "authenticate with", key.Type())
	return nil, nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <user> <host:port> <command>", os.Args[0])
	}

	client, session, err := connectToHost(os.Args[1], os.Args[2])
	// ERROR IS HERE
	if err != nil {
		panic(err)
	}
	out, err := session.CombinedOutput(os.Args[3])
	print("Output: ")
	fmt.Println(string(out))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	client.Close()
}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	var pass string
	fmt.Print("Password: ")
	fmt.Scanf("%s\n", &pass)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		log.Println("Couldn't dail")
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		log.Println("Couldn't crate new session")
		return nil, nil, err
	}

	return client, session, nil
}
