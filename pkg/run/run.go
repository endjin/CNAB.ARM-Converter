package run

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/cnabio/cnab-go/credentials"
	"github.com/cnabio/cnab-go/schema"
	"github.com/cnabio/cnab-go/valuesource"
	"github.com/endjin/CNAB.ARM-Converter/pkg/common"
)

type config struct {
	cnabBundleTag        string
	cnabAction           string
	cnabInstallationName string
}

type parameterSet struct {
	SchemaVersion schema.Version         `json:"schemaVersion" yaml:"schemaVersion"`
	Name          string                 `json:"name" yaml:"name"`
	Created       time.Time              `json:"created" yaml:"created"`
	Modified      time.Time              `json:"modified" yaml:"modified"`
	Parameters    []valuesource.Strategy `json:"parameters" yaml:"parameters"`
}

//Run runs Porter with the Azure driver, using environment variables
func Run() error {

	config, err := getConfig()
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	cnabBundleTag := config.cnabBundleTag
	cnabAction := config.cnabAction
	cnabInstallationName := config.cnabInstallationName

	cmdParams := buildPorterCommandParams(cnabInstallationName, cnabAction, cnabBundleTag)

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

func buildPorterCommandParams(cnabInstallationName string, cnabAction string, cnabBundleTag string) []string {
	credsPath, err := generateCredsFile(cnabInstallationName)
	if err != nil {
		log.Fatalf("generateCredsFile command failed with %s\n", err)
	}

	paramsPath, err := generateParamsFile(cnabInstallationName)
	if err != nil {
		log.Fatalf("generateParamsFile command failed with %s\n", err)
	}

	cmdParams := []string{cnabAction, cnabInstallationName, "-d", "azure", "--tag", cnabBundleTag, "--cred", credsPath, "--parameter-set", paramsPath}

	return cmdParams
}

func getConfig() (config, error) {
	var config config
	var missing []string

	if cnabBundleTag, ok := os.LookupEnv(common.GetEnvironmentVariableNames().CnabBundleTag); ok {
		config.cnabBundleTag = cnabBundleTag
	} else {
		missing = append(missing, common.GetEnvironmentVariableNames().CnabBundleTag)
	}

	if cnabAction, ok := os.LookupEnv(common.GetEnvironmentVariableNames().CnabAction); ok {
		config.cnabAction = cnabAction
	} else {
		missing = append(missing, common.GetEnvironmentVariableNames().CnabAction)
	}

	if cnabInstallationName, ok := os.LookupEnv(common.GetEnvironmentVariableNames().CnabInstallationName); ok {
		config.cnabInstallationName = cnabInstallationName
	} else {
		missing = append(missing, common.GetEnvironmentVariableNames().CnabInstallationName)
	}

	var err error
	if len(missing) > 0 {
		err = fmt.Errorf("The following environment variables must be set but are missing: %s", strings.Join(missing, ", "))
	}

	return config, err
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
		var cred valuesource.Strategy
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

			cred = valuesource.Strategy{
				Name: key,
				Source: valuesource.Source{
					Key:   "path",
					Value: path,
				},
			}
		} else {
			key = strings.TrimPrefix(envVar, common.GetEnvironmentVariableNames().CnabCredentialPrefix)
			cred = valuesource.Strategy{
				Name: key,
				Source: valuesource.Source{
					Key:   "env",
					Value: envVar,
				},
			}
		}

		creds.Credentials = append(creds.Credentials, cred)
	}

	credFileName := cnabInstallationName + "-creds.json"
	credPath := path.Join(tempDir, credFileName)

	credData, _ := json.Marshal(creds)

	if err := ioutil.WriteFile(credPath, credData, 0644); err != nil {
		return "", err
	}

	return credPath, nil
}

func generateParamsFile(cnabInstallationName string) (string, error) {
	tempDir, _ := ioutil.TempDir("", "cnabarmdriver")

	cnabParams := getCnabParams()

	paramsFileName := cnabInstallationName + "-params.json"
	paramsPath := path.Join(tempDir, paramsFileName)

	params := parameterSet{
		Name: cnabInstallationName,
	}

	for _, cnabParam := range cnabParams {
		splits := strings.Split(cnabParam, "=")
		envVar := splits[0]
		key := strings.TrimPrefix(envVar, common.GetEnvironmentVariableNames().CnabParameterPrefix)
		params.Parameters = append(params.Parameters, valuesource.Strategy{
			Name: key,
			Source: valuesource.Source{
				Key:   "value",
				Value: os.Getenv(envVar),
			},
		})
	}
	paramsData, _ := json.Marshal(params)

	if err := ioutil.WriteFile(paramsPath, []byte(paramsData), 0644); err != nil {
		return "", err
	}

	return paramsPath, nil
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
