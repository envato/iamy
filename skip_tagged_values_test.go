package main

import "testing"

func TestValidateSkipTaggedValues(t *testing.T) {
	for _, filter := range []string{"DeployedBy=Terraform/Ship-SI", "DeployedBy=Terraform/Ship", "DeployedBy=ship", "Key=", "Key=value=with=equals"} {
		if err := validateSkipTaggedValues([]string{filter}); err != nil {
			t.Errorf("%q: %v", filter, err)
		}
	}
	for _, filter := range []string{"", "DeployedBy", "=ship"} {
		if err := validateSkipTaggedValues([]string{filter}); err == nil {
			t.Errorf("expected %q to be rejected", filter)
		}
	}
}
