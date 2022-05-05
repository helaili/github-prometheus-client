param location string = resourceGroup().location

@minLength(3)
@maxLength(24)
param storageName string

param installation string

resource storage 'Microsoft.Storage/storageAccounts@2021-08-01' = {
  name: storageName
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'

  resource fileServices 'fileServices@2021-08-01' = {
    name: 'default'

    resource data 'shares@2021-08-01' = {
      name: 'prometheus-${installation}-data'
    }

    resource config 'shares@2021-08-01' = {
      name: 'prometheus-${installation}-config'
    }
  }
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
    ipAddress: {
      type: 'Public'
      dnsNameLabel: 'prometheus-${installation}-ghrover'
      ports: [
        {
          port: 9090
          protocol: 'TCP'
        }
      ]
    }
  }
}
