# IAMy

IAMy is a tool for dumping and loading your AWS IAM configuration into YAML files.

This allows you to use an [Infrastructure as Code](https://en.wikipedia.org/wiki/Infrastructure_as_Code) model to manage your IAM configuration. For example, you might use a github repo with a pull request model for changes to IAM config.

This code was originally developed by 99designs ([origin upstream](https://github.com/99designs/iamy.git)), we recognise and appreciate the enormous effort they have put into this tool.
This particular version has been cloned to allow Envato to rapidly develop the features that are important to our use of this tool, we are following the existing semver arrangements for the repository, but we've appended a envato build tag.

# Additional features

Features added to this fork include:
- .iamy-version file support, [Original PR](https://github.com/99designs/iamy/pull/63)
- Flags to skip resources by tag (`--skip-tagged the-tag-name` and `--skip-cfn-tagged`)
- .iamy-flags file support for default flags. Flags are appended to command line supplied flags. Example .iamy-flags file
  contents: `--skip-tagged=iamy-ignore`.
- `iamy fmt`, which formats files to match the result of `iamy pull`
- Add support for specifying [MaxSessionDuration](https://aws.amazon.com/about-aws/whats-new/2018/03/longer-role-sessions/) on a role

# Upcoming features

The additional features we are likely to add to this fork are:
- support for organizations, ous and scps

# Installation

```
brew tap envato/envato-iamy
brew install envato/envato-iamy/iamy
```

# Development Status

Under active development, pull requests welcome.  Open issues for discussions please.

## How it works

IAMy has two main subcommands.

`pull` will sync IAM users, groups and policies from AWS to YAML files

`push` will sync IAM users, groups and policies from YAML files to AWS

For the `push` command, IAMy will output an execution plan as a series of [`aws` cli](https://aws.amazon.com/cli/) commands which can be optionally executed. This turns out to be a very direct and understandable way to display the changes to be made, and means you can pick and choose exactly what commands get actioned.

### Other features

- `fmt` will reformat all relevant files to match the output of `iamy pull`. This is particularly useful for using IAMy for drift detection, as you can use it as a PR check, and/or reformat files before performing a diff.

## Getting started

You can install IAMy on macOS with `brew install iamy`, or with the go toolchain `go get -u github.com/99designs/iamy`.

Because IAMy uses the [aws cli tool](https://aws.amazon.com/cli/), you'll want to install it first.

For configuration, IAMy uses the same [AWS environment variables](http://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html) as the aws cli. You might find [aws-vault](https://github.com/99designs/aws-vault) an excellent complementary tool for managing AWS credentials.


## Example Usage

```bash
$ iamy pull

$ find .
./myaccount-123456789/iam/user/joe.yml

$ mkdir -p myaccount-123456789/iam/user/foo

$ touch myaccount-123456789/iam/user/foo/bar.baz

$ cat << EOD > myaccount-123456789/iam/user/billy.blogs
Policies:
- arn:aws:iam::aws:policy/ReadOnly
EOD

$ iamy push
Commands to push changes to AWS:
        aws iam create-user --path /foo --user-name bar.baz
        aws iam create-user --user-name billy.blogs
        aws iam attach-user-policy --user-name billy.blogs --policy-arn arn:aws:iam::aws:policy/ReadOnly

Exec all aws commands? (y/N) y

> aws iam create-user --path /foo --user-name bar.baz
> aws iam create-user --user-name billy.blogs
> aws iam attach-user-policy --user-name billy.blogs --policy-arn arn:aws:iam::aws:policy/ReadOnly
```

## Accurate cloudformation matching

By default, iamy will use a simple heuristic (does it end with an ID, eg -ABCDEF1234) to determine if a given resource is managed by cloudformation.

This behaviour is good enough for some cases, but if you want slower but more accurate matching pass `--accurate-cfn`
to enumerate all cloudformation stacks and resources to determine exactly which resources are managed.

## Inspiration and similar tools
- https://github.com/percolate/iamer
- https://github.com/hashicorp/terraform
