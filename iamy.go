package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

const versionTooOldError = "Your version of IAMy (%s) is out of date compared to what the local project expects. You should upgrade to %s to use this project.\n"
const buildVersionMismatch = "Your version of IAMy (%s) does not match the build tag (%s) the local project requires. You should upgrade to %s to use this project.\n"

var (
	Version         string = "dev"
	defaultDir      string
	dryRun          *bool
	versionFileName string = ".iamy-version"
	configFileName  string = ".iamy-flags"
)

type logWriter struct{ *log.Logger }

func (w logWriter) Write(b []byte) (int, error) {
	w.Printf("%s", b)
	return len(b), nil
}

type Ui struct {
	*log.Logger
	Error, Debug *log.Logger
	Exit         func(code int)
}

// CFN automatically tags resources with this and other tags:
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resource-tags.html
const cloudformationStackNameTag = "aws:cloudformation:stack-name"

func main() {
	var (
		debug            = kingpin.Flag("debug", "Show debugging output").Bool()
		skipCfnTagged    = kingpin.Flag("skip-cfn-tagged", fmt.Sprintf("Shorthand for --skip-tagged %s", cloudformationStackNameTag)).Bool()
		skipTagged       = kingpin.Flag("skip-tagged", "Skips IAM entities (or buckets associated with bucket policies) tagged with a given tag").Strings()
		includeTagged    = kingpin.Flag("include-tagged", "Includes IAM entities (or buckets associated with bucket policies) tagged with a given tag").Strings()
		skipPathPrefixes = kingpin.Flag("skip-path-prefix", fmt.Sprintf("Skips IAM entities that have a path starting with the supplied prefix, repeat flag for multiple prefixes")).Strings()
		pull             = kingpin.Command("pull", "Syncs IAM users, groups and policies from the active AWS account to files")
		pullDir          = pull.Flag("dir", "The directory to dump yaml files to").Default(defaultDir).Short('d').String()
		pullCanDelete    = pull.Flag("delete", "Delete extraneous files from destination dir").Bool()
		lookupCfn        = pull.Flag("accurate-cfn", "Fetch all known resource names from cloudformation to get exact filtering").Bool()
		push             = kingpin.Command("push", "Syncs IAM users, groups and policies from files to the active AWS account")
		pushDir          = push.Flag("dir", "The directory to load yaml files from").Default(defaultDir).Short('d').ExistingDir()
		format           = kingpin.Command("fmt", "Update YAML files to match expected format")
		formatDir        = format.Flag("dir", "The base directory to format").Default(defaultDir).Short('d').ExistingDir()
		formatCanDelete  = format.Flag("delete", "Delete extraneous files from destination dir").Bool()
	)
	dryRun = kingpin.Flag("dry-run", "Show what would happen, but don't prompt to do it").Bool()

	kingpin.Version(Version)
	kingpin.CommandLine.Help =
		`Read and write AWS IAM users, policies, groups and roles from YAML files.`

	ui := Ui{
		Logger: log.New(os.Stdout, "", 0),
		Error:  log.New(os.Stderr, "", 0),
		Debug:  log.New(ioutil.Discard, "", 0),
		Exit:   os.Exit,
	}

	args := os.Args[1:]
	var configFileArgs []string
	if _, err := os.Stat(configFileName); err == nil {
		configFileArgs, err = kingpin.ExpandArgsFromFile(configFileName)
		if err != nil {
			panic(err)
		}
		args = append(args, configFileArgs...)
	}
	cmd, err := kingpin.CommandLine.Parse(args)
	if err != nil {
		panic(err)
	}

	if *debug {
		ui.Debug = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
		log.SetFlags(0)
		log.SetOutput(&logWriter{ui.Debug})
		if len(configFileArgs) > 0 {
			ui.Debug.Printf("Found flags in %s: %s", configFileName, configFileArgs)
		}
	} else {
		log.SetOutput(ioutil.Discard)
	}

	if err := checkVersion(); err != nil {
		panic(err)
	}

	if *skipCfnTagged {
		*skipTagged = append(*skipTagged, cloudformationStackNameTag)
	}

	switch cmd {
	case push.FullCommand():
		PushCommand(ui, PushCommandInput{
			Dir:                  *pushDir,
			HeuristicCfnMatching: !*lookupCfn,
			SkipTagged:           *skipTagged,
			IncludeTagged:        *includeTagged,
			SkipPathPrefixes:     *skipPathPrefixes,
		})

	case pull.FullCommand():
		PullCommand(ui, PullCommandInput{
			Dir:                  *pullDir,
			CanDelete:            *pullCanDelete,
			HeuristicCfnMatching: !*lookupCfn,
			SkipTagged:           *skipTagged,
			IncludeTagged:        *includeTagged,
			SkipPathPrefixes:     *skipPathPrefixes,
		})

	case format.FullCommand():
		FormatCommand(ui, FormatCommandInput{
			Dir:       *formatDir,
			CanDelete: *formatCanDelete,
		})
	}
}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		panic(err)
	}
	defaultDir = filepath.Clean(dir)
}
