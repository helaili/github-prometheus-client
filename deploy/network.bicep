@description('Location for all resources.')
param location string = resourceGroup().location

param vaultName string
param certName string
param appGatewayIdentityName string 

param backendHTTPSettingName string = 'backendHTTPSetting'
param backendPoolName string = 'backendPool'
param backendSubnetName string = 'backendSubnet'
param frontendHTTPListenerName string = 'frontendHTTPListener'
param appGatewaySubnetName string = 'appGatewaySubnet'
param appGatewayPublicFrontendIPName string = 'appGatewayPublicFrontendIP'
param frontendSSLPortName string = 'frontendSSLPort'

var virtualNetworkName = 'GHRoverVNet'
var publicIPAddressName = 'GHRoverPublicIP'
var appGatewayName = 'GHRoverAppGateway'
var virtualNetworkPrefix = '10.0.0.0/16'
var subnetPrefix = '10.0.0.0/24'
var backendSubnetPrefix = '10.0.1.0/24'


resource networkingSecretsKeyVault 'Microsoft.KeyVault/vaults@2019-09-01' existing = {
  scope: resourceGroup()
  name: vaultName
}

resource cert 'Microsoft.KeyVault/vaults/secrets@2021-11-01-preview' existing = {
  name: certName
  parent: networkingSecretsKeyVault
}

resource existing_identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
  name : appGatewayIdentityName
}

resource publicIPAddress 'Microsoft.Network/publicIPAddresses@2021-05-01' = {
  name: publicIPAddressName
  location: location
  sku: {
    name: 'Standard'
  }
  properties: {
    publicIPAddressVersion: 'IPv4'
    publicIPAllocationMethod: 'Static'
    idleTimeoutInMinutes: 4
  }
}

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2021-05-01' = {
  name: virtualNetworkName
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        virtualNetworkPrefix
      ]
    }
    subnets: [
      {
        name: appGatewaySubnetName
        properties: {
          addressPrefix: subnetPrefix
          privateEndpointNetworkPolicies: 'Enabled'
          privateLinkServiceNetworkPolicies: 'Enabled'
        }
      }
      {
        name: backendSubnetName
        properties: {
          addressPrefix: backendSubnetPrefix
          privateEndpointNetworkPolicies: 'Enabled'
          privateLinkServiceNetworkPolicies: 'Enabled'
        }
      }
    ]
    enableDdosProtection: false
    enableVmProtection: false
  }
}

resource appGateway 'Microsoft.Network/applicationGateways@2021-05-01' = {
  name: appGatewayName
  location: location
  identity: {
    type:'UserAssigned'
    userAssignedIdentities:{
      '${existing_identity.id}' : {}
    }
  }
  properties: {
    sku: {
      name: 'Standard_v2'
      tier: 'Standard_v2'
    }
    gatewayIPConfigurations: [
      {
        name: 'appGatewayIPConfig'
        properties: {
          subnet: {
            id: resourceId('Microsoft.Network/virtualNetworks/subnets', virtualNetworkName, appGatewaySubnetName)
          }
        }
      }
    ]
    frontendIPConfigurations: [
      {
        name: appGatewayPublicFrontendIPName
        properties: {
          privateIPAllocationMethod: 'Dynamic'
          publicIPAddress: {
            id: resourceId('Microsoft.Network/publicIPAddresses', publicIPAddressName)
          }
        }
      }
    ]
    frontendPorts: [
      {
        name: frontendSSLPortName
        properties: {
          port: 443
        }
      }
    ]
    backendAddressPools: [
      {
        name: backendPoolName
        properties: {}
      }
    ]
    backendHttpSettingsCollection: [
      {
        name: backendHTTPSettingName
        properties: {
          port: 80
          protocol: 'Http'
          cookieBasedAffinity: 'Disabled'
          pickHostNameFromBackendAddress: false
          requestTimeout: 20
        }
      }
    ]
    httpListeners: [
      {
        name: frontendHTTPListenerName
        properties: {
          frontendIPConfiguration: {
            id: resourceId('Microsoft.Network/appGateways/frontendIPConfigurations', appGatewayName, appGatewayPublicFrontendIPName)
          }
          frontendPort: {
            id: resourceId('Microsoft.Network/appGateways/frontendPorts', appGatewayName, frontendSSLPortName)
          }
          protocol: 'Https'
          sslCertificate: {
            id: resourceId('Microsoft.Network/appGateways/sslCertificates', appGatewayName, cert.name)
          }
          requireServerNameIndication: false
        }
      }
    ]
    requestRoutingRules: [
      {
        name: 'routingRule'
        properties: {
          ruleType: 'Basic'
          httpListener: {
            id: resourceId('Microsoft.Network/appGateways/httpListeners', appGatewayName, frontendHTTPListenerName)
          }
          backendAddressPool: {
            id: resourceId('Microsoft.Network/appGateways/backendAddressPools', appGatewayName, backendPoolName)
          }
          backendHttpSettings: {
            id: resourceId('Microsoft.Network/appGateways/backendHttpSettingsCollection', appGatewayName, backendHTTPSettingName)
          }
        }
      }
    ]
    sslCertificates: [
      {
        name: cert.name
        properties: {
          keyVaultSecretId: cert.properties.secretUri
        }
      }
    ]
    enableHttp2: false
    autoscaleConfiguration: {
      minCapacity: 0
      maxCapacity: 10
    }
  }
  dependsOn: [
    virtualNetwork
    publicIPAddress
  ]
}
