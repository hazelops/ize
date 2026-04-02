package ecscron

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type networkConfig struct {
	SecurityGroups struct {
		Value string `json:"value"`
	} `json:"security_groups"`
	Subnets struct {
		Value [][]string `json:"value"`
	} `json:"subnets"`
	VpcPrivateSubnets struct {
		Value []string `json:"value"`
	} `json:"vpc_private_subnets"`
	VpcPublicSubnets struct {
		Value []string `json:"value"`
	} `json:"vpc_public_subnets"`
}

func getNetworkConfigFromSSM(svc ssmiface.SSMAPI, env string) (*networkConfig, error) {
	resp, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("can't get terraform output: %w", err)
	}

	value, err := base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		return nil, fmt.Errorf("can't decode terraform output: %w", err)
	}

	var nc networkConfig
	if err := json.Unmarshal(value, &nc); err != nil {
		return nil, fmt.Errorf("can't unmarshal network configuration: %w", err)
	}
	return &nc, nil
}
