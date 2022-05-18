param location string = resourceGroup().location

param registryUsername string
@secure()
param registryPassword string

param port int = 8080
param app_id string
@secure()
param webhook_secret string
@secure()
param private_key string
param env string
param redisServerName string 

resource redis 'Microsoft.Cache/redis@2021-06-01' existing = {
  name: redisServerName
}

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' = {
  name: 'github-prometheus-client'
  location: location

  properties: {
    sku: 'Standard'
    imageRegistryCredentials: [
      {
        server: 'ghcr.io'
        username: registryUsername
        password: registryPassword
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
              name: 'REDIS_ADDRESS'
              value: '${redis.properties.hostName}:${redis.properties.sslPort}'
            }
            {
              name: 'REDIS_PASSWORD'
              secureValue: redis.properties.accessKeys.primaryKey
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

