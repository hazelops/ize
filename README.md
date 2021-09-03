# Ize Tool (WIP)

This tool is designed to be an opinionated infrastructure wrapper that allows to use multiple tools in one infra: terraform, serverless, waypoint. 
It combines build and deploy workflows in one.

This tool is using configuration file that describes the workflows.

## Quickstart
- GO version should be 1.16+
- `GOPATH` environment variable is set to `~/go` 

### Ize initialization
```shell
go mod download
make install
```

(acts as an ideation doc, stuff is not working)
### Application Lifecycle

```shell
ize build <goblin>
ize deploy <goblin>
```

### Operations Lifecycle
#### Establish SSM tunnel
```shell
ize tunnel up
ize tunnel down
```

#### Upload secrets
```shell
ize secret set
ize secret get
```

#### Deploy Infra
```shell
ize deploy infra
ize destroy infra
```
