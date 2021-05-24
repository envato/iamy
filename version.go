package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"

	"github.com/blang/semver/v4"
)

func checkVersion() error {
	if Version != "dev" {
		requiredVersion, err := fetchRequiredVersion(versionFileName)
		if err != nil {
			return err
		}
		currentVersion, err := semver.ParseTolerant(Version)
		// Ignore Prepatches - iamy uses them to mark builds from uncommitted trees
		currentVersion.Pre = nil
		if err != nil {
			return err
		}
		ok, msg := versionOk(currentVersion, requiredVersion)
		if !ok {
			return msg
		}
	}

	return nil
}

func fetchRequiredVersion(filename string) (semver.Version, error) {
	if _, err := os.Stat(versionFileName); !os.IsNotExist(err) {
		log.Printf("%s found", filename)
		fileBytes, _ := ioutil.ReadFile(filename)
		fileContents := string(fileBytes)

		if fileContents != "" {
			re := regexp.MustCompile(`\d\.\d+\.\d\+?\w*`)
			match := re.FindStringSubmatch(fileContents)
			return semver.Make(match[0])
		}
	}
	return semver.Parse("0.0.0")
}

func versionOk(current semver.Version, required semver.Version) (bool, error) {
	if current.LT(required) {
		msg := fmt.Errorf(versionTooOldError, current, required)
		return false, msg
	}
	if len(required.Build) > 0 {
		// Pay attention to build tags as well if they are required
		if !reflect.DeepEqual(required.Build, current.Build) {
			msg := fmt.Errorf(buildVersionMismatch, current, required.Build, required)
			return false, msg
		}
	}
	return true, nil
}
