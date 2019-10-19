package template

const (
	//Image is the value of the Image property for the container that runs the ARM driver
	Image = "cnabquickstarts.azurecr.io/cnabarmdriver:latest"
)

// NewCnabArmDriverTemplate creates a new instance of Template for running a CNAB bundle using cnab-azure-driver
func NewCnabArmDriverTemplate() Template {

	resources := []Resource{
		{
			Condition:  "[equals(parameters('cnab_state_storage_account_resource_group'),resourceGroup().name)]",
			Type:       "Microsoft.Storage/storageAccounts",
			Name:       "[parameters('cnab_state_storage_account_name')]",
			APIVersion: "2019-04-01",
			Location:   "[parameters('location')]",
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
			Location:   "[parameters('location')]",
			DependsOn: []string{
				"[parameters('cnab_state_storage_account_name')]",
			},
			Properties: ContainerGroupProperties{
				Containers: []Container{
					{
						Name: ContainerName,
						Properties: ContainerProperties{
							Image: Image,
							Resources: Resources{
								Requests: Requests{
									CPU:        "1.0",
									MemoryInGb: "1.5",
								},
							},
							EnvironmentVariables: []EnvironmentVariable{
								{
									Name:  CnabActionEnvVarName,
									Value: "[parameters('cnab_action')]",
								},
								{
									Name:  CnabInstallationNameEnvVarName,
									Value: "[parameters('cnab_installation_name')]",
								},
								{
									Name:  "CNAB_AZURE_LOCATION",
									Value: "[parameters('cnab_azure_location')]",
								},
								{
									Name:  "CNAB_AZURE_SUBSCRIPTION_ID",
									Value: "[subscription().subscriptionId]",
								},
								{
									Name:  "CNAB_AZURE_TENANT_ID",
									Value: "[subscription().tenantId]",
								},
								{
									Name:  "CNAB_STATE_STORAGE_ACCOUNT_NAME",
									Value: "[parameters('cnab_state_storage_account_name')]",
								},
								{
									Name:        "CNAB_STATE_STORAGE_ACCOUNT_KEY",
									SecureValue: "[parameters('cnab_state_storage_account_key')]",
								},
								{
									Name:  "CNAB_STATE_SHARE_NAME",
									Value: "[parameters('cnab_state_share_name')]",
								},
								{
									Name:  "VERBOSE",
									Value: "false",
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

		// TODO:The allowed values should be generated automatically based on ACI availability
		"location": {
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
		},
		// TODO:The allowed values should be generated automatically based on ACI availability
		"cnab_azure_location": {
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
		},
		"cnab_action": {
			Type:         "string",
			DefaultValue: "install",
			Metadata: &Metadata{
				Description: "The name of the action to be performed on the application instance.",
			},
		},
		"containerGroupName": {
			Type: "string",
			Metadata: &Metadata{
				Description: "Name for the container group",
			},
			DefaultValue: "[concat('cg-',uniqueString(resourceGroup().id, newGuid()))]",
		},
		"containerName": {
			Type: "string",
			Metadata: &Metadata{
				Description: "Name for the container",
			},
			DefaultValue: "[concat('cn-',uniqueString(resourceGroup().id, newGuid()))]",
		},
		"cnab_state_storage_account_name": {
			Type: "string",
			Metadata: &Metadata{
				Description: "The storage account name for the account for the CNAB state to be stored in, by default this will be in the current resource group and will be created if it does not exist",
			},
			DefaultValue: "[concat('cnabstate',uniqueString(resourceGroup().id))]",
		},
		"cnab_state_storage_account_key": {
			Type: "string",
			Metadata: &Metadata{
				Description: "The storage account key for the account for the CNAB state to be stored in, if this is left blank it will be looked up at runtime",
			},
			DefaultValue: "",
		},
		"cnab_state_share_name": {
			Type: "string",
			Metadata: &Metadata{
				Description: "The file share name in the storage account for the CNAB state to be stored in",
			},
			DefaultValue: "",
		},
		"cnab_state_storage_account_resource_group": {
			Type: "string",
			Metadata: &Metadata{
				Description: "The resource group name for the storage account for the CNAB state to be stored in, by default this will be in the current resource group, if this is changed to a different resource group the storage account is expected to already exist",
			},
			DefaultValue: "[resourceGroup().name]",
		},
	}

	output := Outputs{
		Output{
			Type:  "string",
			Value: "[concat('az container logs -g ',resourceGroup().name,' -n ',parameters('containerGroupName'),'  --container-name ',parameters('containerName'), ' --follow')]",
		},
	}

	template := Template{
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Resources:      resources,
		Parameters:     parameters,
		Outputs:        output,
	}

	return template
}
