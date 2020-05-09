// Package static contains utility to work with ThundeRatz
// static files server https://static.thunderatz.org/
package static

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Conn represents a connection to the static server
type Conn struct {
	BasePath string

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// New initializes a new connection to the static server
func New(addr, user, password, base string, port int) (*Conn, error) {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", addr, port), &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &Conn{
		BasePath: base,

		sshClient:  conn,
		sftpClient: client,
	}, nil
}

// Mkdir creates the sirectory in the remote path, relative to static base path
func (sc *Conn) Mkdir(path string, makeParents bool) error {
	fullPath := sc.BasePath + path

	if makeParents {
		return sc.sftpClient.MkdirAll(fullPath)
	}
	return sc.sftpClient.Mkdir(fullPath)
}

// List returns a list of files form path realtive to static base path
func (sc *Conn) List(path string) ([]os.FileInfo, error) {
	fullPath := sc.BasePath + "/" + path

	return sc.sftpClient.ReadDir(fullPath)
}

// Get downloads a file from the remote path, relative to the static base path
func (sc *Conn) Get(remote string) error {
	paths := strings.Split(remote, "/")
	filename := paths[len(paths)-1]
	fullPath := sc.BasePath + "/" + remote

	dstFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := sc.sftpClient.Open(fullPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	fmt.Printf("Got %d bytes\n", bytes)

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// Put uploads a file to the remote path, relative to the static base path
func (sc *Conn) Put(local, remote string, overwrite bool) error {
	paths := strings.Split(local, "/")
	filename := paths[len(paths)-1]

	fullPath := sc.BasePath + "/" + remote + "/" + filename

	file, _ := sc.sftpClient.Stat(fullPath)

	if file != nil && !overwrite {
		return fmt.Errorf("File already exists")
	}

	srcFile, err := os.Open(local)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := sc.sftpClient.Create(fullPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	fmt.Printf("Put %d bytes\n", bytes)

	return nil
}

// Close ends the connection to the static server
func (sc *Conn) Close() {
	sc.sshClient.Close()
	sc.sftpClient.Close()
}
