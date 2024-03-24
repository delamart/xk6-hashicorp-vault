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