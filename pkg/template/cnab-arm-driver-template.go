package template

import (
	"fmt"

	"github.com/endjin/CNAB.ARM-Converter/pkg/common"
)

const (
	//CnabArmDriverImageName is the image name for the docker image that runs the ARM driver
	CnabArmDriverImageName = "cnabquickstarts.azurecr.io/cnabarmdriver"
)

// NewCnabArmDriverTemplate creates a new instance of Template for running a CNAB bundle using cnab-azure-driver
func NewCnabArmDriverTemplate(bundleName string, bundleTag string, containerImageName string, containerImageVersion string, simplify bool) Template {

	resources := []Resource{
		{
			Condition:  "[equals(variables('cnab_state_storage_account_resource_group'),resourceGroup().name)]",
			Type:       "Microsoft.Storage/storageAccounts",
			Name:       "[variables('cnab_state_storage_account_name')]",
			APIVersion: "2019-04-01",
			Location:   "[variables('location')]",
			Sku: &Sku{
				Name: "Standard_LRS",
			},
			Kind: "StorageV2",
			Properties: StorageProperties{
				Encryption: Encryption{
					KeySource: "Microsoft.Storage",
					Services: Services{
						File: File{
							Enabled: true,
						},
					},
				},
			},
		},
		{
			Name:       ContainerGroupName,
			Type:       "Microsoft.ContainerInstance/containerGroups",
			APIVersion: "2018-10-01",
			Location:   "[variables('location')]",
			DependsOn: []string{
				"[variables('cnab_state_storage_account_name')]",
			},
			Properties: ContainerGroupProperties{
				Containers: []Container{
					{
						Name: ContainerName,
						Properties: ContainerProperties{
							Resources: Resources{
								Requests: Requests{
									CPU:        "1.0",
									MemoryInGb: "1.5",
								},
							},
							EnvironmentVariables: []EnvironmentVariable{
								{
									Name:  common.GetEnvironmentVariableNames().CnabAction,
									Value: "[variables('cnab_action')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabInstallationName,
									Value: "[variables('cnab_installation_name')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabAzureLocation,
									Value: "[variables('cnab_azure_location')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabAzureClientID,
									Value: "[variables('cnab_azure_client_id')]",
								},
								{
									Name:        common.GetEnvironmentVariableNames().CnabAzureClientSecret,
									SecureValue: "[variables('cnab_azure_client_secret')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabAzureSubscriptionID,
									Value: "[variables('cnab_azure_subscription_id')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabAzureTenantID,
									Value: "[variables('cnab_azure_tenant_id')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabStateStorageAccountName,
									Value: "[variables('cnab_state_storage_account_name')]",
								},
								{
									Name:        common.GetEnvironmentVariableNames().CnabStateStorageAccountKey,
									SecureValue: "[variables('cnab_state_storage_account_key')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabStateStorageAccountResourceGroup,
									Value: "[variables('cnab_state_storage_account_resource_group')]",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabStateShareName,
									Value: "[variables('cnab_state_share_name')]",
								},
								{
									Name:  "VERBOSE",
									Value: "false",
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabBundleName,
									Value: bundleName,
								},
								{
									Name:  common.GetEnvironmentVariableNames().CnabBundleTag,
									Value: bundleTag,
								},
							},
						},
					},
				},
				OsType:        "Linux",
				RestartPolicy: "Never",
			},
		},
	}

	parameters := map[string]Parameter{
		"cnab_action": {
			Type:         "string",
			DefaultValue: "install",
			Metadata: &Metadata{
				Description: "The name of the action to be performed on the application instance.",
			},
		},
		"cnab_azure_client_id": {
			Type:         "string",
			DefaultValue: "",
			Metadata: &Metadata{
				Description: "AAD Client ID for Azure account authentication - used to authenticate to Azure using Service Principal for ACI creation.",
			},
		},
		"cnab_azure_client_secret": {
			Type:         "securestring",
			DefaultValue: "",
			Metadata: &Metadata{
				Description: "AAD Client Secret for Azure account authentication - used to authenticate to Azure using Service Principal for ACI creation.",
			},
		},
	}

	if !simplify {
		// TODO:The allowed values should be generated automatically based on ACI availability
		parameters["location"] = Parameter{
			Type:         "string",
			DefaultValue: "[resourceGroup().Location]",
			AllowedValues: []string{
				"westus",
				"eastus",
				"westeurope",
				"westus2",
				"northeurope",
				"southeastasia",
				"eastus2",
				"centralus",
				"australiaeast",
				"uksouth",
				"southcentralus",
				"centralindia",
				"southindia",
				"northcentralus",
				"eastasia",
				"canadacentral",
				"japaneast",
			},
			Metadata: &Metadata{
				Description: "The location in which the resources will be created.",
			},
		}

		// TODO:The allowed values should be generated automatically based on ACI availability
		parameters["cnab_azure_location"] = Parameter{
			Type:         "string",
			DefaultValue: "[resourceGroup().Location]",
			AllowedValues: []string{
				"westus",
				"eastus",
				"westeurope",
				"westus2",
				"northeurope",
				"southeastasia",
				"eastus2",
				"centralus",
				"australiaeast",
				"uksouth",
				"southcentralus",
				"centralindia",
				"southindia",
				"northcentralus",
				"eastasia",
				"canadacentral",
				"japaneast",
			},
			Metadata: &Metadata{
				Description: "The location which the cnab-azure driver will use to create ACI.",
			},
		}

		parameters["cnab_azure_subscription_id"] = Parameter{
			Type:         "string",
			DefaultValue: "[subscription().subscriptionId]",
			Metadata: &Metadata{
				Description: "Azure Subscription Id - this is the subscription to be used for ACI creation, if not specified the first (random) subscription is used.",
			},
		}

		parameters["cnab_azure_tenant_id"] = Parameter{
			Type:         "string",
			DefaultValue: "[subscription().tenantId]",
			Metadata: &Metadata{
				Description: "Azure AAD Tenant Id Azure account authentication - used to authenticate to Azure using Service Principal or Device Code for ACI creation.",
			},
		}

		parameters["cnab_installation_name"] = Parameter{
			Type:         "string",
			DefaultValue: bundleName,
			Metadata: &Metadata{
				Description: "The name of the application instance.",
			},
		}

		parameters["containerGroupName"] = Parameter{
			Type: "string",
			Metadata: &Metadata{
				Description: "Name for the container group",
			},
			DefaultValue: "[concat('cg-',uniqueString(resourceGroup().id, newGuid()))]",
		}

		parameters["containerName"] = Parameter{
			Type: "string",
			Metadata: &Metadata{
				Description: "Name for the container",
			},
			DefaultValue: "[concat('cn-',uniqueString(resourceGroup().id, newGuid()))]",
		}

		parameters["cnab_state_storage_account_name"] = Parameter{
			Type: "string",
			Metadata: &Metadata{
				Description: "The storage account name for the account for the CNAB state to be stored in, by default this will be in the current resource group and will be created if it does not exist",
			},
			DefaultValue: "[concat('cnabstate',uniqueString(resourceGroup().id))]",
		}

		parameters["cnab_state_storage_account_key"] = Parameter{
			Type: "string",
			Metadata: &Metadata{
				Description: "The storage account key for the account for the CNAB state to be stored in, if this is left blank it will be looked up at runtime",
			},
			DefaultValue: "",
		}

		parameters["cnab_state_share_name"] = Parameter{
			Type: "string",
			Metadata: &Metadata{
				Description: "The file share name in the storage account for the CNAB state to be stored in",
			},
			DefaultValue: "",
		}

		parameters["cnab_state_storage_account_resource_group"] = Parameter{
			Type: "string",
			Metadata: &Metadata{
				Description: "The resource group name for the storage account for the CNAB state to be stored in, by default this will be in the current resource group, if this is changed to a different resource group the storage account is expected to already exist",
			},
			DefaultValue: "[resourceGroup().name]",
		}

	}

	output := Outputs{
		Output{
			Type:  "string",
			Value: "[concat('az container logs -g ',resourceGroup().name,' -n ',variables('containerGroupName'),'  --container-name ',variables('containerName'), ' --follow')]",
		},
	}

	template := Template{
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Resources:      resources,
		Parameters:     parameters,
		Outputs:        output,
	}

	template.setContainerImage(containerImageName, containerImageVersion)

	if simplify {
		template.addSimpleVariables(bundleName, bundleTag)
	} else {
		template.addAdvancedVariables()
	}

	return template
}

func (template *Template) addAdvancedVariables() {
	variables := map[string]string{
		"cnab_action":                               "[parameters('cnab_action')]",
		"cnab_azure_client_id":                      "[parameters('cnab_azure_client_id')]",
		"cnab_azure_client_secret":                  "[parameters('cnab_azure_client_secret')]",
		"cnab_azure_location":                       "[parameters('cnab_azure_location')]",
		"cnab_azure_subscription_id":                "[parameters('cnab_azure_subscription_id')]",
		"cnab_azure_tenant_id":                      "[parameters('cnab_azure_tenant_id')]",
		"cnab_installation_name":                    "[parameters('cnab_installation_name')]",
		"cnab_state_share_name":                     "[parameters('cnab_state_share_name')]",
		"cnab_state_storage_account_key":            "[parameters('cnab_state_storage_account_key')]",
		"cnab_state_storage_account_name":           "[parameters('cnab_state_storage_account_name')]",
		"cnab_state_storage_account_resource_group": "[parameters('cnab_state_storage_account_resource_group')]",
		"containerGroupName":                        "[parameters('containerGroupName')]",
		"containerName":                             "[parameters('containerName')]",
		"location":                                  "[parameters('location')]",
	}

	template.Variables = variables
}

func (template *Template) addSimpleVariables(bundleName string, bundleTag string) {
	variables := map[string]string{
		"cnab_action":                               "[parameters('cnab_action')]",
		"cnab_azure_client_id":                      "[parameters('cnab_azure_client_id')]",
		"cnab_azure_client_secret":                  "[parameters('cnab_azure_client_secret')]",
		"cnab_azure_location":                       "[resourceGroup().Location]",
		"cnab_azure_subscription_id":                "[subscription().subscriptionId]",
		"cnab_azure_tenant_id":                      "[subscription().tenantId]",
		"cnab_installation_name":                    bundleName,
		"cnab_state_share_name":                     bundleName,
		"cnab_state_storage_account_key":            "",
		"cnab_state_storage_account_name":           "[concat('cnabstate',uniqueString(resourceGroup().id))]",
		"cnab_state_storage_account_resource_group": "[resourceGroup().name]",
		"containerGroupName":                        fmt.Sprintf("[concat('cg-',uniqueString(resourceGroup().id, '%s', '%s'))]", bundleName, bundleTag),
		"containerName":                             fmt.Sprintf("[concat('cn-',uniqueString(resourceGroup().id, '%s', '%s'))]", bundleName, bundleTag),
		"location":                                  "[resourceGroup().Location]",
	}

	template.Variables = variables
}
