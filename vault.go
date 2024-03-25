package vault

import (
	"context"
	"os"

	"github.com/dop251/goja"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/hashicorp-vault", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU
	}

	Vault struct {
		client *vault.Client
		ctx    context.Context
	}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu: vu,
	}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"Vault": mi.newClient,
		},
	}
}

func (mi *ModuleInstance) newClient(c goja.ConstructorCall) *goja.Object {
	rt := mi.vu.Runtime()
	ctx := mi.vu.Context()

	v := &Vault{}
	var address string
	err := rt.ExportTo(c.Argument(0), &address)
	if err != nil {
		common.Throw(rt, err)
	}
	v.ctx = ctx
	v.client, err = vault.New(
		vault.WithAddress(address),
	)
	if err != nil {
		common.Throw(rt, err)
	}
	return rt.ToValue(v).ToObject(rt)
}

// Vault set Token
func (v *Vault) SetToken(token string) error {
	if err := v.client.SetToken(token); err != nil {
		return err
	}
	return nil
}

// Vault AppRole Login
func (v *Vault) AppRoleLogin(roleid, secretid, mount string) error {
	if mount == "" {
		mount = "approle"
	}
	r, err := v.client.Auth.AppRoleLogin(
		v.ctx,
		schema.AppRoleLoginRequest{
			RoleId:   roleid,
			SecretId: secretid,
		},
		vault.WithMountPath(mount),
	)
	if err != nil {
		return err
	}
	if err = v.client.SetToken(r.Auth.ClientToken); err != nil {
		return err
	}
	return nil
}

// Vault Kubernetes Login
func (v *Vault) KubernetesLogin(role, jwt, mount string) error {
	if jwt == "" {
		token, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			return err
		}
		jwt = string(token)
	}
	if mount == "" {
		mount = "kubernetes"
	}
	r, err := v.client.Auth.KubernetesLogin(
		v.ctx,
		schema.KubernetesLoginRequest{
			Jwt:  jwt,
			Role: role,
		},
		vault.WithMountPath(mount),
	)
	if err != nil {
		return err
	}
	if err = v.client.SetToken(r.Auth.ClientToken); err != nil {
		return err
	}
	return nil
}

// Vault Read
func (v *Vault) Read(path string) (interface{}, error) {
	r, err := v.client.Read(v.ctx, path)
	if err != nil {
		return false, err
	}
	if r != nil {
		return r.Data, nil
	}
	return nil, nil
}

// Vault List
func (v *Vault) List(path string) (interface{}, error) {
	r, err := v.client.List(v.ctx, path)
	if err != nil {
		return false, err
	}
	if r != nil {
		return r.Data, nil
	}
	return nil, nil
}

// Vault Write
func (v *Vault) Write(path string, body map[string]interface{}) (interface{}, error) {
	r, err := v.client.Write(v.ctx, path, body)
	if err != nil {
		return false, err
	}
	if r != nil {
		return r.Data, nil
	}
	return nil, nil
}

// Vault Delete
func (v *Vault) Delete(path string) (interface{}, error) {
	r, err := v.client.Delete(v.ctx, path)
	if err != nil {
		return false, err
	}
	if r != nil {
		return r.Data, nil
	}
	return nil, nil
}
