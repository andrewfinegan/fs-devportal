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
  github:
    gitHubRawContentHost: "https://raw.githubusercontent.com"
    gitHubContentBranch: "<content_branch_name>"
    gitHubSourceOwner: "Fiserv"
    gitHubSourceRepo: "<content_github_repo_name>"
    gitHubBaseBranch: "main"
    gitHubAuthToken: "<token>"
    gitHubUser: "username"
  
  content:
    tenantProviderApiFileName: "tenant_api.json"
    apiSpecYamlFileName: "api_swagger.yaml"
    productLayoutFileName: "product_layout.yaml"
    ```
## Code
- routers.go: Contains the routes that tenants need to implement.
