# vault-secrets-shuffle

the application fetches the nodes definition from PuppetDB, generate a random secret and store them to Vault. 
It is meant to be used in conjunction with [hiera_vault](https://github.com/petems/petems-hiera_vault) 

#WIP: please wait


You need to create a configuration file as following (beware of file permissions: it contains Vault token): 

```
[vault_params]
# Vault parameters
vault_token = xxxxxxxxxxx
vault_host = vault.yourdomain.org
vault_path = test/toast
```

and you can run the tool:
```
vault-secrets-shuffle --config /path/to/file.conf
```
