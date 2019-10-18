package template

const (
	//Image is the value of the Image property for the container that runs duffle in the generated template
	Image = "cnabquickstartstest.azurecr.io/simongdavies/run-duffle:latest"
)

// NewDuffleAciDriverTemplate creates a new instance of Template for running a CNAB bundle using duffle-aci-docker
func NewDuffleAciDriverTemplate() Template {

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
									Name:  "CNAB_ACTION",
									Value: "[parameters('cnab_action')]",
								},
								{
									Name:  "CNAB_INSTALLATION_NAME",
									Value: "[parameters('cnab_installation_name')]",
								},
								{
									Name:  "ACI_LOCATION",
									Value: "[parameters('location')]",
								},
								{
									Name:  "AZURE_SUBSCRIPTION_ID",
									Value: "[subscription().subscriptionId]",
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

		// TODO:This needs to be renamed to ACI_LOCATION once changes to driver are done so that ACI adn resources can be in different locations
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
