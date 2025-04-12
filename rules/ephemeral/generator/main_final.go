//go:build generators

package main

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
)

type writeOnlyArgument struct {
	OriginalAttribute    string
	WriteOnlyAlternative string
}

func main() {
	awsProvider := utils.LoadProviderSchema("../../tools/provider-schema/schema.json")

	resourcesWithWriteOnly := map[string][]writeOnlyArgument{}
	// Iterate over all resources in the AWS provider schema
	for resourceName, resource := range awsProvider.ResourceSchemas {
		if arguments := writeOnlyArguments(resource); len(arguments) > 0 {
			// gather sensitive attributes with write only argument alternatives
			resourcesWithWriteOnly[resourceName] = findReplaceableAttribute(arguments, resource)
		}
	}

	// Generate the write-only rule file
	utils.GenerateFile("../../rules/ephemeral/aws_write_only_arguments.go", "../../rules/ephemeral/aws_write_only_arguments_rule.go.tmpl", resourcesWithWriteOnly)

	// Generate the write-only test file
	utils.GenerateFile("../../rules/ephemeral/aws_write_only_arguments_test.go", "../../rules/ephemeral/aws_write_only_arguments_rule_test.go.tmpl", resourcesWithWriteOnly)

	ephemeralResourcesAsDataAlternative := []string{}
	// Iterate over all ephemeral resources in the AWS provider schema
	for resourceName, _ := range awsProvider.EphemeralResourceSchemas {
		if awsProvider.DataSourceSchemas[resourceName] != nil {
			ephemeralResourcesAsDataAlternative = append(ephemeralResourcesAsDataAlternative, resourceName)
		}
	}

	// Generate the ephemeral resource rule file
	utils.GenerateFile("../../rules/ephemeral/aws_ephemeral_resources.go", "../../rules/ephemeral/aws_ephemeral_resources_rule.go.tmpl", ephemeralResourcesAsDataAlternative)

	// Generate the ephemeral resource test file
	utils.GenerateFile("../../rules/ephemeral/aws_ephemeral_resources_test.go", "../../rules/ephemeral/aws_ephemeral_resources_rule_test.go.tmpl", ephemeralResourcesAsDataAlternative)
}

func findReplaceableAttribute(arguments []string, resource *tfjson.Schema) []writeOnlyArgument {
	writeOnlyArguments := []writeOnlyArgument{}

	for _, argument := range arguments {
		// Check if the argument ends with "_wo" and if the original attribute without "_wo" suffix exists in the resource schema
		if attribute := strings.TrimSuffix(argument, "_wo"); strings.HasSuffix(argument, "_wo") && resource.Block.Attributes[attribute] != nil {
			writeOnlyArguments = append(writeOnlyArguments, writeOnlyArgument{
				OriginalAttribute:    attribute,
				WriteOnlyAlternative: argument,
			})
		}
	}

	return writeOnlyArguments
}

func writeOnlyArguments(resource *tfjson.Schema) []string {
	if resource == nil || resource.Block == nil {
		return []string{}
	}

	writeOnlyArguments := []string{}

	// Check if the resource has any write-only attributes
	for name, attribute := range resource.Block.Attributes {
		if attribute.WriteOnly {
			writeOnlyArguments = append(writeOnlyArguments, name)
		}
	}

	return writeOnlyArguments
}
