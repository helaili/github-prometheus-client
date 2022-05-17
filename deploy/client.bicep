param location string = resourceGroup().location

param registry_username string
@secure()
param registry_password string

param port int = 8080
param app_id string
@secure()
param webhook_secret string
@secure()
param private_key string
param env string
param redisServer string
param redisPort int = 6379

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' = {
  name: 'github-prometheus-client'
  location: location

  properties: {
    sku: 'Standard'
    imageRegistryCredentials: [
      {
        server: 'ghcr.io'
        username: registry_username
        password: registry_password
      }
    ]
    containers: [
      {
        name: 'github-prometheus-client'
        properties: {
          image: 'ghcr.io/helaili/github-prometheus-client:main'
          ports: [
            {
              port: port
              protocol: 'TCP'
            }
          ]
          resources: {
            requests: {
              cpu: 1
              memoryInGB: 1
            }
          }
          environmentVariables: [
            {
              name: 'APP_ID'
              value: app_id
            }
            {
              name: 'WEBHOOK_SECRET'
              secureValue: webhook_secret
            }
            {
              name: 'PRIVATE_KEY'
              secureValue: private_key
            }
            {
              name: 'PORT'
              value: string(port)
            }
            {
              name: 'GITHUB_PROMETHEUS_CLIENT_ENV'
              value: env
            }
            {
              name: 'REDIS_SERVER'
              value: redisServer
            }
            {
              name: 'REDIS_PORT'
              value: string(redisPort)
            }
          ]
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Always'
    subnetIds: [
      {
        id: resourceId('Microsoft.Network/virtualNetworks/subnets', 'myVNet', 'myAGSubnet')
      }
    ]
    ipAddress: {
      type: 'Private'
      dnsNameLabel: 'ghrover'
      ports: [
        {
          port: port
          protocol: 'TCP'
        }
      ]
    }
  }
}

