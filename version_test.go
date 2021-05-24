package main

import "testing"

func TestCheckVersionIsHigher(t *testing.T) {
	versionFileName = "fixtures/0.0.1"
	Version = "0.0.2"
	if err := checkVersion(); err != nil {
		t.Errorf("Received an unexpected error %s", err)
	}
}

func TestCheckVersionIsLower(t *testing.T) {
	versionFileName = "fixtures/0.0.2"
	Version = "0.0.1"
	if err := checkVersion(); err == nil {
		t.Error("Received no error", err)
	}
}

func TestCheckVersionIsTheSameWithBuildTag(t *testing.T) {
	versionFileName = "fixtures/0.0.2"
	Version = "0.0.2+envato"
	if err := checkVersion(); err != nil {
		t.Errorf("Received an unexpected error %s", err)
	}
}

func TestCheckVersionIsLowerWithBuildTag(t *testing.T) {
	versionFileName = "fixtures/0.0.2"
	Version = "0.0.1+envato"
	if err := checkVersion(); err == nil {
		t.Error("Received no error", err)
	}
}

func TestCheckVersionIsTheSameWithMatchingBuildTag(t *testing.T) {
	versionFileName = "fixtures/0.0.2.buildtag"
	Version = "0.0.2+buildtag"
	if err := checkVersion(); err != nil {
		t.Errorf("Received an unexpected error %s", err)
	}
}

func TestCheckVersionIsTheSameWithoutMatchingBuildTag(t *testing.T) {
	versionFileName = "fixtures/0.0.2.buildtag"
	Version = "0.0.2+notbuildtag"
	if err := checkVersion(); err == nil {
		t.Error("Received no error", err)
	}
}

func TestCheckVersionWithPreRelease(t *testing.T) {
	versionFileName = "fixtures/0.0.2"
	Version = "0.0.2-dirty"
	if err := checkVersion(); err != nil {
		t.Errorf("Received an unexpected error %s", err)
	}
}
