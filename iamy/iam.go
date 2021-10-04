package iamy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

type iamClient struct {
	iamiface.IAMAPI
}

func newIamClient(sess *session.Session) *iamClient {
	return &iamClient{
		iam.New(sess),
	}
}

func (c *iamClient) getPolicyDescription(arn string) (string, error) {
	resp, err := c.GetPolicy(&iam.GetPolicyInput{PolicyArn: &arn})
	if err == nil && resp.Policy.Description != nil {
		return *resp.Policy.Description, nil
	}
	return "", err
}

func (c *iamClient) getRole(name string) (string, int, error) {
	resp, err := c.GetRole(&iam.GetRoleInput{RoleName: &name})
	var sessionDuration int64
	var description string
	// 3600 is the default, so let's ignore it
	if resp.Role.MaxSessionDuration != nil && *resp.Role.MaxSessionDuration != 3600 {
		sessionDuration = *resp.Role.MaxSessionDuration
	}

	if resp.Role.Description != nil {
		description = *resp.Role.Description
	}
	return description, int(sessionDuration), err
}

func (c *iamClient) MustGetSecurityCredsForUser(username string) (accessKeyIds, mfaIds []string, hasLoginProfile bool) {
	// access keys
	listUsersResp, err := c.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(username),
	})
	if err != nil {
		panic(err)
	}
	for _, m := range listUsersResp.AccessKeyMetadata {
		accessKeyIds = append(accessKeyIds, *m.AccessKeyId)
	}

	// mfa devices
	mfaResp, err := c.ListMFADevices(&iam.ListMFADevicesInput{
		UserName: aws.String(username),
	})
	if err != nil {
		panic(err)
	}
	for _, m := range mfaResp.MFADevices {
		mfaIds = append(mfaIds, *m.SerialNumber)
	}

	// login profile
	_, err = c.GetLoginProfile(&iam.GetLoginProfileInput{
		UserName: aws.String(username),
	})
	if err == nil {
		hasLoginProfile = true
	}

	return
}
