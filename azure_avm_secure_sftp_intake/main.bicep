targetScope = 'resourceGroup'

@description('Globally unique Azure Storage account name.')
@minLength(3)
@maxLength(24)
param storageAccountName string

@description('Azure region for the SFTP landing account.')
param location string = resourceGroup().location

@description('Private quarantine container used as the SFTP home directory.')
@minLength(3)
@maxLength(63)
param quarantineContainerName string = 'sftp-quarantine'

@description('Local SFTP user name scoped to the quarantine container.')
@minLength(3)
param sftpUserName string

@secure()
@description('OpenSSH public key. Password authentication is prohibited.')
param sshPublicKey string

@description('Subnet resource ID for the Blob private endpoint.')
param privateEndpointSubnetResourceId string

@description('Private DNS zone resource ID for privatelink.blob.core.windows.net.')
param privateDnsZoneResourceId string

@description('Log Analytics workspace resource ID for storage data-plane diagnostics.')
param logAnalyticsWorkspaceResourceId string

@description('Soft-delete retention for landing bytes. The separate retained evidence store remains mandatory.')
@minValue(7)
@maxValue(365)
param softDeleteRetentionDays int = 30

module storage 'br/public:avm/res/storage/storage-account:0.33.0' = {
  name: 'secure-sftp-${uniqueString(storageAccountName, resourceGroup().id)}'
  params: {
    name: storageAccountName
    location: location
    skuName: 'Standard_ZRS'
    kind: 'StorageV2'
    enableHierarchicalNamespace: true
    enableSftp: true
    isLocalUserEnabled: true
    allowSharedKeyAccess: false
    defaultToOAuthAuthentication: true
    allowBlobPublicAccess: false
    publicNetworkAccess: 'Disabled'
    supportsHttpsTrafficOnly: true
    minimumTlsVersion: 'TLS1_2'
    requireInfrastructureEncryption: true
    networkAcls: {
      bypass: 'None'
      defaultAction: 'Deny'
      ipRules: []
      virtualNetworkRules: []
    }
    privateEndpoints: [
      {
        service: 'blob'
        subnetResourceId: privateEndpointSubnetResourceId
        privateDnsZoneGroup: {
          privateDnsZoneGroupConfigs: [
            {
              privateDnsZoneResourceId: privateDnsZoneResourceId
            }
          ]
        }
      }
    ]
    blobServices: {
      containerDeleteRetentionPolicyEnabled: true
      containerDeleteRetentionPolicyDays: softDeleteRetentionDays
      deleteRetentionPolicyEnabled: true
      deleteRetentionPolicyDays: softDeleteRetentionDays
      isVersioningEnabled: false
      containers: [
        {
          name: quarantineContainerName
          publicAccess: 'None'
        }
      ]
      diagnosticSettings: [
        {
          name: 'blob-audit'
          workspaceResourceId: logAnalyticsWorkspaceResourceId
          logCategoriesAndGroups: [
            { category: 'StorageRead' }
            { category: 'StorageWrite' }
            { category: 'StorageDelete' }
          ]
          metricCategories: [
            { category: 'Transaction' }
          ]
        }
      ]
    }
    localUsers: [
      {
        name: sftpUserName
        hasSharedKey: false
        hasSshKey: true
        hasSshPassword: false
        homeDirectory: quarantineContainerName
        sshAuthorizedKeys: [
          {
            description: 'project-managed upload key'
            key: sshPublicKey
          }
        ]
        permissionScopes: [
          {
            permissions: 'cw'
            service: 'blob'
            resourceName: quarantineContainerName
          }
        ]
      }
    ]
  }
}

output storageAccountResourceId string = storage.outputs.resourceId
output privateBlobEndpoint string = storage.outputs.primaryBlobEndpoint
