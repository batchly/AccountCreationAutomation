package stack

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func CreateStack(svc *cloudformation.CloudFormation) (resp *cloudformation.CreateStackOutput, err error) {
	params := &cloudformation.CreateStackInput{
		StackName: aws.String("Batchly-IAM-Permissions"), // Required
		Capabilities: []*string{
			aws.String("CAPABILITY_NAMED_IAM"), // Required
			// More values...
		},
		DisableRollback: aws.Bool(false),

		TemplateURL:      aws.String("https://s3.amazonaws.com/batchly-iam/batchly-iam-cloudformation.template"),
		TimeoutInMinutes: aws.Int64(5),
	}
	resp, err = svc.CreateStack(params)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
