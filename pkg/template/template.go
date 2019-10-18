package template

import "fmt"

const (
	//ContainerGroupName is the value of the ContainerGroup Resource Name property in the generated template
	ContainerGroupName = "[parameters('containerGroupName')]"
	//ContainerName is the value of the Container Resource Name property for the container that runs duffle in the generated template
	ContainerName = "[parameters('containerName')]"
)

// Template defines an ARM Template that can run a CNAB Bundle
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
	MinValue      *int        `json:"minValue,omitempty"`
	MaxValue      *int        `json:"maxValue,omitempty"`
	MinLength     *int        `json:"minLength,omitempty"`
	MaxLength     *int        `json:"maxLength,omitempty"`
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

// SetContainerEnvironmentVariable sets a environment variable for the container instance
func (template *Template) SetContainerEnvironmentVariable(environmentVariable EnvironmentVariable) error {
	container, err := findContainer(template)
	if err != nil {
		return err
	}

	container.Properties.EnvironmentVariables = append(container.Properties.EnvironmentVariables, environmentVariable)

	return nil
}

func findContainer(template *Template) (*Container, error) {
	for i := range template.Resources {
		resource := &template.Resources[i]
		if resource.Name == ContainerGroupName {
			if containerGroup, ok := resource.Properties.(ContainerGroupProperties); ok {
				for j := range containerGroup.Containers {
					container := &containerGroup.Containers[j]
					if container.Name == ContainerName {
						return container, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("Container not found in the temaple")
}
