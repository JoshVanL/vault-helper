package kubernetes

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	vault "github.com/hashicorp/vault/api"
)

type Backend interface {
	Ensure() error
	Path() string
}

type Kubernetes struct {
	clusterID string // clusterID is required parameter, lowercase only, [a-z0-9-]+

	etcdKubernetesPKI *PKI
	etcdOverlayPKI    *PKI
	kubernetesPKI     *PKI
	secretsGeneric    *Generic
}

var _ Backend = &PKI{}
var _ Backend = &Generic{}

func New(clusterID string) *Kubernetes {

	// TODO: validate clusterID

	k := &Kubernetes{
		clusterID: clusterID,
	}

	vaultClient, err := vault.NewClient(nil)
	if err != nil {
		errors.New("Unable to create vault client")
		return nil
	}

	k.etcdKubernetesPKI = NewPKI(k, "etcd-k8s", vaultClient)
	k.etcdOverlayPKI = NewPKI(k, "etcd-overlay", vaultClient)
	k.kubernetesPKI = NewPKI(k, "k8s", vaultClient)

	k.secretsGeneric = NewGeneric(k)

	return k
}

func (k *Kubernetes) backends() []Backend {
	return []Backend{
		k.etcdKubernetesPKI,
		k.etcdOverlayPKI,
		k.kubernetesPKI,
	}
}

func (k *Kubernetes) Ensure() error {
	var result error
	for _, backend := range k.backends() {
		if err := backend.Ensure(); err != nil {
			result = multierror.Append(result, fmt.Errorf("backend %s: %s", backend.Path(), err))
		}
	}
	return result
}

func (k *Kubernetes) Path() string {
	return k.clusterID
}

func NewGeneric(k *Kubernetes) *Generic {
	return &Generic{kubernetes: k}
}

func NewPKI(k *Kubernetes, pkiName string, vaultClient *vault.Client) *PKI {
	return &PKI{
		pkiName:     pkiName,
		kubernetes:  k,
		vaultClient: vaultClient,
	}
}

type PKI struct {
	pkiName     string
	kubernetes  *Kubernetes
	vaultClient *vault.Client
}

func (p *PKI) Ensure() error {

	mount, err := GetMountByPath(p.vaultClient, p.Path())
	if err != nil {
		return err
	}

	if mount == nil {
		err := p.vaultClient.Sys().Mount(
			p.Path(),
			&vault.MountInput{
				Description: "Kubernetes " + p.kubernetes.clusterID + "/" + p.pkiName + " CA",
				Type:        "pki",
			},
		)
		if err != nil {
			return fmt.Errorf("error creating mount: %s", err)
		}
		return nil
	} else {
		if mount.Type != "pki" {
			return fmt.Errorf("mount '%s' already existing with wrong type '%s'", p.Path(), mount.Type)
		}
		return fmt.Errorf("mount '%s' already existing", p.Path())
	}

	//tuneMountRequired := false

	//if mount.Config.DefaultLeaseTTL != int(p.DefaultLeaseTTL.Seconds()) {
	//	tuneMountRequired = true
	//}
	//if mount.Config.MaxLeaseTTL != int(p.MaxLeaseTTL.Seconds()) {
	//	tuneMountRequired = true
	//}

	//if tuneMountRequired {
	//	mountConfig := p.getMountConfigInput()
	//	err := p.vaultClient.Sys().TuneMount(p.path, mountConfig)
	//	if err != nil {
	//		return fmt.Errorf("error tuning mount config: %s", err.Error())
	//	}
	//	p.log.Infof("tuned mount config=%+v")
	//}

	//return errors.New("implement me")
}

func (p *PKI) Path() string {
	return filepath.Join(p.kubernetes.Path(), "pki", p.pkiName)
}

type Generic struct {
	kubernetes *Kubernetes
}

func (g *Generic) Ensure() error {
	return errors.New("implement me")
}

func (g *Generic) Path() string {
	return filepath.Join(g.kubernetes.Path(), "generic")
}

func GetMountByPath(vaultClient *vault.Client, mountPath string) (*vault.MountOutput, error) {

	mounts, err := vaultClient.Sys().ListMounts()
	if err != nil {
		return nil, fmt.Errorf("error listing mounts: %s", err)
	}

	var mount *vault.MountOutput
	for key, _ := range mounts {
		if filepath.Clean(key) == filepath.Clean(mountPath) {
			mount = mounts[key]
			break
		}
	}

	return mount, nil
}
