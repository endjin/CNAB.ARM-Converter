package common

// EnvironmentVariableNames defines environment variables names
type EnvironmentVariableNames struct {
	CnabParameterPrefix         string
	CnabCredentialPrefix        string
	CnabCredentialFilePrefix    string
	CnabAction                  string
	CnabInstallationName        string
	CnabBundleName              string
	CnabBundleTag               string
	CnabAzureLocation           string
	CnabAzureClientID           string
	CnabAzureClientSecret       string
	CnabAzureSubscriptionID     string
	CnabAzureTenantID           string
	CnabStateStorageAccountName string
	CnabStateStorageAccountKey  string
	CnabStateShareName          string
	Verbose                     string
}

// GetEnvironmentVariableNames returns environment variable names
func GetEnvironmentVariableNames() EnvironmentVariableNames {
	return EnvironmentVariableNames{
		CnabParameterPrefix:         "CNAB_PARAM_",
		CnabCredentialPrefix:        "CNAB_CRED_",
		CnabCredentialFilePrefix:    "CNAB_CRED_FILE_",
		CnabAction:                  "CNAB_ACTION",
		CnabInstallationName:        "CNAB_INSTALLATION_NAME",
		CnabBundleName:              "CNAB_BUNDLE_NAME",
		CnabBundleTag:               "CNAB_BUNDLE_TAG",
		CnabAzureLocation:           "CNAB_AZURE_LOCATION",
		CnabAzureClientID:           "CNAB_AZURE_CLIENT_ID",
		CnabAzureClientSecret:       "CNAB_AZURE_CLIENT_SECRET",
		CnabAzureSubscriptionID:     "CNAB_AZURE_SUBSCRIPTION_ID",
		CnabAzureTenantID:           "CNAB_AZURE_TENANT_ID",
		CnabStateStorageAccountName: "CNAB_STATE_STORAGE_ACCOUNT_NAME",
		CnabStateStorageAccountKey:  "CNAB_STATE_STORAGE_ACCOUNT_KEY",
		CnabStateShareName:          "CNAB_STATE_SHARE_NAME",
		Verbose:                     "VERBOSE",
	}
}
