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

	// TODO validate environment variables are set

	cnabBundleTag := os.Getenv(template.CnabBundleTagEnvVar)
	cnabAction := os.Getenv(template.CnabActionEnvVarName)
	cnabInstallationName := os.Getenv(template.CnabInstallationNameEnvVarName)

	cnabParams := getCnabParams()

	credsPath, err := generateCredsFile(cnabInstallationName)
	if err != nil {
		log.Fatalf("generateCredsFile command failed with %s\n", err)
	}

	cmdParams := []string{cnabAction, cnabInstallationName, "-d", "azure", "--tag", cnabBundleTag, "--cred", credsPath}
	for i := range cnabParams {
		cmdParams = append(cmdParams, "--param")
		cmdParams = append(cmdParams, strings.TrimPrefix(cnabParams[i], "CNAB_PARAM_"))
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
	cnabCreds := getCnabCreds()

	creds := credentials.CredentialSet{
		Name: cnabInstallationName,
	}

	for _, cnabCred := range cnabCreds {
		splits := strings.Split(cnabCred, "=")
		envVar := splits[0]
		key := strings.TrimPrefix(envVar, "CNAB_CRED_")

		credentialStrategy := credentials.CredentialStrategy{
			Name: key,
			Source: credentials.Source{
				EnvVar: envVar,
			},
		}

		creds.Credentials = append(creds.Credentials, credentialStrategy)
	}

	credFileName := cnabInstallationName + ".yaml"
	tempDir, _ := ioutil.TempDir("", "cnabarmdriver")
	credPath := path.Join(tempDir, credFileName)

	credData, _ := yaml.Marshal(creds)

	if err := ioutil.WriteFile(credPath, credData, 0644); err != nil {
		return "", err
	}

	return credPath, nil
}

func getCnabParams() []string {
	return getEnvVarsStartingWith("CNAB_PARAM_")
}

func getCnabCreds() []string {
	return getEnvVarsStartingWith("CNAB_CRED_")
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
