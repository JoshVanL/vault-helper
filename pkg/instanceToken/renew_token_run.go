package instanceToken

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

const FlagTokenRole = "role"
const FlagVaultConfigPath = "config-path"

func (i *InstanceToken) Run(cmd *cobra.Command, args []string) error {

	if len(args) > 0 {
		i.SetClusterID(args[0])
	} else {
		return errors.New("No cluster id was given")
	}

	value, err := cmd.PersistentFlags().GetString(FlagTokenRole)
	if err != nil {
		return fmt.Errorf("Error parsing %s '%s': %s", FlagTokenRole, value, err)
	}
	if value == "" {
		return fmt.Errorf("No token role was given. Token role is required for this command:\n --%s", FlagTokenRole)
	}
	i.SetRole(value)

	value, err = cmd.PersistentFlags().GetString(FlagVaultConfigPath)
	if err != nil {
		return fmt.Errorf("Error parsing %s '%s': %s", FlagVaultConfigPath, value, err)
	}
	if value != "" {
		i.SetVaultConfigPath(value)
	}

	return i.TokenRenewRun()

}