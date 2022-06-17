@description('Location for all resources.')
param location string = resourceGroup().location

param vaultName string
param certName string
param appGatewayName string
param virtualNetworkName string
param appGatewayIdentityName string 
param appGatewaySubnetName string
param prometheusClientBackendPoolName string
param prometheusBackendPoolName string

param appGatewayPublicFrontendIPName string = 'appGatewayPublicFrontendIP'
param backendHTTPSettingName string = 'backendHTTPSetting'
param backendSubnetName string
param frontendHTTPListenerName string = 'frontendHTTPListener'
param frontendSSLPortName string = 'frontendSSLPort'

var publicIPAddressName = 'GHRoverPublicIP'
var virtualNetworkPrefix = '10.0.0.0/16'
var subnetPrefix = '10.0.0.0/24'
var backendSubnetPrefix = '10.0.1.0/24'


param prometheusIPs array = []
param prometheusClientIPs array = []

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

// This was not tested
resource privateDNSZone 'Microsoft.Network/privateDnsZones@2020-06-01' = {
  name: 'ghrover-private.com'
  location: location
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
          delegations: [
            {
              name: 'ContainerInstance'
              properties: {
                serviceName: 'Microsoft.ContainerInstance/containerGroups'
              }
            }
          ]
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
        name: prometheusClientBackendPoolName
        properties: {
          backendAddresses: prometheusClientIPs
        }
      }
      {
        name: prometheusBackendPoolName
        properties: {
          backendAddresses: prometheusIPs
        }
      }
    ]
    probes: [
      {
        name: 'ping'
        properties: {
          host: '127.0.0.1'
          path: '/ping'
          interval: 30
          timeout: 30
          unhealthyThreshold: 3
        }
      }
    ]
    backendHttpSettingsCollection: [
      {
        name: backendHTTPSettingName
        properties: {
          port: 8080
          protocol: 'Http'
          cookieBasedAffinity: 'Disabled'
          pickHostNameFromBackendAddress: false
          requestTimeout: 20
          probe: {
            id: resourceId('Microsoft.Network/applicationGateways/probes', appGatewayName, 'ping')
          }
        }
      }
      {
        name: 'prometheusBackendSettings'
        properties: {
          port: 9090
          protocol: 'Http'
          cookieBasedAffinity: 'Disabled'
          requestTimeout: 20
        }
      }
    ]
    urlPathMaps: [
      {
        name: 'prometheus-client-path-map'
        properties: {
          defaultBackendAddressPool: {
            id: resourceId('Microsoft.Network/applicationGateways/backendAddressPools', appGatewayName, prometheusClientBackendPoolName)
          }
          defaultBackendHttpSettings: {
            id: resourceId('Microsoft.Network/applicationGateways/backendHttpSettingsCollection', appGatewayName, backendHTTPSettingName)
          }
          pathRules: [
            {
              name: 'ping'
              properties: {
                backendHttpSettings: {
                  id: resourceId('Microsoft.Network/applicationGateways/backendHttpSettingsCollection', appGatewayName, backendHTTPSettingName)
                }
                backendAddressPool: {
                  id: resourceId('Microsoft.Network/applicationGateways/backendAddressPools', appGatewayName, prometheusClientBackendPoolName)
                }
                paths: [
                  '/ping'
                  '/webhook'
                ]
              }
            }
          ]
        }
      }
    ]
    httpListeners: [
      {
        name: frontendHTTPListenerName
        properties: {
          protocol: 'Https'
          frontendIPConfiguration: {
            id: resourceId('Microsoft.Network/applicationGateways/frontendIPConfigurations', appGatewayName, appGatewayPublicFrontendIPName)
          }
          frontendPort: {
            id: resourceId('Microsoft.Network/applicationGateways/frontendPorts', appGatewayName, frontendSSLPortName)
          }
          sslCertificate: {
            id: resourceId('Microsoft.Network/applicationGateways/sslCertificates', appGatewayName, cert.name)
          }
          requireServerNameIndication: false
        }
      }
    ]
    requestRoutingRules: [
      {
        name: 'routingRule'
        properties: {
          ruleType: 'PathBasedRouting'
          httpListener: {
            id: resourceId('Microsoft.Network/applicationGateways/httpListeners', appGatewayName, frontendHTTPListenerName)
          }
          urlPathMap: {
            id: resourceId('Microsoft.Network/applicationGateways/urlPathMaps', appGatewayName, 'prometheus-client-path-map')
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
