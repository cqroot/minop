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
	"path"
	"path/filepath"
	"time"

	"github.com/cqroot/minop/pkg/host"
	"github.com/cqroot/minop/pkg/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Remote represents a SSH/SFTP client for remote server operations
type Remote struct {
	Hostname string
	Port     int
	Username string
	Password string
	Logger   *log.Logger
	client   *ssh.Client  // SSH client
	sftp     *sftp.Client // SFTP client
}

// New creates a new Remote instance and establishes connections
func New(h host.Host, logger *log.Logger) (*Remote, error) {
	remoteLogger := log.NewFromLogger(logger.With().Str("host", fmt.Sprintf("%s@%s:%d", h.User, h.Address, h.Port)).Logger())
	r := &Remote{
		Hostname: h.Address,
		Port:     h.Port,
		Username: h.User,
		Password: h.Password,
		Logger:   remoteLogger,
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

	// Format connection string and dial SSH
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.Address, h.Port), sshConfig)
	if err != nil {
		r.Logger.Error().Err(err).Msg("SSH dial error")
		return nil, fmt.Errorf("SSH dial error: %w", err)
	}
	r.client = conn

	// Create SFTP client
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		r.Logger.Error().Err(err).Msg("SFTP client error")
		conn.Close() // Close SSH connection if SFTP fails
		return nil, fmt.Errorf("SFTP client error: %w", err)
	}
	r.sftp = sftpClient

	return r, nil
}

func ToUnixPath(pathStr string) string {
	return path.Clean(filepath.ToSlash(pathStr))
}

// Close closes the SSH and SFTP connections
func (r *Remote) Close() error {
	var errs []error

	// Close SFTP client if it exists
	if r.sftp != nil {
		if err := r.sftp.Close(); err != nil {
			errs = append(errs, fmt.Errorf("SFTP close error: %w", err))
		}
	}

	// Close SSH client if it exists
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("SSH close error: %w", err))
		}
	}

	// Return combined errors if any occurred
	if len(errs) > 0 {
		return fmt.Errorf("multiple errors closing connections: %v", errs)
	}

	return nil
}

// ExecuteCommand executes a command on remote host and returns the output
func (r *Remote) ExecuteCommand(cmd string) (int, string, string, error) {
	session, err := r.client.NewSession()
	if err != nil {
		r.Logger.Error().Err(err).Msg("create session error")
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
		err = nil // Command executed but returned non-zero exit status
	} else if err != nil {
		// Other types of errors (connection issues, etc.)
		r.Logger.Error().Err(err).Msg("command execution error")
		return 0, "", "", fmt.Errorf("command execution error: %w", err)
	}

	return exitStatus, stdout.String(), stderr.String(), err
}

// determineOptimalBufferSize calculates optimal buffer size based on file size
func determineOptimalBufferSize(fileSize int64) int {
	// For small files (< 1MB), use 32KB buffer
	if fileSize < 1024*1024 {
		return 32 * 1024 // 32KB
	}
	// For medium files (1MB - 10MB), use 128KB buffer
	if fileSize < 10*1024*1024 {
		return 128 * 1024 // 128KB
	}
	// For large files (10MB - 100MB), use 512KB buffer
	if fileSize < 100*1024*1024 {
		return 512 * 1024 // 512KB
	}
	// For very large files (> 100MB), use 1MB buffer
	return 1024 * 1024 // 1MB
}

// UploadFile uploads a local file to remote path with buffer optimization
func (r *Remote) UploadFile(localPath, remotePath string) error {
	remotePath = ToUnixPath(remotePath)

	startTime := time.Now()
	r.Logger.Debug().
		Str("local", localPath).
		Str("remote", remotePath).
		Msg("starting file upload")

	defer func() {
		elapsed := time.Since(startTime)
		r.Logger.Info().
			Str("local", localPath).
			Str("remote", remotePath).
			Dur("elapsed", elapsed).
			Msg("file upload completed")
	}()

	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		r.Logger.Error().Err(err).Msg("open local file error")
		return fmt.Errorf("open local file error: %w", err)
	}
	defer localFile.Close()

	// Get file info to check size
	fileInfo, err := localFile.Stat()
	if err != nil {
		r.Logger.Error().Err(err).Msg("get file info error")
		return fmt.Errorf("get file info error: %w", err)
	}

	// Ensure remote directory exists
	remoteDir := ToUnixPath(filepath.Dir(remotePath))
	if err := r.ensureRemoteDir(remoteDir); err != nil {
		return fmt.Errorf("ensure remote directory error: %w", err)
	}

	// Create remote file
	remoteFile, err := r.sftp.Create(remotePath)
	if err != nil {
		r.Logger.Error().Err(err).Msg("create remote file error")
		return fmt.Errorf("create remote file error: %w", err)
	}
	defer remoteFile.Close()

	// Use buffered copy with optimal buffer size
	bufferSize := determineOptimalBufferSize(fileInfo.Size())
	_, err = io.CopyBuffer(remoteFile, localFile, make([]byte, bufferSize))
	if err != nil {
		r.Logger.Error().Err(err).Msg("copy file content error")
		return fmt.Errorf("copy file content error: %w", err)
	}

	r.Logger.Info().Str("local", localPath).Str("remote", remotePath).Msg("file uploaded successfully")
	return nil
}

// ensureRemoteDir ensures that the remote directory exists, creating it if necessary
func (r *Remote) ensureRemoteDir(remoteDir string) error {
	// Skip if directory is empty (root)
	if remoteDir == "" || remoteDir == "." || remoteDir == "/" {
		return nil
	}

	// Check if directory already exists
	fileInfo, err := r.sftp.Stat(remoteDir)
	if err == nil {
		if fileInfo.IsDir() {
			return nil
		}
		return fmt.Errorf("remote path exists but is not a directory: %s", remoteDir)
	}

	// If error is not "file doesn't exist", return it
	if !os.IsNotExist(err) {
		return fmt.Errorf("check remote directory error: %w", err)
	}

	// Create directory (and parent directories if needed)
	if err := r.sftp.MkdirAll(remoteDir); err != nil {
		// Double-check if directory was created by another process
		if _, checkErr := r.sftp.Stat(remoteDir); checkErr == nil {
			return nil
		}
		return fmt.Errorf("create remote directory error: %w", err)
	}

	return nil
}

// UploadDirectory uploads a local directory recursively to remote path with better error handling
func (r *Remote) UploadDirectory(localDir, remoteDir string) error {
	remoteDir = ToUnixPath(remoteDir)

	localInfo, err := os.Stat(localDir)
	if err != nil {
		r.Logger.Error().Err(err).Str("path", localDir).Msg("local directory error")
		return fmt.Errorf("local directory error: %w", err)
	}

	if !localInfo.IsDir() {
		return fmt.Errorf("local path is not a directory: %s", localDir)
	}

	if err := r.ensureRemoteDir(remoteDir); err != nil {
		return err
	}

	r.Logger.Debug().Str("local", localDir).Str("remote", remoteDir).Msg("starting directory upload")

	var uploadErrors []error

	// Walk through local directory recursively
	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			r.Logger.Warn().Err(err).Str("path", path).Msg("skip path due to error")
			uploadErrors = append(uploadErrors, err)
			return nil
		}

		// Calculate relative path from local directory root
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			r.Logger.Warn().Err(err).Str("path", path).Msg("get relative path error")
			uploadErrors = append(uploadErrors, err)
			return nil
		}

		// Convert to slash-separated path for remote server compatibility
		relPath = filepath.ToSlash(relPath)
		remotePath := filepath.ToSlash(filepath.Join(remoteDir, relPath))

		if info.IsDir() {
			// Ensure remote directory exists
			if err := r.ensureRemoteDir(remotePath); err != nil {
				r.Logger.Warn().Err(err).Str("path", remotePath).Msg("create remote directory error")
				uploadErrors = append(uploadErrors, err)
			}
			return nil
		}

		// Upload file
		r.Logger.Debug().Str("local", path).Str("remote", remotePath).Msg("uploading file")
		if err := r.UploadFile(path, remotePath); err != nil {
			r.Logger.Warn().Err(err).Str("path", path).Msg("upload file error")
			uploadErrors = append(uploadErrors, err)
		}

		return nil
	})
	if err != nil {
		uploadErrors = append(uploadErrors, err)
	}

	// Report errors if any occurred during upload
	if len(uploadErrors) > 0 {
		r.Logger.Error().Int("err_count", len(uploadErrors)).Msg("directory upload completed with errors")
		for i, err := range uploadErrors {
			if i < 5 {
				r.Logger.Error().Int("index", i).Err(err).Msg("")
			}
		}
		return fmt.Errorf("directory upload completed with %d errors", len(uploadErrors))
	}

	r.Logger.Info().Str("local", localDir).Str("remote", remoteDir).Msg("directory upload completed successfully")
	return nil
}
