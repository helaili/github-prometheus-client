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
