package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const accountAlias = "myaccount-123"

var isDryRun = false
var testDir = ""

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	err = teardown()
	if err != nil {
		log.Println(err)
	}

	os.Exit(code)
}

func setup() error {
	dryRun = &isDryRun
	tempDir, err := os.MkdirTemp("", "iamy-test-normalize")
	if err != nil {
		return err
	}

	policyDir := filepath.Join(tempDir, accountAlias, "iam", "policy")

	err = os.MkdirAll(policyDir, 0755)
	if err != nil {
		return err
	}

	log.Println(os.Getwd())

	data, err := ioutil.ReadFile("fixtures/resources/iam-policy-before-normalization.yaml")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(policyDir, "TestPolicy.yaml"), data, 0644)
	if err != nil {
		return err
	}

	testDir = tempDir

	return nil
}

func teardown() error {
	if testDir != "" {
		err := os.RemoveAll(testDir)
		return err
	}

	return nil
}

func newUi() Ui {
	return Ui{
		Logger: log.New(ioutil.Discard, "", 0),
		Error:  log.New(ioutil.Discard, "", 0),
		Debug:  log.New(ioutil.Discard, "", 0),
		Exit:   os.Exit, // TODO: probably wrong
	}
}

func TestNormalization(t *testing.T) {
	expected, err := ioutil.ReadFile("fixtures/resources/iam-policy-after-normalization.yaml")
	if err != nil {
		t.Fatal(err)
	}

	input := NormalizeCommandInput{
		Dir:       testDir,
		CanDelete: true,
	}

	NormalizeCommand(newUi(), input)

	actual, err := ioutil.ReadFile(filepath.Join(testDir, accountAlias, "iam", "policy", "TestPolicy.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}
