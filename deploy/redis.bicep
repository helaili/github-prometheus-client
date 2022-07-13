param location string = resourceGroup().location
param name string 

resource redis 'Microsoft.Cache/redis@2021-06-01' = {
  name: name
  location: location
  properties: {
    enableNonSslPort: true
    sku: {
      name: 'Basic'
      capacity: 0
      family: 'C'
    }
  }
}
