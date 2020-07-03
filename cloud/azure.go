package cloud

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// AzureEnvironment interface
type AzureEnvironment struct {
	ClusterID           string
	SubscriptionID      string
	Region              string
	ClientID            string
	ClientSecret        string
	Tenant              string
	ResourceGroup       string
	VMNet               string
	VMSubnet            string
	VMSize              string
	ImageURI            string
	ImageStorageAccount string
	WorkerNodes         int64
	EnvParams           []string
}

func (e *AzureEnvironment) getNicClient() (network.InterfacesClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := network.NewInterfacesClient(e.SubscriptionID)
	authorizer, err := authConfig.Authorizer()
	client.Authorizer = authorizer
	return client, err
}

func (e *AzureEnvironment) getPublicIPClient() (network.PublicIPAddressesClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := network.NewPublicIPAddressesClient(e.SubscriptionID)
	authorizer, err := authConfig.Authorizer()
	client.Authorizer = authorizer
	return client, err
}

func (e *AzureEnvironment) getVMClient() (compute.VirtualMachinesClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := compute.NewVirtualMachinesClient(e.SubscriptionID)
	authorizer, err := authConfig.Authorizer()
	client.Authorizer = authorizer
	return client, err
}

func (e *AzureEnvironment) getSubnetClient() (network.SubnetsClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := network.NewSubnetsClient(e.SubscriptionID)
	authorizer, err := authConfig.Authorizer()
	client.Authorizer = authorizer
	return client, err
}

func (e *AzureEnvironment) getDiskClient() (compute.DisksClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := compute.NewDisksClient(e.SubscriptionID)
	authorizer, err := authConfig.Authorizer()
	client.Authorizer = authorizer
	return client, err
}

func (e *AzureEnvironment) getSubnet(ctx context.Context, vnetName string, subnetName string) (network.Subnet, error) {
	subnetsClient, _ := e.getSubnetClient()
	return subnetsClient.Get(ctx, e.ResourceGroup, vnetName, subnetName, "")
}

func (e *AzureEnvironment) createNIC(name string) (string, error) {
	cli, err := e.getNicClient()
	if err != nil {
		return "", err
	}

	subnet, err := e.getSubnet(context.Background(), e.VMNet, e.VMSubnet)
	if err != nil {
		return "", err
	}

	nicPath := "/subscriptions/" + e.SubscriptionID +
		"/resourceGroups/" + e.ResourceGroup +
		"/providers/Microsoft.Network/networkInterfaces/" + name

	nicParams := network.Interface{
		Name:     to.StringPtr(name),
		Location: to.StringPtr(e.Region),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				{
					Name: to.StringPtr(name),
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						Subnet:                    &subnet,
						PrivateIPAllocationMethod: network.Dynamic,
					},
				},
			},
		},
	}

	future, err := cli.CreateOrUpdate(context.Background(), e.ResourceGroup, name, nicParams)
	if err != nil {
		return "", err
	}

	err = future.WaitForCompletionRef(context.Background(), cli.Client)
	return nicPath, err
}

func (e *AzureEnvironment) deleteNIC(name string) error {
	cli, err := e.getNicClient()
	if err != nil {
		return err
	}
	future, err := cli.Delete(context.Background(), e.ResourceGroup, name)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(context.Background(), cli.Client)
}

func (e *AzureEnvironment) createDisk(name string) (string, error) {
	cli, err := e.getDiskClient()
	if err != nil {
		return "", err
	}

	storageAccountID := "/subscriptions/" + e.SubscriptionID +
		"/resourceGroups/" + e.ResourceGroup +
		"/providers/Microsoft.Storage/storageAccounts/" +
		e.ImageStorageAccount

	diskPath := "/subscriptions/" + e.SubscriptionID +
		"/resourceGroups/" + e.ResourceGroup +
		"/providers/Microsoft.Compute/disks/" + name

	disk := compute.Disk{
		Location: to.StringPtr(e.Region),
		Name:     to.StringPtr(name),
		DiskProperties: &compute.DiskProperties{
			DiskSizeGB: to.Int32Ptr(30),
			CreationData: &compute.CreationData{
				CreateOption:     compute.Import,
				SourceURI:        to.StringPtr(e.ImageURI),
				StorageAccountID: to.StringPtr(storageAccountID),
			},
		},
	}

	future, err := cli.CreateOrUpdate(context.Background(), e.ResourceGroup, name, disk)
	if err != nil {
		return "", err
	}

	return diskPath, future.WaitForCompletionRef(context.Background(), cli.Client)
}

func (e *AzureEnvironment) getPrivateIP(name string) (string, error) {
	cli, err := e.getNicClient()
	if err != nil {
		return "", err
	}

	nics, err := cli.List(context.Background(), e.ResourceGroup)
	if err != nil {
		return "", err
	}

	for _, el := range nics.Values() {
		if name == *el.Name {
			return *(*el.IPConfigurations)[0].PrivateIPAddress, nil
		}
	}

	return "", errors.New("private IP not found for VM " + name)
}

func (e *AzureEnvironment) launchVM(name string, waitForCompletion bool) (string, error) {
	cli, err := e.getVMClient()

	if err != nil {
		return "", err
	}
	ctx := context.Background()

	nic, err := e.createNIC(name)
	if err != nil {
		return "", err
	}

	disk, err := e.createDisk(name)
	if err != nil {
		return "", err
	}

	vmParameters := compute.VirtualMachine{
		Location: to.StringPtr(e.Region),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypesStandardD8sV3,
			},

			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Name:         to.StringPtr(name),
					CreateOption: compute.DiskCreateOptionTypesAttach,
					OsType:       compute.Linux,
					ManagedDisk: &compute.ManagedDiskParameters{
						StorageAccountType: compute.StorageAccountTypesStandardLRS,
						ID:                 to.StringPtr(disk),
					},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: to.StringPtr(nic),
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
	}
	future, err := cli.CreateOrUpdate(ctx, e.ResourceGroup, e.ClusterID, vmParameters)
	if waitForCompletion {
		err = future.WaitForCompletionRef(ctx, cli.Client)
		if err != nil {
			return "", err
		}

		privateIP, err := e.getPrivateIP(name)
		if err != nil {
			return "", err
		}

		return privateIP, nil
	}

	return "", err
}

// CreateCluster - creates spark clusters
func (e *AzureEnvironment) CreateCluster() (string, error) {
	return e.launchVM(e.ClusterID+"-master", true)
}

// DestroyCluster - destroys spark clusters
func (e *AzureEnvironment) DestroyCluster() error {
	return nil
}

func (e *AzureEnvironment) getClusterNodes() ([]string, error) {
	return []string{}, nil
}
