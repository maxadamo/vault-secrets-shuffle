# vault-secrets-shuffle

Fetches nodes definitions from PuppetDB, generate random secrets for each host and store them to Vault.

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

# PuppetDB parameters
puppetdb_host = puppetdb.yourdomain.org
puppetdb_port = 8080

# Password properties
pass_lenght = 10
min_digits = 2
max_digits = 6
min_symbols = 0
max_symbols = 0
```

## usage

you can run the tool with `--help` to check all options:

```bash
vault-secrets-shuffle --help
Vault Secrets Shuffler:
  - iterates all VMs registered in PuppetDB
  - generate generate random secrets different for each host
  - upload the secrets to vault.

Usage:
  vault-secrets-shuffle --config=CONFIG [--kv=kv] [--write=WRITE] [--debug]
  vault-secrets-shuffle -v | --version
  vault-secrets-shuffle -b | --build
  vault-secrets-shuffle -h | --help

Options:
  -h --help           Show this screen
  -c --config=CONFIG  Config file
  -w --write=WRITE    Output file (OPTIONAL)
  -k --kv=kv          Keystore Version. [default: 2]
  -d --debug          Print password and full key path (OPTIONAL)
  -v --version        Print version exit
  -b --build          Print version and build information and exit
```

or you can simply run:

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

Some change is already on the work on [petems/petems-hiera_vault#43](https://github.com/petems/petems-hiera_vault/pull/43)

These changes will allow to use Kv v2, which is safer to use (as it has password history)
