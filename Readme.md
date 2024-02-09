# Summary

This is a helper program for logging into AWS accounts.
It also dumps credentials for the command line into `~/.aws/env` and then exports them into environment variables.

# How to

Populate the `~/.aws/config` as specified here:
https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html

To use the program, either:
a. Simply call this command from the shell:
```
eval $(cat ~/.aws/env-flush) && export AWS_PROFILE=<profile> && awslogin && eval $(cat ~/.aws/env)
```
b. Create the following function in bashrc: 
```
awslogin_wrapper() {
    if [ $# != 1 ]; then
        echo "Usage: awslogin_wrapper <profile>"
        exit 1
    fi

    case $1 in
        "--help")
            echo "Usage: awslogin_wrapper <profile>"
            exit 1
            ;;
        "-h")
            echo "Usage: awslogin_wrapper <profile>"
            exit 1
            ;;
    esac

    eval $(cat ~/.aws/env-flush) && export AWS_PROFILE=$1 && awslogin && eval $(cat ~/.aws/env)
}
```
Where `<profile>` is one of the profiles configured under `~/.aws/config`.
And where `~/.aws/env-flush` has the following contents:
```
export AWS_REGION=""
export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
export AWS_SESSION_TOKEN=""
export AWS_CREDENTIAL_EXPIRATION=""

```

`~/.aws/env-flush` is required if a new sso login is required in the same terminal.
If this isn't called, the sso login might use the credentials in the exported variables to log in into the old account.
