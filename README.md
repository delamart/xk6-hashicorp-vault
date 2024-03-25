[![Go Reference](https://pkg.go.dev/badge/github.com/delamart/xk6-hashicorp-vault.svg)](https://pkg.go.dev/github.com/delamart/xk6-hashicorp-vault)
[![Version Badge](https://img.shields.io/github/v/release/delamart/xk6-hashicorp-vault?style=flat-square)](https://github.com/delamart/xk6-hashicorp-vault/releases)

# xk6-hashicorp-vault
A k6 extension for interacting with Hashicorp Vault servers while testing. Mostly a wrapper for [vault-client-go](https://github.com/hashicorp/vault-client-go).

## Build

To build a custom `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

1. Download [xk6](https://github.com/grafana/xk6):
  
    ```bash
    go install go.k6.io/xk6/cmd/xk6@latest
    ```

2. [Build the k6 binary](https://github.com/grafana/xk6#command-usage):
  
    ```bash
    xk6 build --with github.com/delamart/xk6-hashicorp-vault
    ```

    The `xk6 build` command creates a k6 binary that includes the xk6-hashicorp-vault extension in your local folder.


### Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```shell
git clone git@github.com:delamart/xk6-hashicorp-vault.git
cd xk6-hashicorp-vault
make
```

Using the `k6` binary with `xk6-hashicorp-vault`, run the k6 test as usual:

```bash
./k6 run k8s-test-script.js

```
# Usage

* Create a new Vault instance which will create a vault client to a specific server url
* Set the authentication using one of the three methods
  * `vault.setToken( token )`
  * `vault.AppRoleLogin( roleid, secretid, mount )` mount defaults to 'approle'
  * `vault.KubernetesLogin( role, jwt, mount )` jwt defaults to the content of `/var/run/secrets/kubernetes.io/serviceaccount/token`, mount defaults to 'kubernetes'
* Use the provided methods
  * `vault.read( path )`
  * `vault.write( path, data )`
  * `vault.list( path )`
  * `vault.delete( path )`

## Examples

### Connect, write, read, list and delete a secret
```javascript
import { Vault } from 'k6/x/hashicorp-vault';

// run a DEV vault server with docker using:
// docker run --rm --cap-add=IPC_LOCK -p 8200:8200 hashicorp/vault server -dev -dev-root-token-id=root
export default function () {
    const vault = new Vault('http://localhost:8200')
    vault.setToken('root')

    var w = vault.write('secret/data/test', {'data': {'key': 'value'}})
    console.log(w.version)

    var r = vault.read('secret/data/test')
    console.log(r.data.key)

    var l = vault.list('secret/metadata')
    console.log(l.keys)

    var d = vault.delete('secret/data/test')
    console.log(d)
}
```