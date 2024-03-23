package vault

import (
	"errors"

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
		// vault is the exported type
		vault *Vault
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

// NewModuleInstance implements the modules.Module interface returning a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:    vu,
		vault: &Vault{vu: vu},
	}
}

// Vault is the type for our custom API.
type Vault struct {
	vu     modules.VU // provides methods for accessing internal k6 objects
	client *vault.Client
}

// Vault setup client
func (v *Vault) Setup(address string) {
	var err error
	runtime := v.vu.Runtime()
	v.client, err = vault.New(
		vault.WithAddress(address),
	)
	if err != nil {
		common.Throw(runtime, err)
	}
}

var ErrUninitializedClient = errors.New("Vault client has not been initialized, use setup(address)")

// Vault set Token
func (v *Vault) SetToken(token string) {
	var err error
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	if err = v.client.SetToken(token); err != nil {
		common.Throw(runtime, err)
	}
}

// Vault AppRole Login
func (v *Vault) AppRoleLogin(roleid, secretid, mount string) {
	ctx := v.vu.Context()
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	r, err := v.client.Auth.AppRoleLogin(
		ctx,
		schema.AppRoleLoginRequest{
			RoleId:   roleid,
			SecretId: secretid,
		},
		vault.WithMountPath(mount),
	)
	if err != nil {
		common.Throw(runtime, err)
	}
	if err = v.client.SetToken(r.Auth.ClientToken); err != nil {
		common.Throw(runtime, err)
	}
}

// Vault Kubernetes Login
func (v *Vault) KubernetesLogin(jwt, role, mount string) {
	ctx := v.vu.Context()
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	r, err := v.client.Auth.KubernetesLogin(
		ctx,
		schema.KubernetesLoginRequest{
			Jwt:  jwt,
			Role: role,
		},
		vault.WithMountPath(mount),
	)
	if err != nil {
		common.Throw(runtime, err)
	}
	if err = v.client.SetToken(r.Auth.ClientToken); err != nil {
		common.Throw(runtime, err)
	}
}

// Vault Read
func (v *Vault) Read(path string) interface{} {
	ctx := v.vu.Context()
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	r, err := v.client.Read(ctx, path)
	if err != nil {
		common.Throw(runtime, err)
	}
	return r.Data
}

// Vault List
func (v *Vault) List(path string) interface{} {
	ctx := v.vu.Context()
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	r, err := v.client.List(ctx, path)
	if err != nil {
		common.Throw(runtime, err)
	}
	return r.Data
}

// Vault Write
func (v *Vault) Write(path string, body map[string]interface{}) interface{} {
	ctx := v.vu.Context()
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	r, err := v.client.Write(ctx, path, body)
	if err != nil {
		common.Throw(runtime, err)
	}
	return r.Data
}

// Vault Delete
func (v *Vault) Delete(path string) interface{} {
	ctx := v.vu.Context()
	runtime := v.vu.Runtime()
	if v.client == nil {
		common.Throw(runtime, ErrUninitializedClient)
	}
	r, err := v.client.Delete(ctx, path)
	if err != nil {
		common.Throw(runtime, err)
	}
	return r.Data
}

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.vault,
	}
}
