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
	awsEnv := homePath + "/.aws/env"

	aws_login()
	time.Sleep(500 * time.Millisecond)
	exportVars := aws_credentials()
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
	out, err := exec.Command("bash", "-c", "aws sso login").CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		panic(err)
	}

	fmt.Printf("%s\n", string(out))
}

func aws_credentials() string {
	out, err := exec.Command("bash", "-c", "aws configure export-credentials --format=env").
		CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
		panic(err)
	}

	exportVars := string(out)

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
