package iamy

import "testing"

func TestCfnMangedResources(t *testing.T) {
	t.Run("With fetched CFN resource lists", func(t *testing.T) {

		cfn := cfnClient{
			managedResources: map[string]CfnResourceTypes{
				"foobar": []CfnResourceType{CfnIamPolicy, CfnIamRole},
			},
		}

		if cfn.IsManagedResource(CfnIamUser, "foobar") {
			t.Fatal("different object types with same name is not managed")
		}

		if !cfn.IsManagedResource(CfnIamPolicy, "foobar") {
			t.Fatal("matching object and type should be managed")
		}
	})

	t.Run("With heuristic matching", func(t *testing.T) {
		cfn := cfnClient{}

		if cfn.IsManagedResource(CfnIamUser, "foobar") {
			t.Fatal("names without id suffix are not managed")
		}

		if !cfn.IsManagedResource(CfnIamPolicy, "foobar-ABCDEFGH1234567") {
			t.Fatal("names with id suffix are managed")
		}

		if cfn.IsManagedResource(CfnIamPolicy, "foobar-abcdefgh1234567") {
			t.Fatal("names with suffix containing only lowercase letters are not managed")
		}

		if cfn.IsManagedResource(CfnIamPolicy, "elasticbeanstalk-us-east-1-298865909318") {
			t.Fatal("names ending with account numbers are managed")
		}

		if cfn.IsManagedResource(CfnIamPolicy, "Elasticbeanstalk-us-East-1-298865909318") {
			t.Fatal("names ending with account numbers are managed, even if the name contains uppercase letters")
		}
	})
}
