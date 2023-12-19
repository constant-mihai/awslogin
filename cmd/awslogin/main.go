package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/constant-mihai/aws-login/pkg/renderer"
)

var (
	region = flag.String("region", "eu-west-1", "AWS region")
)

func main() {
	flag.Parse()

	homePath := os.Getenv("HOME")
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

	profile := ""
	renderer.Render(parse(file), func(choice string) {
		profile = choice
	})

	// TODO: I could pass aws_login as a cb, but I would have to figure out
	// how to print the output. The bubble tea renderer will swallow it.
	aws_login(profile)
	time.Sleep(500 * time.Millisecond)
	aws_credentials(profile, awsCredentials)
	set_env_variables(profile, awsEnv)
}

func set_env_variables(profile string, awsEnv string) {
	aws_profile := "export AWS_PROFILE=" + profile
	aws_region := "export AWS_REGION=" + *region

	file, err := os.OpenFile(awsEnv, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(aws_profile + "\n" + aws_region + "\n")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", aws_profile)
	fmt.Printf("%s\n", aws_region)
}

func aws_login(profile string) {
	out, err := exec.Command("bash", "-c", "aws sso login --sso-session emnify --profile "+
		profile).CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		panic(err)
	}

	fmt.Printf("%s\n", string(out))
}

// TODO: decide whether to parse the file and extend / update as required.
// This would mean that multiple tokens would be available at the same time,
// which could potentially lead to applying changes in the wrong account.
func aws_credentials(profile string, awsCredentials string) {
	file, err := os.OpenFile(awsCredentials, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("[" + profile + "]\n")
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("bash", "-c", "aws configure export-credentials --format env-no-export").
		CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		panic(err)
	}

	outStr := string(out)
	outStr = strings.Replace(outStr, "AWS_ACCESS_KEY_ID", "aws_access_key_id", 1)
	outStr = strings.Replace(outStr, "AWS_SECRET_ACCESS_KEY", "aws_secret_access_key", 1)
	outStr = strings.Replace(outStr, "AWS_SESSION_TOKEN", "aws_session_token", 1)

	_, err = file.WriteString(outStr)
	if err != nil {
		panic(err)
	}
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
