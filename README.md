# Fiserv Dev Portal

    Sample tenants code for Fiserv Dev Portal

## Run Tenant Application
- To Start tenant:
```Shell
    ./start.sh
```
- Run unittest
```Shell
    ./unittest.sh
```

## Configurations
- /resources/config.yml : Update below configuration to point to tenant content github repo
    ```Shell 

contentpath:
  tenantProviderApiFile:     "tenant_api.json"
  apiSpecYamlFile:            "api_swagger.yaml"
  productLayoutFile:          "product_layout.yaml"


github:
  gitHubRawContentHost:       "https://raw.githubusercontent.com"
  gitHubSourceOwner:          "Fiserv"
  gitHubSourceRepo:           "<tenant_repo>"
  gitHubContentBranch:        "main" 
  gitHubAuthToken:            "<auth_token>"
  gitHubUser:                 "<username>"

    ```
## Code
- routers.go: Contains the routes that tenants need to implement.
