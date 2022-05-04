export RESOURCE_GROUP=GH-Rover-Staging
export LOCATION=eastus2

az group delete --name $RESOURCE_GROUP --yes

az group create --name $RESOURCE_GROUP --location $LOCATION

az deployment group create --resource-group $RESOURCE_GROUP --template-file redis.bicep