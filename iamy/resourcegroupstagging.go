package iamy

import (
	"log"

        "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
)

type resourceGroupsTaggingAPIClient struct {
	resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
}

func newResourceGroupsTaggingAPIClient(sess *session.Session) *resourceGroupsTaggingAPIClient {
        // Force us of us-east-1 otherwise tags will not be returned for global resources
	return &resourceGroupsTaggingAPIClient{
		resourcegroupstaggingapi.New(sess, aws.NewConfig().WithRegion("us-east-1")),
	}
}

func (c *resourceGroupsTaggingAPIClient) getMultiplePolicyTags(arns []*string) (map[string]map[string]string, error) {
	queryArns := make([]string, 0)
	for _, s := range arns {
		queryArns = append(queryArns, *s)
	}
	log.Println("Fetching tags for:", queryArns)

	res := make(map[string]map[string]string)
	if len(arns) == 0 {
		return res, nil
	}
	resp, err := c.GetResources(&resourcegroupstaggingapi.GetResourcesInput{ResourceARNList: arns})
	if err != nil {
		return nil, err
	}
	for _, mapping := range resp.ResourceTagMappingList {
		res[*mapping.ResourceARN] = make(map[string]string)
		for _, tag := range mapping.Tags {
			res[*mapping.ResourceARN][*tag.Key] = *tag.Value
		}
	}
	return res, nil
}
