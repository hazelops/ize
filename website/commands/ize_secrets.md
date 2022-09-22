## ize secrets

Manage secrets

### Options

```
  -h, --help   help for secrets
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
* [ize secrets edit](ize_secrets_edit.md)	 - Edit secrets file
* [ize secrets pull](ize_secrets_pull.md)	 - Pull secrets to a a local file (like SSM)
* [ize secrets push](ize_secrets_push.md)	 - Push secrets to a key-value storage (like SSM)
* [ize secrets rm](ize_secrets_rm.md)	 - Remove secrets from storage

###### Auto generated by spf13/cobra on 22-Sep-2022