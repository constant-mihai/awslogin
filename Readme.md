# Summary

This is a helper program for logging into AWS accounts.
It also dumps credentials for the command line into ~/.aws/credentials.
It will overwrite the credentials file without making a backup!
It will also create a file of the type ~/.aws/env, which are env variables to be exported.

# How to

Populate the ~/.aws/config as specified here:
https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
The awslogin utilitary will parse this file and offer a list of possible accounts.

To export `AWS_PROFILE` and `AWS_REGION` do:
```
awslogin && eval $(cat ~/.aws/env)
```
