package stack

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func DescribeStack(svc *cloudformation.CloudFormation, stackid string) (resp *cloudformation.DescribeStacksOutput, err error) {

	params := &cloudformation.DescribeStacksInput{

		StackName: aws.String(stackid),
	}
	resp, err = svc.DescribeStacks(params)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
