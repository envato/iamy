package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const accountAlias = "myaccount-123"

var isDryRun = false
var testDir = ""

func mockUi(t *testing.T) Ui {
	return Ui{
		Logger: log.New(io.Discard, "", 0),
		Error:  log.New(io.Discard, "", 0),
		Debug:  log.New(io.Discard, "", 0),
		Exit:   func(code int) { t.Errorf("ui.Exit called with status %d", code) },
	}
}

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

	data, err := os.ReadFile("fixtures/resources/iam-policy-before-normalization.yaml")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(policyDir, "TestPolicy.yaml"), data, 0644)
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

func TestNormalization(t *testing.T) {
	expected, err := os.ReadFile("fixtures/resources/iam-policy-after-normalization.yaml")
	if err != nil {
		t.Fatal(err)
	}

	input := NormalizeCommandInput{
		Dir:       testDir,
		CanDelete: true,
	}

	NormalizeCommand(mockUi(t), input)

	actual, err := os.ReadFile(filepath.Join(testDir, accountAlias, "iam", "policy", "TestPolicy.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}
