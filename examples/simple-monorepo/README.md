# IZE Simple Example

### Commands for generate terraform files

```shell
ize -c ./examples/simple-monorepo/.infra/env/testnut/ize.toml env terraform
```

### Commands for deploy infrastructure

```shell
ize -c .infra/env/testnut/ize.toml deploy infra
ize -c .infra/env/testnut/ize.toml destroy infra
```

### Set config path via ENV
```shell
export IZE_CONFIG_FILE=<path to your infra folder>/.infra/env/dev/ize.toml 
```

### Establish SSM tunnel
```shell
ize tunnel ssh-key
ize tunnel up
```

### Upload secrets
```shell
ize secret set --file .infra/env/testnut/secrets/example-service.json --type ssm
```

### Remove secrets
```shell
ize secret remove --type ssm --path /testnut/example-service
```

### Remove secrets
```shell
ize secret remove --type ssm --path /testnut/example-service
```

### Commands for desytoy infrastructure
```shell
ize destoy infra
```