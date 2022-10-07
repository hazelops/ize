## ize



### Synopsis

  [37m[1mWelcome to IZE[0m[37m[0m
  [34mDocs:[0m https://ize.sh/docs
  [32mVersion:[0m 1.1.5 (5a8df5bd)
  
  Opinionated tool for infrastructure and code.
  
  This tool is designed as a simple wrapper around popular tools,
  so they can be easily integrated in one infra: terraform,
  ECS deployment, serverless, and others.
  
  It combines infra, build and deploy workflows in one
  and is too simple to be considered sophisticated.
  So let's not do it but rather embrace the simplicity and minimalism.

### Options

```
  -p, --aws-profile string         (required) set AWS profile (overrides value in ize.toml and IZE_AWS_PROFILE / AWS_PROFILE if any of them are set)
  -r, --aws-region string          (required) set AWS region (overrides value in ize.toml and IZE_AWS_REGION / AWS_REGION if any of them are set)
  -c, --config-file string         set config file name
  -e, --env string                 (required) set environment name (overrides value set in IZE_ENV / ENV if any of them are set)
  -h, --help                       help for ize
  -l, --log-level string           set log level. Possible levels: info, debug, trace, panic, warn, error, fatal(default)
  -n, --namespace string           (required) set namespace (overrides value in ize.toml and IZE_NAMESPACE / NAMESPACE if any of them are set)
      --plain-text                 enable plain text
      --prefer-runtime string      set prefer runtime (native or docker) (default "native")
  -t, --tag string                 set tag
      --terraform-version string   set terraform-version
```

### SEE ALSO

* [ize build](ize_build.md)	 - build apps
* [ize configure](ize_configure.md)	 - Generate global configuration file
* [ize console](ize_console.md)	 - Connect to a container in the ECS
* [ize deploy](ize_deploy.md)	 - Manage deployments
* [ize down](ize_down.md)	 - Destroy application
* [ize exec](ize_exec.md)	 - Execute command in ECS container
* [ize gen](ize_gen.md)	 - Generate something
* [ize init](ize_init.md)	 - Initialize project
* [ize logs](ize_logs.md)	 - Stream logs of container in the ECS
* [ize push](ize_push.md)	 - push app's image
* [ize secrets](ize_secrets.md)	 - Manage secrets
* [ize status](ize_status.md)	 - Show debug information
* [ize terraform](ize_terraform.md)	 - Run terraform
* [ize tunnel](ize_tunnel.md)	 - Tunnel management
* [ize up](ize_up.md)	 - Bring full application up from the bottom to the top.
* [ize validate](ize_validate.md)	 - Validate configuration (only for test)
* [ize version](ize_version.md)	 - Show IZE version

###### Auto generated by spf13/cobra on 22-Sep-2022