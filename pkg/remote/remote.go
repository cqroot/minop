/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package remote

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cqroot/minop/pkg/host"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Remote represents a SSH/SFTP client for remote server operations
type Remote struct {
	Hostname string
	Port     int
	Username string
	Password string
	client   *ssh.Client  // SSH client
	sftp     *sftp.Client // SFTP client
}

// New creates a new Remote instance and establishes connections
func New(h host.Host) (*Remote, error) {
	client := &Remote{
		Hostname: h.Address,
		Port:     h.Port,
		Username: h.User,
		Password: h.Password,
	}

	// Establish SSH connection
	sshConfig := &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.Address, h.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH dial error: %w", err)
	}
	client.client = conn

	// Create SFTP client
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("SFTP client error: %w", err)
	}
	client.sftp = sftpClient

	return client, nil
}

// ExecuteCommand executes a command on remote host and returns the output
func (r *Remote) ExecuteCommand(cmd string) (int, string, string, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return 0, "", "", fmt.Errorf("create session error: %w", err)
	}
	defer session.Close()

	var (
		exitStatus int = 0
		stdout     bytes.Buffer
		stderr     bytes.Buffer
	)
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	var e *ssh.ExitError
	if err != nil && errors.As(err, &e) {
		exitStatus = e.ExitStatus()
		err = nil
	}

	return exitStatus, stdout.String(), stderr.String(), err
}

// UploadFile uploads a local file to remote path
func (r *Remote) UploadFile(localPath, remotePath string) error {
	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local file error: %w", err)
	}
	defer localFile.Close()

	// Create remote file
	remoteFile, err := r.sftp.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create remote file error: %w", err)
	}
	defer remoteFile.Close()

	// Copy file content
	_, err = io.Copy(remoteFile, localFile)
	return err
}

// DownloadFile downloads a remote file to local path
func (r *Remote) DownloadFile(remotePath, localPath string) error {
	// Open remote file
	remoteFile, err := r.sftp.Open(remotePath)
	if err != nil {
		return fmt.Errorf("open remote file error: %w", err)
	}
	defer remoteFile.Close()

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create local file error: %w", err)
	}
	defer localFile.Close()

	// Copy file content
	_, err = io.Copy(localFile, remoteFile)
	return err
}

// UploadDirectory uploads a local directory recursively to remote path
func (r *Remote) UploadDirectory(localDir, remoteDir string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(localDir, path)
		remotePath := filepath.Join(remoteDir, relPath)

		if info.IsDir() {
			return r.sftp.Mkdir(remotePath)
		}

		return r.UploadFile(path, remotePath)
	})
}

// DownloadDirectory downloads a remote directory recursively to local path
func (r *Remote) DownloadDirectory(remoteDir, localDir string) error {
	walker := r.sftp.Walk(remoteDir)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}

		remotePath := walker.Path()
		relPath, _ := filepath.Rel(remoteDir, remotePath)
		localPath := filepath.Join(localDir, relPath)

		if walker.Stat().IsDir() {
			err := os.MkdirAll(localPath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := r.DownloadFile(remotePath, localPath); err != nil {
			return err
		}
	}
	return nil
}

// Close closes all connections
func (r *Remote) Close() error {
	if r.sftp != nil {
		r.sftp.Close()
	}
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
