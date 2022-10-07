## ize secrets pull

Pull secrets to a a local file (like SSM)

### Synopsis

This command pulls secrets from a key-value storage to a local file (like SSM)

```
ize secrets pull [flags]
```

### Options

```
      --backend string   backend type (default=ssm) (default "ssm")
      --file string      file with secrets
      --force            allow values overwrite
  -h, --help             help for pull
      --path string      path where to store secrets (/<env>/<app> by default)
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

* [ize secrets](ize_secrets.md)	 - Manage secrets

###### Auto generated by spf13/cobra on 22-Sep-2022