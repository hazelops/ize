package schema

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/xeipuuv/gojsonschema"
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
		return fmt.Errorf("%s in %s", m, i)
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
