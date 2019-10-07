package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/spf13/cobra"

	"github.com/simongdavies/atfcnab/pkg/template"
)

const (
	bundlecontainerregistry = "cnabquickstartstest.azurecr.io/"
)

var bundleloc string
var outputloc string
var overwrite bool
var indent bool

var rootCmd = &cobra.Command{
	Use:   "atfcnab",
	Short: "atfcnab generates an ARM template for executing a CNAB package using Azure ACI",
	Long:  `atfcnab generates an ARM template which can be used to execute Duffle in a container using ACI to perform actions on a CNAB Package, which in turn executes the CNAB Actions using the Duffle ACI Driver   `,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return generateTemplate(bundleloc, outputloc, overwrite, indent)

	},
}

func init() {
	rootCmd.Flags().StringVarP(&bundleloc, "bundle", "b", "bundle.json", "name of bundle file to generate template for , default is bundle.json")
	rootCmd.Flags().StringVarP(&outputloc, "file", "f", "azuredeploy.json", "file name for generated template,default is azuredeploy.json")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "specifies if to overwrite the output file if it already exists, default is false")
	rootCmd.Flags().BoolVarP(&indent, "indent", "i", false, "specifies if the json output should be indented")
}

// Execute runs the template generator
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func generateTemplate(bundleloc string, outputfile string, overwrite bool, indent bool) error {

	// TODO support http uri and registry based bundle

	bundle, err := loadBundle(bundleloc)

	if err != nil {
		return err
	}

	if err = checkOutputFile(outputloc, overwrite); err != nil {
		return err
	}

	// TODO need to translate new json schema based parameter format to ARM template format

	generatedTemplate := template.NewTemplate()

	// TODO need to fix this when duffle and porter support bundle push/install from registry
	bundleName, _ := getBundleName(bundle)
	environmentVariable := template.EnvironmentVariable{
		Name:  "CNAB_BUNDLE_NAME",
		Value: bundleName,
	}

	// Set the default installation name to be the bundle name

	installationName := strings.ReplaceAll(bundleName, "/", "-")
	generatedTemplate.Parameters["cnab_installation_name"] = template.Parameter{
		Type:         "string",
		DefaultValue: installationName,
		Metadata: &template.Metadata{
			Description: "The name of the application instance.",
		},
	}

	generatedTemplate.SetContainerEnvironmentVariable(environmentVariable)

	for n, p := range bundle.Parameters {

		// Parameter names cannot contain - as they are converted into environment variables set on duffle ACI container

		// porter-debug is added automatically so can only be modified by updating porter

		if n == "porter-debug" {
			continue
		}

		if strings.Contains(n, "-") {
			return fmt.Errorf("Invalid Parameter name: %s.ARM template generation requires parameter names that can be used as environment variables", n)
		}

		// Location parameter is added to template definition automatically as ACI uses it
		if n != "location" {

			var metadata template.Metadata
			if p.Metadata != nil && p.Metadata.Description != "" {
				metadata = template.Metadata{
					Description: p.Metadata.Description,
				}
			}

			var allowedValues interface{}
			if p.AllowedValues != nil {
				allowedValues = p.AllowedValues
			}

			var defaultValue interface{}
			if p.DefaultValue != nil {
				defaultValue = p.DefaultValue
			}

			generatedTemplate.Parameters[n] = template.Parameter{
				Type:          p.DataType,
				AllowedValues: allowedValues,
				DefaultValue:  defaultValue,
				Metadata:      &metadata,
			}
		}

		environmentVariable := template.EnvironmentVariable{
			Name:  strings.ToUpper(n),
			Value: fmt.Sprintf("[parameters('%s')]", n),
		}

		generatedTemplate.SetContainerEnvironmentVariable(environmentVariable)

	}

	for n := range bundle.Credentials {

		if strings.Contains(n, "-") {
			return fmt.Errorf("Invalid Credential name: %s.ARM template generation requires credential names that can be used as environment variables", n)
		}

		var environmentVariable template.EnvironmentVariable

		// TODO update to support description and required attributes once CNAB go is updated

		// Handle TenantId and SubscriptionId as default values from ARM template functions
		if n == "azure_subscription_id" || n == "azure_tenant_id" {
			environmentVariable = template.EnvironmentVariable{
				Name:  strings.ToUpper(n),
				Value: fmt.Sprintf("[subscription().%sId]", strings.TrimSuffix(strings.TrimPrefix(n, "azure_"), "_id")),
			}
		} else {
			generatedTemplate.Parameters[n] = template.Parameter{
				Type: "securestring",
			}
			environmentVariable = template.EnvironmentVariable{
				Name:        strings.ToUpper(n),
				SecureValue: fmt.Sprintf("[parameters('%s')]", n),
			}
		}

		generatedTemplate.SetContainerEnvironmentVariable(environmentVariable)
	}

	var data []byte
	if indent {
		data, _ = json.MarshalIndent(generatedTemplate, "", "\t")
	} else {
		data, _ = json.Marshal(generatedTemplate)
	}

	if err := ioutil.WriteFile(outputloc, data, 0644); err != nil {
		return err
	}

	return nil
}
func getBundleName(bundle *bundle.Bundle) (string, error) {

	for _, i := range bundle.InvocationImages {
		if i.ImageType == "docker" {
			if i.Digest == "" {
				return strings.TrimPrefix(strings.Split(i.Image, ":")[0], bundlecontainerregistry), nil
			}
			return strings.TrimPrefix(strings.Split(i.Image, "@")[0], bundlecontainerregistry), nil
		}
	}

	return "", fmt.Errorf("Cannot get bundle name from invocationImages: %v", bundle.InvocationImages)
}
func loadBundle(source string) (*bundle.Bundle, error) {
	_, err := os.Stat(source)
	if err == nil {
		jsonFile, _ := os.Open(source)
		bundle, err := bundle.ParseReader(jsonFile)
		return &bundle, err
	}
	return nil, err
}
func checkOutputFile(dest string, overwrite bool) error {
	if _, err := os.Stat(dest); err == nil {
		if !overwrite {
			return fmt.Errorf("File %s exists and overwrite not specified", dest)
		}
	} else {
		if !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
