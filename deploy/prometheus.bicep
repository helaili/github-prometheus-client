param location string = resourceGroup().location

@minLength(3)
@maxLength(24)
param storageName string

param installation string
param virtualNetworkName string
param backendSubnetName string

resource storage 'Microsoft.Storage/storageAccounts@2021-08-01' existing = {
  name: storageName
}

resource containerGroup 'Microsoft.ContainerInstance/containerGroups@2021-09-01' = {
  name: 'prometheus-${installation}-ghrover'
  location: location

  properties: {
    containers: [
      {
        name: 'prometheus-${installation}-ghrover'
        properties: {
          image: 'prom/prometheus:latest'
          ports: [
            {
              port: 9090
              protocol: 'TCP'
            }
          ]
          resources: {
            requests: {
              cpu: 1
              memoryInGB: 1
            }
          }
          volumeMounts: [
            {
              name: 'prometheus-config'
              mountPath: '/etc/prometheus' 
            }
            {
              name: 'prometheus-data'
              mountPath: '/prometheus' 
            }
          ]
        }
      }
    ]
    volumes: [
      {
        name: 'prometheus-config'
        azureFile: {
          shareName: 'prometheus-${installation}-config'
          storageAccountName: storageName
          storageAccountKey: storage.listKeys().keys[0].value
        }
      }
      {
        name: 'prometheus-data'
        azureFile: {
          shareName: 'prometheus-${installation}-data'
          storageAccountName: storageName
          storageAccountKey: storage.listKeys().keys[0].value
        }
      }
    ]
    osType: 'Linux'
    restartPolicy: 'Always'
    subnetIds: [
      {
        id: resourceId('Microsoft.Network/virtualNetworks/subnets', virtualNetworkName, backendSubnetName)
      }
    ]
    ipAddress: {
      type: 'Private'
      ports: [
        {
          port: 9090
          protocol: 'TCP'
        }
      ]
    }
  }
}

