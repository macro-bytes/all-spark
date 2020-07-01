package cloud

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// AzureEnvironment interface
type AzureEnvironment struct {
	ClusterID      string
	SubscriptionID string
	Region         string
	ClientID       string
	ClientSecret   string
	Tenant         string
	ResourceGroup  string
	VMNet          string
	VMSubnet       string
	VMSize         string
	ImagePublisher string
	ImageOffer     string
	ImageSku       string
	ImageVersion   string
	WorkerNodes    int64
	EnvParams      []string
}

func (e *AzureEnvironment) getNicClient() (network.InterfacesClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := network.NewInterfacesClient(e.SubscriptionID)
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

func (e *AzureEnvironment) launchVM() (string, error) {
	cli, err := e.getVMClient()
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	vmParameters := compute.VirtualMachine{
		Location: to.StringPtr(e.Region),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypesStandardD8sV3,
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: &compute.ImageReference{
					Publisher: to.StringPtr(e.ImagePublisher),
					Offer:     to.StringPtr(e.ImageOffer),
					Sku:       to.StringPtr(e.ImageSku),
					Version:   to.StringPtr(e.ImageVersion),
				},
			},
			OsProfile: &compute.OSProfile{
				ComputerName:  to.StringPtr(e.ClusterID),
				AdminUsername: to.StringPtr("foobar"),
				AdminPassword: to.StringPtr("W9EHid7dfHTi47Ud"),
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: to.StringPtr("/subscriptions/5bdf2ce9-dc93-4028-9251-ded8d49af5bb/resourceGroups/allspark/providers/Microsoft.Network/networkInterfaces/allsparknetworkinface"),
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
	}
	future, err := cli.CreateOrUpdate(ctx, e.ResourceGroup, e.ClusterID, vmParameters)
	if err != nil {
		return "", err
	}

	return "", future.WaitForCompletionRef(ctx, cli.Client)
}

// CreateCluster - creates spark clusters
func (e *AzureEnvironment) CreateCluster() (string, error) {
	return "", nil
}

// DestroyCluster - destroys spark clusters
func (e *AzureEnvironment) DestroyCluster() error {
	return nil
}

func (e *AzureEnvironment) getClusterNodes() ([]string, error) {
	return []string{}, nil
}
