# vault-secrets-shuffle

Fetches the nodes definition from PuppetDB, generate random secrets different for each host and store them to Vault.

It is meant to be used in conjunction with [hiera_vault](https://github.com/petems/petems-hiera_vault)

## configuration

you have:

- a kv v2 store on Vault
- puppet Hiera connected to Vault through [hiera_vault](https://github.com/petems/petems-hiera_vault) and your lookups include certnames/fqdn
- a configuration file with one `vault` section as following (beware of file permissions):

```ini
[vault]
# Vault parameters
vault_token = xxxxxxxxxxx
vault_ssl = true
vault_host = vault.yourdomain.org
vault_port = 443
vault_path = test/toast
vault_keyname = vault_root_password

# Password properties
pass_lenght = 10
min_digits = 2
max_digits = 6
min_symbols = 0
max_symbols = 0

# PuppetDB parameters
puppetdb_host = puppetdb.yourdomain.org
puppetdb_port = 8080
```

## usage

you can run the tool with `--help` to check all options, or you can just run:

```bash
vault-secrets-shuffle --config /path/to/file.conf
```

## compatibility

tested against:

- puppetdb 6.2
- vault 1.0.2

## build

you can use `build.sh` from this repo

## doubts/stoppers

not tested with self-signed certificate

while we could use a KV v1, this is not very safe with a bulk action.

Hence we depends on the following issues to use the KV v2 (which means that you can store keys on a v2 Vault backend, but you can't use them effectively from hiera):

- [petems/petems-hiera_vault#23](https://github.com/petems/petems-hiera_vault/issues/23)
- [hashicorp/vault-ruby#194](https://github.com/hashicorp/vault-ruby/issues/194)
- [hashicorp/vault-ruby#195](https://github.com/hashicorp/vault-ruby/issues/195)
- [hashicorp/vault-ruby#196](https://github.com/hashicorp/vault-ruby/issues/196)
