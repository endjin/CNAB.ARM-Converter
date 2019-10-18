package run

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/deislabs/cnab-go/credentials"
	"github.com/endjin/CNAB.ARM-Converter/pkg/template"
	"gopkg.in/yaml.v2"
)

//Run runs Porter with the Azure driver, using environment variables
func Run() error {

	cnabBundleName := os.Getenv(template.CnabBundleNameEnvVar)
	cnabAction := os.Getenv(template.CnabActionEnvVarName)
	cnabInstallationName := os.Getenv(template.CnabInstallationNameEnvVarName)

	cnabParams := getCnabParams()
	params := strings.Join(cnabParams, " ")

	generateCredsFile(cnabInstallationName)

	// porter install <name> -d azure --tag <bundle> --param key1=value1 --cred <credFile>

	cmd := exec.Command("porter", cnabAction, cnabInstallationName, "-d", "azure", "--tag", cnabBundleName, "--param", params, "--cred", cnabInstallationName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("porter command failed with %s\n", err)
	}

	return nil
}

func generateCredsFile(cnabInstallationName string) error {
	porterHome := os.Getenv("PORTER_HOME")

	cnabCreds := getCnabCreds()

	creds := credentials.CredentialSet{
		Name: cnabInstallationName,
	}

	for _, cnabCred := range cnabCreds {
		splits := strings.Split(cnabCred, "=")
		key := splits[0]

		credentialStrategy := credentials.CredentialStrategy{
			Name: key,
			Source: credentials.Source{
				EnvVar: key,
			},
		}

		creds.Credentials = append(creds.Credentials, credentialStrategy)
	}

	credFileName := cnabInstallationName + ".yaml"
	credPath := path.Join(porterHome, "credentials", credFileName)

	credData, _ := yaml.Marshal(creds)

	if err := ioutil.WriteFile(credPath, credData, 0644); err != nil {
		return err
	}

	return nil
}

func getCnabParams() []string {
	return getEnvVarsStartingWith("CNAB_PARAM_")
}

func getCnabCreds() []string {
	return getEnvVarsStartingWith("CNAB_CRED_")
}

func getEnvVarsStartingWith(prefix string) []string {
	environmentVariables := os.Environ()
	filterFunc := func(s string) bool { return !strings.HasPrefix(s, prefix) }
	envVars := filter(environmentVariables, filterFunc)

	for _, envVar := range envVars {
		envVar = strings.TrimPrefix(envVar, prefix)
	}

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
