package schema

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/exp/slices"
	"strings"

	// Enable support for embedded static resources
	_ "embed"
)

//go:embed ize-spec.json
var Schema string

// Validate uses the jsonschema to validate the configuration
func Validate(config map[string]interface{}) error {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	if err := compiler.AddResource("schema.json", strings.NewReader(Schema)); err != nil {
		panic(err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		panic(err)
	}
	err = schema.ValidateInterface(config)
	if err != nil {
		i, m := GetErrorMessage(err.(*jsonschema.ValidationError))
		if i == "#" {
			i = "root"
		} else {
			i = strings.ReplaceAll(i[2:], "/", ".")
		}
		errMsg := fmt.Sprintf("%s in %s of config file (or environment variables)", m, i)
		if strings.Contains(errMsg, "additionalProperties") {
			errMsg += ". The following options are available:\n"
			properties := GetSchema()
			if i != "root" {
				properties = properties[strings.Split(i, ".")[0]].Items
			}
			for k := range properties {
				errMsg += fmt.Sprintf("- %s\n", k)
			}
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}

func GetErrorMessage(err *jsonschema.ValidationError) (string, string) {
	if len(err.Causes) == 0 {
		return err.InstancePtr, err.Message
	}
	return GetErrorMessage(err.Causes[0])
}

func GetJsonSchema() interface{} {
	schemaLoader := gojsonschema.NewStringLoader(Schema)

	json, err := schemaLoader.LoadJSON()
	if err != nil {
		return err
	}

	return json
}

func GetSchema() Items {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	if err := compiler.AddResource("schema.json", strings.NewReader(Schema)); err != nil {
		panic(err)
	}
	compiler.ExtractAnnotations = true
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		panic(err)
	}

	items := Items{}
	getProperties(items, schema)

	return items
}

func getProperties(items Items, schema *jsonschema.Schema) {
	for k, v := range schema.Properties {
		if strings.Contains(v.Description, "deprecated") {
			continue
		}
		if !slices.Contains(v.Types, "object") {
			r := slices.Contains(schema.Required, k)
			items[k] = Item{
				Default:     v.Default,
				Required:    r,
				Description: v.Description,
			}
		} else {
			r := slices.Contains(schema.Required, k)
			i := Items{}
			if len(v.PatternProperties) == 0 {
				getProperties(i, v)
			} else {
				for _, p := range v.PatternProperties {
					getProperties(i, p.Ref)
				}
			}
			items[k] = Item{
				Default:     v.Default,
				Required:    r,
				Description: v.Description,
				Items:       i,
			}
		}
	}
}

type Item struct {
	Default     interface{}
	Required    bool
	Description string
	Items       Items
}

type Items map[string]Item
