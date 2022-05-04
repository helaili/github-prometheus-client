param location string = resourceGroup().location

resource redis 'Microsoft.Cache/redis@2021-06-01' = {
  name: 'ghrover-staging-redis'
  location: location
  properties: {
    sku: {
      name: 'Basic'
      capacity: 0
      family: 'C'
    }
  }
}
