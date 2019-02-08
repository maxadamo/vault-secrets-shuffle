# vault-secrets-shuffle

the application fetches the nodes definition from PuppetDB, generate a random secret and store them to Vault.
It is meant to be used in conjunction with [hiera_vault](https://github.com/petems/petems-hiera_vault)

# WIP: please wait

You need to create a configuration file as following (beware of file permissions: it contains the token):

```ini
[vault_params]
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

and you can run the tool (right now, it will only print hosts and passwords):

```bash
vault-secrets-shuffle --config /path/to/file.conf
```
