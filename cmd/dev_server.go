// Copyright Jetstack Ltd. See LICENSE for details.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/coreos/go-systemd/daemon"
	"github.com/jetstack/vault-helper/pkg/dev_server"
	"github.com/jetstack/vault-helper/pkg/kubernetes"
)

// initCmd represents the init command
var devServerCmd = &cobra.Command{
	Use:   "dev-server [cluster ID]",
	Short: "Run a vault server in development mode with kubernetes PKI created.",
	Run: func(cmd *cobra.Command, args []string) {
		log, err := LogLevel(cmd)
		if err != nil {
			Must(err)
		}

		if len(args) < 1 {
			Must(errors.New("no cluster ID was given"))
		}

		wait, err := cmd.PersistentFlags().GetBool(dev_server.FlagWaitSignal)
		if err != nil {
			Must(fmt.Errorf("error finding wait value: %v", err))
		}

		port, err := cmd.PersistentFlags().GetInt(dev_server.FlagPortNumber)
		if err != nil {
			Must(fmt.Errorf("error finding port value: %v", err))
		}
		if port > 65536 {
			Must(fmt.Errorf("invalid port %d > 65536", port))
		}
		if port < 1 {
			Must(fmt.Errorf("invalid port %d < 1", port))
		}

		v := dev_server.New(log)
		v.Vault.SetPort(port)
		if err := v.Vault.Start("vault"); err != nil {
			Must(fmt.Errorf("unable to initialise dev vault: %s", err))
		}

		v.Kubernetes = kubernetes.New(v.Vault.Client(), v.Log)
		v.Kubernetes.SetClusterID(args[0])

		if err := setFlagsKubernetes(v.Kubernetes, cmd); err != nil {
			Must(err)
		}

		if err := v.Kubernetes.Ensure(); err != nil {
			Must(err)
		}

		for n, t := range v.Kubernetes.InitTokens() {
			log.Infof(n + "-init_token := " + t)
		}

		daemon.SdNotify(false, "READY=1")

		if wait {
			waitSignal(v)
		}
	},
}

func init() {
	devServerCmd.PersistentFlags().Duration(kubernetes.FlagMaxValidityCA, time.Hour*24*365*20, "Maxium validity for CA certificates")
	devServerCmd.Flag(kubernetes.FlagMaxValidityCA).Shorthand = "c"

	devServerCmd.PersistentFlags().Duration(kubernetes.FlagMaxValidityAdmin, time.Hour*24*365, "Maxium validity for admin certificates")
	devServerCmd.Flag(kubernetes.FlagMaxValidityAdmin).Shorthand = "d"

	devServerCmd.PersistentFlags().Duration(kubernetes.FlagMaxValidityComponents, time.Hour*24*30, "Maxium validity for component certificates")
	devServerCmd.Flag(kubernetes.FlagMaxValidityComponents).Shorthand = "s"

	devServerCmd.PersistentFlags().String(kubernetes.FlagInitTokenEtcd, "", "Set init-token-etcd   (Default to new token)")
	devServerCmd.Flag(kubernetes.FlagInitTokenEtcd).Shorthand = "e"

	devServerCmd.PersistentFlags().String(kubernetes.FlagInitTokenWorker, "", "Set init-token-worker (Default to new token)")
	devServerCmd.Flag(kubernetes.FlagInitTokenWorker).Shorthand = "o"

	devServerCmd.PersistentFlags().String(kubernetes.FlagInitTokenMaster, "", "Set init-token-master (Default to new token)")
	devServerCmd.Flag(kubernetes.FlagInitTokenMaster).Shorthand = "m"

	devServerCmd.PersistentFlags().String(kubernetes.FlagInitTokenAll, "", "Set init-token-all    (Default to new token)")
	devServerCmd.Flag(kubernetes.FlagInitTokenAll).Shorthand = "a"

	devServerCmd.PersistentFlags().Bool(dev_server.FlagWaitSignal, true, "Wait for TERM + QUIT signal has been given before termination")
	devServerCmd.Flag(dev_server.FlagWaitSignal).Shorthand = "w"

	devServerCmd.PersistentFlags().Int(dev_server.FlagPortNumber, 8200, "Set the port number to connect to vault")
	devServerCmd.Flag(dev_server.FlagPortNumber).Shorthand = "t"

	RootCmd.AddCommand(devServerCmd)
}

func waitSignal(v *dev_server.DevVault) {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-signal_chan
			switch s {
			case syscall.SIGTERM:
				exit_chan <- 0

			case syscall.SIGQUIT:
				exit_chan <- 0
			}
		}
	}()

	<-exit_chan
	v.Vault.Stop()
}
