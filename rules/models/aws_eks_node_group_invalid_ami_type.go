// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEksNodeGroupInvalidAMITypeRule checks the pattern is valid
type AwsEksNodeGroupInvalidAMITypeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsEksNodeGroupInvalidAMITypeRule returns new rule with default attributes
func NewAwsEksNodeGroupInvalidAMITypeRule() *AwsEksNodeGroupInvalidAMITypeRule {
	return &AwsEksNodeGroupInvalidAMITypeRule{
		resourceType:  "aws_eks_node_group",
		attributeName: "ami_type",
		enum: []string{
			"AL2_x86_64",
			"AL2_x86_64_GPU",
			"AL2_ARM_64",
			"CUSTOM",
			"BOTTLEROCKET_ARM_64",
			"BOTTLEROCKET_x86_64",
			"BOTTLEROCKET_ARM_64_FIPS",
			"BOTTLEROCKET_x86_64_FIPS",
			"BOTTLEROCKET_ARM_64_NVIDIA",
			"BOTTLEROCKET_x86_64_NVIDIA",
			"WINDOWS_CORE_2019_x86_64",
			"WINDOWS_FULL_2019_x86_64",
			"WINDOWS_CORE_2022_x86_64",
			"WINDOWS_FULL_2022_x86_64",
			"AL2023_x86_64_STANDARD",
			"AL2023_ARM_64_STANDARD",
			"AL2023_x86_64_NEURON",
			"AL2023_x86_64_NVIDIA",
		},
	}
}

// Name returns the rule name
func (r *AwsEksNodeGroupInvalidAMITypeRule) Name() string {
	return "aws_eks_node_group_invalid_ami_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsEksNodeGroupInvalidAMITypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsEksNodeGroupInvalidAMITypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsEksNodeGroupInvalidAMITypeRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsEksNodeGroupInvalidAMITypeRule) Check(runner tflint.Runner) error {
	logger.Trace("Check `%s` rule", r.Name())

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func (val string) error {
			found := false
			for _, item := range r.enum {
				if item == val {
					found = true
				}
			}
			if !found {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" is an invalid value as ami_type`, truncateLongMessage(val)),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
