# IZE Simple Example

### Set config path via ENV
```shell
export IZE_CONFIG_FILE=<path to your infra folder>/.infra/env/dev/ize.toml 
```

### Commands for generate terraform files

```shell
ize env terraform
```

### Commands to deploy/destroy infrastructure

```shell
ize deploy infra
ize destroy infra
```

### Commands to deploy/destroy "goblin" project
```shell
ize deploy goblin
ize destroy goblin
```

### Establish SSM tunnel
__Note:__ ssh key at `~/.ssh/id_rsa` should be created before establishing tunnel or please use the `--ssh-private-key` flag
```shell
ize tunnel up
ize tunnel down
```

### Upload/Remove secrets
```shell
ize secret set --file .infra/env/testnut/secrets/example-service.json --type ssm
ize secret remove --type ssm --path /testnut/example-service
```
