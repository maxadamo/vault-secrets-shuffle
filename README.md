# vault-secrets-shuffle

Fetches the nodes definition from PuppetDB, generate a random secret and store them to Vault.

It is meant to be used in conjunction with [hiera_vault](https://github.com/petems/petems-hiera_vault)

# WIP: please wait

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

you can run the tool as following:

```bash
vault-secrets-shuffle --config /path/to/file.conf
```

## compatibility

tested against:

- puppetdb 6.2
- vault 1.0.2
