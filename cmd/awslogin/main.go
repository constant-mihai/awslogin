package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	region   = flag.String("region", "eu-west-1", "AWS region")
	homePath = os.Getenv("HOME")
)

func main() {
	flag.Parse()

	if homePath == "" {
		panic("home env variable is missing")
	}
	awsConfig := homePath + "/.aws/config"
	awsCredentials := homePath + "/.aws/credentials"
	awsEnv := homePath + "/.aws/env"

	file, err := os.Open(awsConfig)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	aws_login()
	time.Sleep(500 * time.Millisecond)
	exportVars := aws_credentials(awsCredentials)
	set_env_variables(exportVars, awsEnv)
}

func set_env_variables(exportVars string, awsEnv string) {
	awsRegion := "export AWS_REGION=" + *region

	file, err := os.OpenFile(awsEnv, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(awsRegion + "\n" + exportVars + "\n")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", awsRegion)
	fmt.Printf("%s\n", exportVars)
}

func aws_login() {
	file, err := os.Open(homePath + "/.aws/config")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	out, err := exec.Command("bash", "-c", "aws sso login").CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		panic(err)
	}

	fmt.Printf("%s\n", string(out))
}

// TODO: decide whether to parse the file and extend / update as required.
// This would mean that multiple tokens would be available at the same time,
// which could potentially lead to applying changes in the wrong account.
func aws_credentials(awsCredentials string) string {
	file, err := os.OpenFile(awsCredentials, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("[awslogin]\n")
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("bash", "-c", "aws configure export-credentials --format=env").
		CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		panic(err)
	}

	outStr := string(out)
	exportVars := outStr
	outStr = strings.ReplaceAll(outStr, "export ", "")
	outStr = strings.Replace(outStr, "AWS_ACCESS_KEY_ID", "aws_access_key_id", 1)
	outStr = strings.Replace(outStr, "AWS_SECRET_ACCESS_KEY", "aws_secret_access_key", 1)
	outStr = strings.Replace(outStr, "AWS_SESSION_TOKEN", "aws_session_token", 1)
	outStr = strings.Replace(outStr, "AWS_CREDENTIAL_EXPIRATION", "aws_credential_expiration", 1)

	_, err = file.WriteString(outStr)
	if err != nil {
		panic(err)
	}

	return exportVars
}

func parse(file *os.File) []string {
	profiles := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "profile") {
			profile := strings.TrimPrefix(scanner.Text(), "[")
			profile = strings.TrimSuffix(profile, "]")
			_, profile, _ = strings.Cut(profile, " ")
			profiles = append(profiles, profile)
		}
	}
	return profiles
}
