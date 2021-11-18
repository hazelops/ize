package config

import (
	"github.com/hashicorp/hcl/v2"
)

type hclInfra struct {
	Provider string   `hcl:"provider,label"`
	Name     string   `hcl:"name,label"`
	Body     hcl.Body `hcl:",body"`
	Remain   hcl.Body `hcl:",remain"`
}
