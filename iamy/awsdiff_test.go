package iamy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadDataFrom(p string) *AccountData {
	d, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	yamlLoader := YamlLoadDumper{
		Dir: filepath.Join(d, "testdata", "awsdiff", p),
	}

	dd, err := yamlLoader.Load()
	if err != nil {
		panic(err.Error())
	}
	if len(dd) < 1 {
		return &AccountData{}
	}
	return &dd[0]
}

func TestPolicyIsDetachedFromRoleBeforeUpdate(t *testing.T) {
	localData := loadDataFrom("testcase1-local")
	remoteData := loadDataFrom("testcase1-remote")
	awsCmds := AwsCliCmdsForSync(remoteData, localData)

	expected := strings.Join([]string{
		"aws iam detach-role-policy --role-name testrole --policy-arn arn:aws:iam::123:policy/test",
		"aws iam delete-policy --policy-arn arn:aws:iam::123:policy/test",
	}, "\n")
	actual := awsCmds.String()

	if actual != expected {

		t.Errorf(`Expected:
%v
Actual:
%v`, expected, actual)

	}
}

func TestCreateRoleWithMaxSessionDuration(t *testing.T) {
	localData := loadDataFrom("max-session-duration-local")
	remoteData := loadDataFrom("max-session-duration-remote1")
	awsCmds := AwsCliCmdsForSync(remoteData, localData)

	expected := strings.Join([]string{
		`aws iam create-role --role-name testrole --path / --assume-role-policy-document '{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123:root"
      },
      "Sid": ""
    }
  ],
  "Version": "2012-10-17"
}' --max-session-duration 7200`,
	}, "\n")
	actual := awsCmds.String()

	if actual != expected {

		t.Errorf(`Expected:
%v
Actual:
%v`, expected, actual)

	}
}

func TestUpdateRoleWithMaxSessionDuration(t *testing.T) {
	localData := loadDataFrom("max-session-duration-local")
	remoteData := loadDataFrom("max-session-duration-remote2")
	awsCmds := AwsCliCmdsForSync(remoteData, localData)

	expected := strings.Join([]string{
		"aws iam update-role --role-name testrole --max-session-duration 7200",
	}, "\n")
	actual := awsCmds.String()

	if actual != expected {

		t.Errorf(`Expected:
%v
Actual:
%v`, expected, actual)

	}
}
