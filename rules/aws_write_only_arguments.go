package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
)

// AwsWriteOnlyArgumentsRule checks if a write-only argument is available for sensitive input attributes
type AwsWriteOnlyArgumentsRule struct {
	tflint.DefaultRule

	writeOnlyArguments map[string][]writeOnlyArgument
}

type writeOnlyArgument struct {
	originalAttribute    string
	writeOnlyAlternative string
}

// NewAwsWriteOnlyArgumentsRule returns new rule with default attributes
func NewAwsWriteOnlyArgumentsRule() *AwsWriteOnlyArgumentsRule {
	writeOnlyArguments := map[string][]writeOnlyArgument{
		"aws_db_instance": {
			{
				originalAttribute:    "password",
				writeOnlyAlternative: "password_wo",
			},
		},
		"aws_docdb_cluster": {
			{
				originalAttribute:    "master_password",
				writeOnlyAlternative: "master_password_wo",
			},
		},
		"aws_rds_cluster": {
			{
				originalAttribute:    "master_password",
				writeOnlyAlternative: "master_password_wo",
			},
		},
		"aws_redshift_cluster": {
			{
				originalAttribute:    "master_password",
				writeOnlyAlternative: "master_password_wo",
			},
		},
		"aws_redshiftserverless_namespace": {
			{
				originalAttribute:    "admin_user_password",
				writeOnlyAlternative: "admin_user_password_wo",
			},
		},
		"aws_secretsmanager_secret_version": {
			{
				originalAttribute:    "secret_string",
				writeOnlyAlternative: "secret_string_wo",
			},
		},
		"aws_ssm_parameter": {
			{
				originalAttribute:    "value",
				writeOnlyAlternative: "value_wo",
			},
		},
	}
	return &AwsWriteOnlyArgumentsRule{
		writeOnlyArguments: writeOnlyArguments,
	}
}

// Name returns the rule name
func (r *AwsWriteOnlyArgumentsRule) Name() string {
	return "aws_write_only_arguments"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsWriteOnlyArgumentsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsWriteOnlyArgumentsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsWriteOnlyArgumentsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the sensitive attribute exists
func (r *AwsWriteOnlyArgumentsRule) Check(runner tflint.Runner) error {
	for resourceType, attributes := range r.writeOnlyArguments {
		for _, resourceAttribute := range attributes {
			resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: resourceAttribute.originalAttribute},
				},
			}, nil)
			if err != nil {
				return err
			}

			for _, resource := range resources.Blocks {
				attribute, exists := resource.Body.Attributes[resourceAttribute.originalAttribute]
				if !exists {
					continue
				}

				err := runner.EvaluateExpr(attribute.Expr, func(val cty.Value) error {
					if !val.IsNull() {
						if err := runner.EmitIssueWithFix(
							r,
							fmt.Sprintf("\"%s\" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument \"%s\".", resourceAttribute.originalAttribute, resourceAttribute.writeOnlyAlternative),
							attribute.Expr.Range(),
							func(f tflint.Fixer) error {
								return f.ReplaceText(attribute.NameRange, resourceAttribute.writeOnlyAlternative)
							},
						); err != nil {
							return fmt.Errorf("failed to call EmitIssueWithFix(): %w", err)
						}
					}
					return nil
				}, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
