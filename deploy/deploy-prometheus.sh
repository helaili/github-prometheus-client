export RESOURCE_GROUP=GH-Rover-Staging
export LOCATION=eastus2
export STORAGE_NAME=ghroverstagingstorage
export INSTALLATION=24886277

# Need Bicep to be installed
# az bicep install
# Need the container app extension
# az extension add --name containerapp --upgrade

az deployment group create --resource-group $RESOURCE_GROUP --template-file prometheus.bicep --parameters storageName=$STORAGE_NAME location=$LOCATION installation=$INSTALLATION

STORAGE_KEY=$(az storage account keys list --resource-group $RESOURCE_GROUP --account-name $STORAGE_NAME --query "[0].value" --output tsv)

az storage file upload \
    --account-name $STORAGE_NAME \
    --account-key $STORAGE_KEY \
    --share-name "prometheus-${INSTALLATION}-config" \
    --source "prometheus.yml" \
    --path "prometheus.yml"