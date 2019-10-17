package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/endjin/CNAB.ARM-Converter/pkg/template"
)

const (
	bundlecontainerregistry = "cnabquickstartstest.azurecr.io/"
)

// GenerateTemplate generates ARM template from bundle metadata
func GenerateTemplate(bundleloc string, outputfile string, overwrite bool, indent bool) error {

	// TODO support http uri and registry based bundle

	bundle, err := loadBundle(bundleloc)

	if err != nil {
		return err
	}

	if err = checkOutputFile(outputfile, overwrite); err != nil {
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

	// Sort parameters, because Go randomizes order when iterating a map
	var parameterKeys []string
	for parameterKey := range bundle.Parameters {
		parameterKeys = append(parameterKeys, parameterKey)
	}
	sort.Strings(parameterKeys)

	for _, parameterKey := range parameterKeys {

		parameter := bundle.Parameters[parameterKey]

		// Parameter names cannot contain - as they are converted into environment variables set on duffle ACI container

		// porter-debug is added automatically so can only be modified by updating porter

		if parameterKey == "porter-debug" {
			continue
		}

		if strings.Contains(parameterKey, "-") {
			return fmt.Errorf("Invalid Parameter name: %s.ARM template generation requires parameter names that can be used as environment variables", parameterKey)
		}

		// Location parameter is added to template definition automatically as ACI uses it
		if parameterKey != "location" {

			var metadata template.Metadata
			if parameter.Metadata != nil && parameter.Metadata.Description != "" {
				metadata = template.Metadata{
					Description: parameter.Metadata.Description,
				}
			}

			var allowedValues interface{}
			if parameter.AllowedValues != nil {
				allowedValues = parameter.AllowedValues
			}

			var defaultValue interface{}
			if parameter.DefaultValue != nil {
				defaultValue = parameter.DefaultValue
			}

			generatedTemplate.Parameters[parameterKey] = template.Parameter{
				Type:          parameter.DataType,
				AllowedValues: allowedValues,
				DefaultValue:  defaultValue,
				Metadata:      &metadata,
			}
		}

		environmentVariable := template.EnvironmentVariable{
			Name:  strings.ToUpper(parameterKey),
			Value: fmt.Sprintf("[parameters('%s')]", parameterKey),
		}

		generatedTemplate.SetContainerEnvironmentVariable(environmentVariable)

	}

	// Sort credentials, because Go randomizes order when iterating a map
	var credentialKeys []string
	for credentialKey := range bundle.Credentials {
		credentialKeys = append(credentialKeys, credentialKey)
	}
	sort.Strings(credentialKeys)

	for _, credentialKey := range credentialKeys {

		if strings.Contains(credentialKey, "-") {
			return fmt.Errorf("Invalid Credential name: %s.ARM template generation requires credential names that can be used as environment variables", credentialKey)
		}

		var environmentVariable template.EnvironmentVariable

		// TODO update to support description and required attributes once CNAB go is updated

		// Handle TenantId and SubscriptionId as default values from ARM template functions
		if credentialKey == "azure_subscription_id" || credentialKey == "azure_tenant_id" {
			environmentVariable = template.EnvironmentVariable{
				Name:  strings.ToUpper(credentialKey),
				Value: fmt.Sprintf("[subscription().%sId]", strings.TrimSuffix(strings.TrimPrefix(credentialKey, "azure_"), "_id")),
			}
		} else {
			generatedTemplate.Parameters[credentialKey] = template.Parameter{
				Type: "securestring",
			}
			environmentVariable = template.EnvironmentVariable{
				Name:        strings.ToUpper(credentialKey),
				SecureValue: fmt.Sprintf("[parameters('%s')]", credentialKey),
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

	if err := ioutil.WriteFile(outputfile, data, 0644); err != nil {
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
