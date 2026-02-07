package remote_test

import (
	"os"
	"testing"

	"github.com/cqroot/minop/pkg/remote"
	"github.com/stretchr/testify/require"
)

func IsFileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func TestExecuteCommand(t *testing.T) {
	hostsName := "./testdata/hosts"
	if !IsFileExists(hostsName) {
		t.Log("You need to create a hosts file for testing.")
		return
	}

	hg, err := remote.ParseHostsFile(hostsName)
	require.Nil(t, err)

	for _, hosts := range hg {
		for _, h := range hosts {
			t.Logf("Test ExecuteCommand on %s@%s:%d.\n", h.User, h.Address, h.Port)

			r, err := remote.New(h)
			require.Nil(t, err)

			ret, stdout, stderr, err := r.ExecuteCommand("echo 'minop test stdout'; echo 'minop test stderr' >&2")
			require.Nil(t, err)
			require.Equal(t, 0, ret)
			require.Equal(t, "minop test stdout\n", stdout)
			require.Equal(t, "minop test stderr\n", stderr)
		}
	}
}

func TestUploadFile(t *testing.T) {
	hostsName := "./testdata/hosts"
	if !IsFileExists(hostsName) {
		t.Log("You need to create a hosts file for testing.")
		return
	}

	testFileName := "./testdata/test.txt"
	if !IsFileExists(testFileName) {
		return
	}

	hg, err := remote.ParseHostsFile(hostsName)
	require.Nil(t, err)

	for _, hosts := range hg {
		for _, h := range hosts {
			t.Logf("Test UploadFile on %s@%s:%d.\n", h.User, h.Address, h.Port)

			r, err := remote.New(h)
			require.Nil(t, err)

			err = r.UploadFile(testFileName, "/root/minop.testfile")
			require.Nil(t, err)

			defer func() {
				ret, _, _, err := r.ExecuteCommand("rm -f /root/minop.testfile")
				require.Nil(t, err)
				require.Equal(t, 0, ret)
			}()

			ret, stdout, _, err := r.ExecuteCommand("[ -f /root/minop.testfile ] && echo 'minop test'")
			require.Nil(t, err)
			require.Equal(t, 0, ret)
			require.Equal(t, "minop test\n", stdout)

		}
	}
}

func TestUploadDir(t *testing.T) {
	hostsName := "./testdata/hosts"
	if !IsFileExists(hostsName) {
		t.Log("You need to create a hosts file for testing.")
		return
	}

	testDirName := "./testdata/test.dir"
	if !IsDirExists(testDirName) {
		return
	}

	hg, err := remote.ParseHostsFile(hostsName)
	require.Nil(t, err)

	for _, hosts := range hg {
		for _, h := range hosts {
			t.Logf("Test UploadDir on %s@%s:%d.\n", h.User, h.Address, h.Port)

			r, err := remote.New(h)
			require.Nil(t, err)

			err = r.UploadDir(testDirName, "/root/minop.testdir")
			require.Nil(t, err)

			defer func() {
				ret, _, _, err := r.ExecuteCommand("rm -rf /root/minop.testdir")
				require.Nil(t, err)
				require.Equal(t, 0, ret)
			}()

			ret, stdout, _, err := r.ExecuteCommand("[ -d /root/minop.testdir ] && echo 'minop test'")
			require.Nil(t, err)
			require.Equal(t, 0, ret)
			require.Equal(t, "minop test\n", stdout)

		}
	}
}
