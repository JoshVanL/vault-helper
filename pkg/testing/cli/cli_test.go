// Copyright Jetstack Ltd. See LICENSE for details.
package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"

	"github.com/jetstack/vault-helper/cmd"
	"github.com/jetstack/vault-helper/pkg/testing/vault_dev"
)

const (
	binPath = "src/github.com/jetstack/vault-helper/vault-helper_linux_amd64"
)

var tmpDirs []string

func TestMain(m *testing.M) {

	vault, err := InitVaultDev()
	if err != nil {
		logrus.Fatalf("failed to initiate vault for testing: %v", err)
	}
	logrus.RegisterExitHandler(vault.Stop)

	if err := InitKubernetes(); err != nil {
		logrus.Fatalf("failed to initiate kubernetes for testing: %v", err)
	}

	exitCode := m.Run()

	if err := CleanDirs(); err != nil {
		logrus.Errorf("error cleaning up tmp dirs: %v", err)
	}

	vault.Stop()
	os.Exit(exitCode)
}

func InitVaultDev() (*vault_dev.VaultDev, error) {
	vaultDev := vault_dev.New()

	if err := vaultDev.Start(); err != nil {
		return nil, fmt.Errorf("unable to initialise vault dev server for testing: %v", err)
	}

	addr := fmt.Sprintf("http://127.0.0.1:%d", vaultDev.Port())

	if err := os.Setenv("VAULT_ADDR", addr); err != nil {
		vaultDev.Stop()
		return nil, fmt.Errorf("failed to set vault address environment variable: %v", err)
	}

	if err := os.Setenv("VAULT_TOKEN", "root-token-dev"); err != nil {
		vaultDev.Stop()
		return nil, fmt.Errorf("failed to set vault root token environment variable: %v", err)
	}

	return vaultDev, nil
}

func InitKubernetes() error {
	cmd.RootCmd.SetArgs([]string{
		"setup",
		"test",
		"--init-token-all=all",
		"--init-token-master=master",
		"--init-token-worker=worker",
		"--init-token-etcd=etcd",
	})
	cmd.RootCmd.Execute()

	return nil
}

func RunTest(args []string, pass bool, t *testing.T) {
	dir, err := initTokensDir()
	if err != nil {
		t.Errorf("failed to create tokens directory: %v", err)
		return
	}
	args = append(args, fmt.Sprintf("--config-path=%s", dir))
	cmd.RootCmd.SetArgs(args)

	cmd.Must = func(err error) {
		if err != nil && pass {
			t.Errorf("unexpected error: %v\nargs: %v", err, args)
			return
		}

		if err == nil && !pass {
			t.Errorf("expected error: got=none\nargs: %v", args)
		}
	}

	cmd.RootCmd.Execute()
}

func TmpDir() (string, error) {
	dir, err := ioutil.TempDir("", "test-cluster")
	if err != nil {
		return dir, fmt.Errorf("failed to create token directory: %v", err)
	}
	tmpDirs = append(tmpDirs, dir)

	return dir, nil
}

func CleanDirs() error {
	var result *multierror.Error

	for _, dir := range tmpDirs {
		if err := os.RemoveAll(dir); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func initTokensDir() (string, error) {
	dir, err := TmpDir()
	if err != nil {
		return dir, err
	}

	initTokenFile := fmt.Sprintf("%s/init-token", dir)
	tokenFile := fmt.Sprintf("%s/token", dir)

	if err := ioutil.WriteFile(initTokenFile, []byte("root-token-dev"), 0644); err != nil {
		return dir, fmt.Errorf("failed to write root-token-dev token to file: %v", err)
	}

	f, err := os.Create(tokenFile)
	if err != nil {
		return dir, fmt.Errorf("failed to create token file: %v", err)
	}
	if err := f.Close(); err != nil {
		return dir, fmt.Errorf("failed to close token file: %v", err)
	}

	return dir, nil
}
