package modules

import (
	"fmt"
	"strings"

	awsarn "github.com/aws/aws-sdk-go/aws/arn"
)

type AWSResourceDef struct {
	API       string `json:"api"`
	Resource  string `json:"resource"`
	Region    string `json:"region"`
	AccountID string `json:"accountId"`
}

type awsResourceAPI interface {
	getResourceDetails(arn string)
	deleteResource(arn string)
}

func parseArn(arn string) (*AWSResourceDef, error) {
	a, err := awsarn.Parse(arn)
	if err != nil {
		return nil, err
	}
	if a.Region == "" {
		a.Region = "global"
	}
	if strings.Contains(a.Resource, "/") {
		rParts := strings.Split(a.Resource, "/")
		a.Service = fmt.Sprintf("%s-%s", a.Service, rParts[0])
		a.Resource = strings.Join(rParts[1:], "")
	}

	ara := &AWSResourceDef{
		API:       a.Service,
		Resource:  a.Resource,
		Region:    a.Region,
		AccountID: a.AccountID,
	}

	return ara, nil
}
