# IZE Simple Example

### Set config path via ENV
```shell
export IZE_CONFIG_FILE=<path to your infra folder>/.infra/env/dev/ize.toml 
```

### Commands for generate terraform files

```shell
ize gen tfenv
```

### Commands to deploy/destroy infrastructure

```shell
ize up infra
ize down infra
```

### Commands to deploy/destroy "goblin" project
```shell
ize up goblin
ize down goblin
```

### Establish SSM tunnel
__Note:__ ssh key at `~/.ssh/id_rsa` should be created before establishing tunnel or use the `--ssh-private-key` flag
```shell
ize tunnel up
ize tunnel down
```

### Connect to a container in the ECS
```shell
ize console goblin
```

### Upload/Remove secrets
```shell
ize secrets push goblin --backend ssm --file goblin.json --force
ize secrets rm goblin
```

### Show debug information
```shell
ize status
```
