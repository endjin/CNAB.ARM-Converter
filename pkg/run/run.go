package run

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/deislabs/cnab-go/credentials"
	"github.com/endjin/CNAB.ARM-Converter/pkg/common"
	"gopkg.in/yaml.v2"
)

//Run runs Porter with the Azure driver, using environment variables
func Run() error {

	// TODO validate environment variables are set

	cnabBundleTag := os.Getenv(common.GetEnvironmentVariableNames().CnabBundleTag)
	cnabAction := os.Getenv(common.GetEnvironmentVariableNames().CnabAction)
	cnabInstallationName := os.Getenv(common.GetEnvironmentVariableNames().CnabInstallationName)

	cnabParams := getCnabParams()

	credsPath, err := generateCredsFile(cnabInstallationName)
	if err != nil {
		log.Fatalf("generateCredsFile command failed with %s\n", err)
	}

	cmdParams := []string{cnabAction, cnabInstallationName, "-d", "azure", "--tag", cnabBundleTag, "--cred", credsPath}
	for i := range cnabParams {
		cmdParams = append(cmdParams, "--param")
		cmdParams = append(cmdParams, strings.TrimPrefix(cnabParams[i], common.GetEnvironmentVariableNames().CnabParameterPrefix))
	}

	cmd := exec.Command("porter", cmdParams...)
	log.Println(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("porter command failed with %s\n", err)
	}

	return nil
}

func generateCredsFile(cnabInstallationName string) (string, error) {
	tempDir, _ := ioutil.TempDir("", "cnabarmdriver")

	cnabCreds := getCnabCreds()

	creds := credentials.CredentialSet{
		Name: cnabInstallationName,
	}

	for _, cnabCred := range cnabCreds {
		splits := strings.Split(cnabCred, "=")
		envVar := splits[0]

		var key string
		var credentialStrategy credentials.CredentialStrategy
		if strings.HasPrefix(envVar, common.GetEnvironmentVariableNames().CnabCredentialFilePrefix) {
			key = strings.TrimPrefix(envVar, common.GetEnvironmentVariableNames().CnabCredentialFilePrefix)

			data, err := base64.StdEncoding.DecodeString(os.Getenv(envVar))
			if err != nil {
				return "", fmt.Errorf("Unable to decode %s: %s", key, err)
			}

			path := path.Join(tempDir, key)
			if err := ioutil.WriteFile(path, data, 0644); err != nil {
				return "", err
			}

			credentialStrategy = credentials.CredentialStrategy{
				Name: key,
				Source: credentials.Source{
					Path: path,
				},
			}
		} else {
			key = strings.TrimPrefix(envVar, common.GetEnvironmentVariableNames().CnabCredentialPrefix)
			credentialStrategy = credentials.CredentialStrategy{
				Name: key,
				Source: credentials.Source{
					EnvVar: envVar,
				},
			}
		}

		creds.Credentials = append(creds.Credentials, credentialStrategy)
	}

	credFileName := cnabInstallationName + ".yaml"
	credPath := path.Join(tempDir, credFileName)

	credData, _ := yaml.Marshal(creds)

	if err := ioutil.WriteFile(credPath, credData, 0644); err != nil {
		return "", err
	}

	return credPath, nil
}

func getCnabParams() []string {
	return getEnvVarsStartingWith(common.GetEnvironmentVariableNames().CnabParameterPrefix)
}

func getCnabCreds() []string {
	return getEnvVarsStartingWith(common.GetEnvironmentVariableNames().CnabCredentialPrefix)
}

func getEnvVarsStartingWith(prefix string) []string {
	environmentVariables := os.Environ()
	filterFunc := func(s string) bool { return strings.HasPrefix(s, prefix) }
	envVars := filter(environmentVariables, filterFunc)

	return envVars
}

func filter(array []string, filterCondition func(string) bool) (result []string) {
	for _, item := range array {
		if filterCondition(item) {
			result = append(result, item)
		}
	}
	return
}
