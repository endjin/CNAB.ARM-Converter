package template

const (
	//ContainerGroupName is the value of the ContainerGroup Resource Name property in the generated template
	ContainerGroupName = "[parameters('containerGroupName')]"
	//ContainerName is the value of the Container Resource Name property for the container that runs duffle in the generated template
	ContainerName = "[parameters('containerName')]"
	//Image is the value of the Image property for the container that runs duffle in the generated template
	Image = "cnabquickstartstest.azurecr.io/simongdavies/run-duffle:latest"
)

// Template defines an ARM Template that can run a CNAB Bundle using duffle-aci-docker
type Template struct {
	Schema         string               `json:"$schema"`
	ContentVersion string               `json:"contentVersion"`
	Parameters     map[string]Parameter `json:"parameters"`
	Resources      []Resource           `json:"resources"`
	Outputs        Outputs              `json:"outputs"`
}

// Metadata defines the metadata for a template parameter
type Metadata struct {
	Description string `json:"description,omitempty"`
}

// Parameter defines a template parameter
type Parameter struct {
	Type          string      `json:"type"`
	DefaultValue  interface{} `json:"defaultValue,omitempty"`
	AllowedValues interface{} `json:"allowedValues,omitempty"`
	Metadata      *Metadata   `json:"metadata,omitempty"`
}

// Sku is defines a SKU for template resource
type Sku struct {
	Name string `json:"name,omitempty"`
}

// File defines if encryption is enabled for file shares in a storage account created by the template
type File struct {
	Enabled bool `json:"enabled"`
}

// Services defines Services that can be encrypted in a storage account
type Services struct {
	File File `json:"file"`
}

// Encryption defines the encryption properties for the storage account in the generated template
type Encryption struct {
	KeySource string   `json:"keySource"`
	Services  Services `json:"services"`
}

// StorageProperties defines the properties of the storage account in the generated template
type StorageProperties struct {
	Encryption Encryption `json:"encryption"`
}

// Requests defines the CPU and Memorty requirements of the Container instance in the generated template
type Requests struct {
	CPU        string `json:"cpu"`
	MemoryInGb string `json:"memoryInGb"`
}

// Resources defines the resource requests for the Container instance in the generated template
type Resources struct {
	Requests Requests `json:"requests"`
}

// EnvironmentVariable defines the environment variables that are created for the container in the generated template
type EnvironmentVariable struct {
	Name        string `json:"name"`
	SecureValue string `json:"secureValue,omitempty"`
	Value       string `json:"value,omitempty"`
}

//ContainerProperties define the properties of the container resource in the generated template
type ContainerProperties struct {
	Image                string                `json:"image"`
	Resources            Resources             `json:"resources"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
}

// Container defines the container in the generated template
type Container struct {
	Name       string              `json:"name"`
	Properties ContainerProperties `json:"properties"`
}

//ContainerGroupProperties defines the properties of the Container Group in the generated template
type ContainerGroupProperties struct {
	Containers    []Container `json:"containers"`
	OsType        string      `json:"osType"`
	RestartPolicy string      `json:"restartPolicy"`
}

// Resource defines a resource in the generated template
type Resource struct {
	Condition  string      `json:"condition,omitempty"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	APIVersion string      `json:"apiVersion"`
	Location   string      `json:"location"`
	Sku        *Sku        `json:"sku,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	DependsOn  []string    `json:"dependsOn,omitempty"`
	Properties interface{} `json:"properties"`
}

// Output defines an output in the generated template
type Output struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Outputs defines the outputs in the genreted template
type Outputs struct {
	CNABPackageActionLogsCommand Output `json:"CNAB Package Action Logs Command"`
}

// NewTemplate creates a new instance of Template
func NewTemplate() Template {

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

// SetContainerEnvironmentVariable sets a environment variable for the container instance that runs duffle
func (t *Template) SetContainerEnvironmentVariable(environmentVariable EnvironmentVariable) {
	for _, r := range t.Resources {
		if r.Name == ContainerGroupName {
			if cg, ok := r.Properties.(ContainerGroupProperties); ok {
				for i, c := range cg.Containers {
					if c.Name == ContainerName {
						c.Properties.EnvironmentVariables = append(c.Properties.EnvironmentVariables, environmentVariable)
						cg.Containers[i] = c
					}
				}
			}
		}
	}
}
