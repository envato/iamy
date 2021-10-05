package iamy

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

const cloudformationStackNameTag = "aws:cloudformation:stack-name"
const includeTestTag = "iamy-include"
const testSkipPathPrefix = "/aws-reserved/"

func TestIsSkippableManagedResource(t *testing.T) {
	skippables := []string{
		"myalias-123/iam/role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot.yaml",
		"AWSServiceRoleTest",
		"my-example-role-ABCDEFGH1234567",
	}

	nonSkippables := []string{
		"myalias-123/iam/user/foo/billy.blogs.yaml",
		"myalias-123/s3/my-bucket.yaml",
		"myalias-123/iam/instance-profile/example.yaml",
	}

	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{cloudformationStackNameTag}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{testSkipPathPrefix}}

	for _, name := range skippables {
		t.Run(name, func(t *testing.T) {

			skipped, err := f.isSkippableManagedResource(CfnIamRole, name, map[string]string{}, "/")
			if skipped == false {
				t.Errorf("expected %s to be skipped but got false", name)
			}

			if err == "" {
				t.Errorf("expected %s to output an error message but it was empty", name)
			}
		})
	}

	for _, name := range nonSkippables {
		t.Run(name, func(t *testing.T) {

			skipped, err := f.isSkippableManagedResource(CfnIamRole, name, map[string]string{}, "/")
			if skipped == true {
				t.Errorf("expected %s to not be skipped but got true", name)
			}

			if err != "" {
				t.Errorf("expected %s to not output an error message but got: %s", name, err)
			}
		})
	}
	for _, name := range nonSkippables {
		t.Run(name, func(t *testing.T) {

			skipped, err := f.isSkippableManagedResource(CfnIamRole, name, map[string]string{}, testSkipPathPrefix)
			if skipped == false {
				t.Errorf("expected %s to be skipped due to path but got false", name)
			}

			if err == "" {
				t.Errorf("expected %s to output an error message but it was empty", name)
			}
		})
	}
}

func TestSkippableS3TaggedResources(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{cloudformationStackNameTag}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	skippableTags := map[string]string{cloudformationStackNameTag: "my-stack"}

	skipped, err := f.isSkippableManagedResource(CfnS3Bucket, "my-bucket", skippableTags, "NOSKIP")
	if err == "" {
		t.Errorf("expected an error message but it was empty")
	}
	if skipped == false {
		t.Errorf("expected resource to be skipped but got false")
	}
}

func TestSkippableS3TaggedResources_WithNoSkipTags(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	skippableTags := map[string]string{cloudformationStackNameTag: "my-stack"}

	skipped, err := f.isSkippableManagedResource(CfnS3Bucket, "my-bucket", skippableTags, "NOSKIP")
	if err != "" {
		t.Errorf("expected no error message but it was " + err)
	}
	if skipped == true {
		t.Errorf("expected resource to not be skipped but got true")
	}
}

func TestNonSkippableTaggedResources(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{cloudformationStackNameTag}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	nonSkippableTags := map[string]string{"Name": "blah"}

	skipped, err := f.isSkippableManagedResource(CfnS3Bucket, "my-bucket", nonSkippableTags, "NOSKIP")
	if err != "" {
		t.Errorf("expected no error message but got: %s", err)
	}
	if skipped == true {
		t.Errorf("expected resource to not be skipped but got true")
	}
}

func TestIncludeTagsOverrideSkip(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{cloudformationStackNameTag}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	TestTags := map[string]string{includeTestTag: "true"}

	skipped, err := f.isSkippableManagedResource(CfnS3Bucket, "my-bucket", TestTags, "NOSKIP")
	if err != "" {
		t.Errorf("expected no error message but got: %s", err)
	}
	if skipped == true {
		t.Errorf("expected resource to not be skipped but got true")
	}

	skipped, err = f.isSkippableManagedResource(CfnS3Bucket, "cfn-bucket-ABCDEF123456", TestTags, "NOSKIP")
	if err != "" {
		t.Errorf("expected no error message but got: %s", err)
	}
	if skipped == true {
		t.Errorf("expected resource to not be skipped but got true")
	}
}

func TestSkippableIAMUserResource(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{cloudformationStackNameTag}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	key := cloudformationStackNameTag
	val := "my-stack"
	userName := "my-user"
	path := "/"
	userList := []*iam.UserDetail{
		{Tags: []*iam.Tag{{Key: &key, Value: &val}}, UserName: &userName, Path: &path},
	}

	resp := iam.GetAccountAuthorizationDetailsOutput{UserDetailList: userList}
	f.populateIamData(&resp)
	for _, user := range f.data.Users {
		if user.Name == userName {
			t.Error("Expected to skip user with CFN tags")
		}
	}
}

func TestSkippableIAMUserResource_WithNoSkipTags(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	key := cloudformationStackNameTag
	val := "my-stack"
	userName := "my-user"
	path := "/"
	userList := []*iam.UserDetail{
		{Tags: []*iam.Tag{{Key: &key, Value: &val}}, UserName: &userName, Path: &path},
	}

	resp := iam.GetAccountAuthorizationDetailsOutput{UserDetailList: userList}
	f.populateIamData(&resp)
	foundUser := false
	for _, user := range f.data.Users {
		if user.Name == userName {
			foundUser = true
		}
	}

	if !foundUser {
		t.Error("Expected to not skip user with CFN tags when SkipTagged: []string{}")
	}
}

func TestSkippableIAMRoleResource(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{cloudformationStackNameTag}, IncludeTagged: []string{includeTestTag}, SkipPathPrefixes: []string{}}
	key := cloudformationStackNameTag
	val := "my-stack"
	roleName := "my-role"
	path := "/"
	roleList := []*iam.RoleDetail{
		{Tags: []*iam.Tag{{Key: &key, Value: &val}}, RoleName: &roleName, Path: &path},
	}

	resp := iam.GetAccountAuthorizationDetailsOutput{RoleDetailList: roleList}
	f.populateIamData(&resp)
	for _, role := range f.data.Roles {
		if role.Name == roleName {
			t.Error("Expected to skip role with CFN tags")
		}
	}
}

func TestSkippableIAMRoleResource_WithNoSkipTags(t *testing.T) {
	f := AwsFetcher{cfn: &cfnClient{}, SkipTagged: []string{}, IncludeTagged: []string{includeTestTag}, SkipFetchingPolicyAndRoleDescriptions: true, SkipPathPrefixes: []string{}}
	key := cloudformationStackNameTag
	val := "my-stack"
	roleName := "my-role"
	path := "/"
	str := "{}"
	roleList := []*iam.RoleDetail{
		{Tags: []*iam.Tag{{Key: &key, Value: &val}}, RoleName: &roleName, Path: &path, AssumeRolePolicyDocument: &str},
	}

	resp := iam.GetAccountAuthorizationDetailsOutput{RoleDetailList: roleList}
	f.populateIamData(&resp)
	foundRole := false
	for _, role := range f.data.Roles {
		if role.Name == roleName {
			foundRole = true
		}
	}
	if !foundRole {
		t.Error("Expected to not skip role with CFN tags and SkipTagged: []string{}")
	}
}
