// Copyright Jetstack Ltd. See LICENSE for details.
package kubernetes

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

const (
	pkiType = "pki"
)

type PKI struct {
	pkiName    string
	kubernetes *Kubernetes

	MaxLeaseTTL     time.Duration
	DefaultLeaseTTL time.Duration

	Log *logrus.Entry
}

func NewPKI(k *Kubernetes, pkiName string, logger *logrus.Entry) *PKI {
	return &PKI{
		pkiName:         pkiName,
		kubernetes:      k,
		MaxLeaseTTL:     k.MaxValidityCA,
		DefaultLeaseTTL: k.MaxValidityCA,
		Log:             logger,
	}
}

func (p *PKI) TuneMount(mount *vault.MountOutput) error {
	if p.TuneMountRequired(mount) {
		mountConfig := p.getMountConfigInput()
		err := p.kubernetes.vaultClient.Sys().TuneMount(p.Path(), mountConfig)
		if err != nil {
			return fmt.Errorf("error tuning mount config: %v", err.Error())
		}
		p.Log.Debugf("Tuned Mount: %s", p.pkiName)
		return nil
	}
	p.Log.Debugf("No tune required: %s", p.pkiName)

	return nil
}

func (p *PKI) TuneMountRequired(mount *vault.MountOutput) bool {

	if mount.Config.DefaultLeaseTTL != int(p.DefaultLeaseTTL.Seconds()) {
		return true
	}
	if mount.Config.MaxLeaseTTL != int(p.MaxLeaseTTL.Seconds()) {
		return true
	}

	return false
}

func (p *PKI) Ensure() error {
	mount, err := GetMountByPath(p.kubernetes.vaultClient, p.Path())
	if err != nil {
		return err
	}

	// Mount doesn't Exist
	if mount == nil {
		p.Log.Debugf("No mounts found for: %s", p.pkiName)
		err := p.kubernetes.vaultClient.Sys().Mount(
			p.Path(),
			&vault.MountInput{
				Description: "Kubernetes " + p.kubernetes.clusterID + "/" + p.pkiName + " CA",
				Type:        pkiType,
			},
		)

		if err != nil {
			return fmt.Errorf("failed to create mount: %v", err)
		}
		mount, err = GetMountByPath(p.kubernetes.vaultClient, p.Path())
		if err != nil {
			return err
		}
		p.Log.Infof("Mounted '%s'", p.pkiName)

	} else {
		if mount.Type != pkiType {
			return fmt.Errorf("Mount '%s' already existing with wrong type '%s'", p.Path(), mount.Type)
		}
		p.Log.Debugf("Mount '%s' already existing", p.Path())
	}

	if mount != nil {
		err = p.TuneMount(mount)
		if err != nil {
			return errors.New("failed to tune mount")
		}
	}

	return p.ensureCA()
}

func (p *PKI) EnsureDryRun() (bool, error) {
	mount, err := GetMountByPath(p.kubernetes.vaultClient, p.Path())
	if err != nil {
		return false, err
	}

	// Mount doesn't Exist
	if mount == nil {
		return true, nil

	} else {
		if mount.Type != pkiType {
			return true, nil
		}
	}

	if p.TuneMountRequired(mount) {
		return true, nil
	}

	exist, err := p.caPathExists()
	if err != nil {
		return false, err
	}

	if !exist {
		return true, nil
	}

	return false, nil
}

func (p *PKI) ensureCA() error {
	b, err := p.caPathExists()
	if err != nil {
		return err
	}

	if !b {
		return p.generateCA()
	}

	return nil
}

func (p *PKI) generateCA() error {
	description := "Kubernetes " + p.kubernetes.clusterID + "/" + p.pkiName + " CA"

	data := map[string]interface{}{
		"common_name": description,
		"ttl":         p.getMaxLeaseTTL(),
		"exclude_cn_from_sans": true,
	}

	_, err := p.kubernetes.vaultClient.Logical().Write(p.caGenPath(), data)
	if err != nil {
		return fmt.Errorf("error writing new CA: %v", err)
	}

	return nil
}

func (p *PKI) caPathExists() (bool, error) {
	path := filepath.Join(p.Path(), "cert", "ca")

	s, err := p.kubernetes.vaultClient.Logical().Read(path)
	if err != nil {
		return false, fmt.Errorf("error reading ca path '%s': %v", path, err)
	}

	if s == nil {
		return false, nil
	}
	if val, ok := s.Data["certificate"]; !ok || val == "" {
		return false, nil
	}

	return true, nil
}

func (p *PKI) WriteRole(role *pkiRole) error {
	path := filepath.Join(p.Path(), "roles", role.Name)

	_, err := p.kubernetes.vaultClient.Logical().Write(path, role.Data)
	if err != nil {
		return fmt.Errorf("error writting role '%s' to '%s': %v", role.Name, p.Path(), err)
	}

	return nil
}

func (p *PKI) ReadRole(role *pkiRole) (*vault.Secret, error) {
	path := filepath.Join(p.Path(), "roles", role.Name)

	secret, err := p.kubernetes.vaultClient.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("error reading role '%s' to '%s': %v", role.Name, p.Path(), err)
	}

	return secret, nil
}

func (p *PKI) Path() string {
	return filepath.Join(p.kubernetes.Path(), pkiType, p.pkiName)
}

func (p *PKI) getMountConfigInput() vault.MountConfigInput {
	return vault.MountConfigInput{
		DefaultLeaseTTL: p.getDefaultLeaseTTL(),
		MaxLeaseTTL:     p.getMaxLeaseTTL(),
	}
}

func (p *PKI) getDefaultLeaseTTL() string {
	return fmt.Sprintf("%ds", int(p.DefaultLeaseTTL.Seconds()))
}

func (p *PKI) getMaxLeaseTTL() string {
	return fmt.Sprintf("%ds", int(p.MaxLeaseTTL.Seconds()))
}

func (p *PKI) getTokenPolicyExists(name string) (bool, error) {
	policy, err := p.kubernetes.vaultClient.Sys().GetPolicy(name)
	if err != nil {
		return false, err
	}

	if policy == "" {
		p.Log.Debugf("Policy Not Found: %s", name)
		return false, nil
	}

	p.Log.Debugf("Policy Found: %s", name)

	return true, nil
}

func (p *PKI) caGenPath() string {
	return filepath.Join(p.Path(), "root", "generate", "internal")
}

func (p *PKI) Type() string {
	return pkiType
}

func (p *PKI) Name() string {
	return p.pkiName
}
