## ize exec

Execute command in ECS container

### Synopsis

Connect to a container in the ECS via AWS SSM and run command.
It uses app name as an argument.

```
ize exec [app-name] -- [commands] [flags]
```

### Examples

```
  # Connect to a container in the ECS via AWS SSM and run command.
  ize exec goblin ps aux
```

### Options

```
      --container-name string   set container name
      --ecs-cluster string      set ECS cluster name
  -h, --help                    help for exec
      --task string             set task id
```

### Options inherited from parent commands

```
  -p, --aws-profile string         (required) set AWS profile (overrides value in ize.toml and IZE_AWS_PROFILE / AWS_PROFILE if any of them are set)
  -r, --aws-region string          (required) set AWS region (overrides value in ize.toml and IZE_AWS_REGION / AWS_REGION if any of them are set)
  -c, --config-file string         set config file name
  -e, --env string                 (required) set environment name (overrides value set in IZE_ENV / ENV if any of them are set)
  -l, --log-level string           set log level. Possible levels: info, debug, trace, panic, warn, error, fatal(default)
  -n, --namespace string           (required) set namespace (overrides value in ize.toml and IZE_NAMESPACE / NAMESPACE if any of them are set)
      --plain-text                 enable plain text
      --prefer-runtime string      set prefer runtime (native or docker) (default "native")
      --terraform-version string   set terraform-version
```

### SEE ALSO

* [ize](ize.md)	 - 

###### Auto generated by spf13/cobra on 22-Sep-2022