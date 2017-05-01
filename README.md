# cpr - Chris Pull Request

cpr is designed to be a simple way to author pull requests from the command line.

## Usage

```
$ cpr --help
Usage of cpr:
  -base-branch string
    	The base branch to merge into (master|develop|release|staging) (Required)
  -body string
    	The description of this pull request (Optional)
  -compare-branch string
    	The branch you are attempting to merge (feature|bugfix) (Required)
  -generate-config
    	Use this flag to generate a config for your project.
  -pass string
    	Github api-key (asckoq14rf0n!@$) (Optional)
  -r string
    	A comma separated list of reviewers (Chris,Paul) (Optional)
  -reviewers string
    	A comma separated list of reviewers (Chris,Paul) (Optional)
  -title string
    	The title of this pull request (Required)
  -user string
    	Github username (alistanis) (Optional)
```

## Examples

Creating a pull request.
```
$ cpr --base-branch master --compare-branch develop --title ModifyReadme --body "Modified the readme to include usage"
  ModifyReadme  was created at  2017-05-01 23:24:00 +0000 UTC  by  alistanis
```

Generating a config.

```
$ cpr --generate-config
  Please enter your github username.
  alistanis
  Please enter your github api key.
  Would you like to use encryption for your api key? (y/n)
  y
  The encryption key will be stored at  /Users/christophercooper/.cpr-key
  The cpr config will be stored at  /Users/christophercooper/.cpr.json
```