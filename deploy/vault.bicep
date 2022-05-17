@description('Location for all resources.')
param location string = resourceGroup().location
param applicationGatewayIdentityName string 

param tenantId string 
param vaultName string
param objectId string

resource appGatewayIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
  name : applicationGatewayIdentityName
}

resource vault 'Microsoft.KeyVault/vaults@2021-11-01-preview' = {
  name: vaultName
  location: location
  properties: {
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: tenantId
    enabledForTemplateDeployment: true
    accessPolicies: [
      {
        permissions: {
          certificates: [
            'all'
          ]
          keys: [
            'all'
          ]
          secrets: [
            'all'
          ]
        }
        tenantId: tenantId 
        objectId: objectId
      }
      {
        permissions: {
          certificates: [
            'get'
          ]
          keys: [
            'get'
          ]
          secrets: [
            'get'
          ]
        }
        tenantId: tenantId 
        objectId: appGatewayIdentity.properties.principalId
      }
    ]
  }
}
