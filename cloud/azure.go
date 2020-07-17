package cloud

import (
	"allspark/daemon"
	"allspark/logger"
	"container/list"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/storage/mgmt/storage"

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
	VMSize              compute.VirtualMachineSizeTypes
	ImageStorageAccount string
	DataStorageAccount  string
	ImageContainer      string
	ImageBlob           string
	WorkerNodes         int64
	EnvParams           []string
}

func (e *AzureEnvironment) getStorageClient() (storage.AccountsClient, error) {
	authConfig := auth.NewClientCredentialsConfig(e.ClientID, e.ClientSecret, e.Tenant)
	client := storage.NewAccountsClient(e.SubscriptionID)
	authorizer, err := authConfig.Authorizer()
	client.Authorizer = authorizer
	return client, err
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

func (e *AzureEnvironment) getImageURI() string {
	return "https://" + e.ImageStorageAccount +
		".blob.core.windows.net/" + e.ImageContainer +
		"/" + e.ImageBlob
}

func (e *AzureEnvironment) getSubnet(ctx context.Context, vnetName string, subnetName string) (network.Subnet, error) {
	subnetsClient, _ := e.getSubnetClient()
	return subnetsClient.Get(ctx, e.ResourceGroup, vnetName, subnetName, "")
}

func (e *AzureEnvironment) getPrimaryStorageKey() (string, error) {
	client, _ := e.getStorageClient()

	result, err := client.ListKeys(context.Background(), e.ResourceGroup, e.DataStorageAccount)
	if err != nil {
		return "", err
	}

	keyList := *(result.Keys)
	if len(keyList) == 0 {
		return "", errors.New("storage account contains no keys")
	}

	primaryKey := keyList[0].Value
	return *primaryKey, nil
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

	_, err = cli.CreateOrUpdate(context.Background(), e.ResourceGroup, name, nicParams)
	return nicPath, err
}

func (e *AzureEnvironment) deleteNIC(name string) error {
	cli, err := e.getNicClient()
	if err != nil {
		return err
	}
	_, err = cli.Delete(context.Background(), e.ResourceGroup, name)
	return err
}

func (e *AzureEnvironment) listNICs() ([]string, error) {
	cli, err := e.getNicClient()
	if err != nil {
		return []string{}, err
	}

	result, err := cli.List(context.Background(), e.ResourceGroup)
	if err != nil {
		return []string{}, err
	}

	items := make([]string, 0)
	for _, el := range result.Values() {
		if strings.Contains(*el.Name, e.ClusterID) {
			items = append(items, *el.Name)
		}
	}

	return items, err
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
				SourceURI:        to.StringPtr(e.getImageURI()),
				StorageAccountID: to.StringPtr(storageAccountID),
			},
		},
	}

	_, err = cli.CreateOrUpdate(context.Background(), e.ResourceGroup, name, disk)
	return diskPath, err
}

func (e *AzureEnvironment) deleteDisk(name string) error {
	cli, err := e.getDiskClient()
	if err != nil {
		return err
	}
	_, err = cli.Delete(context.Background(), e.ResourceGroup, name)
	return err
}

func (e *AzureEnvironment) listDisks() (*list.List, error) {
	cli, err := e.getDiskClient()
	if err != nil {
		return nil, err
	}

	result, err := cli.ListByResourceGroup(context.Background(), e.ResourceGroup)
	if err != nil {
		return nil, err
	}

	items := list.New()

	for _, el := range result.Values() {
		if strings.Contains(*el.Name, e.ClusterID) {
			items.PushBack(*el.Name)
		}
	}

	return items, err
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

func (e *AzureEnvironment) createVM(name string, tags map[string]*string,
	waitForCompletion bool) (string, error) {

	cli, err := e.getVMClient()

	if err != nil {
		return "", err
	}
	ctx := context.Background()

	nic, err := e.createNIC(name)
	if err != nil {
		return "", err
	}

	privateIP, err := e.getPrivateIP(name)
	if err != nil {
		return "", err
	}

	disk, err := e.createDisk(name)
	if err != nil {
		return "", err
	}

	vmParameters := compute.VirtualMachine{
		Location: to.StringPtr(e.Region),
		Tags:     tags,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: e.VMSize,
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

	future, err := cli.CreateOrUpdate(ctx, e.ResourceGroup, name, vmParameters)
	if err != nil {
		return "", err
	}

	if waitForCompletion {
		err = future.WaitForCompletionRef(ctx, cli.Client)
		if err != nil {
			return "", err
		}

	}

	return privateIP, nil
}

func (e *AzureEnvironment) deleteVM(name string) {
	cli, err := e.getVMClient()
	if err != nil {
		logger.GetError().Println(err)
	}

	future, err := cli.Delete(context.Background(), e.ResourceGroup, name)
	if err != nil {
		logger.GetError().Println(err)
	}

	err = future.WaitForCompletionRef(context.Background(), cli.Client)
	if err != nil {
		logger.GetError().Println(err)
	}

	err = e.deleteNIC(name)
	if err != nil {
		logger.GetError().Println(err)
	}

	err = e.deleteDisk(name)
	if err != nil {
		logger.GetError().Println(err)
	}
}

func (e *AzureEnvironment) launchMaster() (string, error) {
	tags := make(map[string]*string)

	tags["EXPECTED_WORKERS"] = to.StringPtr(strconv.FormatInt(e.WorkerNodes, 10))
	tags["SPARK_WORKER_PORT"] = to.StringPtr(strconv.FormatInt(sparkWorkerPort, 10))
	tags["CLUSTER_ID"] = to.StringPtr(e.ClusterID)
	tags["ALLSPARK_CALLBACK"] = to.StringPtr(daemon.GetAllSparkConfig().CallbackURL)

	if len(e.DataStorageAccount) > 0 {
		tags["DATA_STORAGE_ACCOUNT"] = to.StringPtr(e.DataStorageAccount)
		storageKey, err := e.getPrimaryStorageKey()
		if err != nil {
			return "", err
		}
		tags["DATA_STORAGE_KEY"] = to.StringPtr(storageKey)
	}

	for _, el := range e.EnvParams {
		buff := strings.SplitN(el, "=", 2)
		tags[buff[0]] = to.StringPtr(buff[1])
	}

	return e.createVM(e.ClusterID+"-master", tags, false)
}

func (e *AzureEnvironment) launchWorkers(masterIP string) error {
	tags := make(map[string]*string)

	tags["MASTER_IP"] = to.StringPtr(masterIP)
	tags["SPARK_WORKER_PORT"] = to.StringPtr(strconv.FormatInt(sparkWorkerPort, 10))

	for _, el := range e.EnvParams {
		buff := strings.SplitN(el, "=", 2)
		tags[buff[0]] = to.StringPtr(buff[1])
	}

	var i int64
	for i = 0; i < e.WorkerNodes; i++ {
		_, err := e.createVM(e.ClusterID+"-worker-"+strconv.FormatInt(i, 10), tags, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateCluster - creates spark clusters
func (e *AzureEnvironment) CreateCluster() (string, error) {
	masterIP, err := e.launchMaster()
	if err != nil {
		return "", err
	}

	return "", e.launchWorkers(masterIP)
}

// DestroyCluster - destroys spark clusters
func (e *AzureEnvironment) DestroyCluster() error {
	vms, err := e.getClusterNodes()
	if err != nil {
		return err
	}

	for _, el := range vms {
		go e.deleteVM(el)
	}

	return nil
}

func (e *AzureEnvironment) getClusterNodes() ([]string, error) {
	cli, err := e.getVMClient()
	if err != nil {
		return nil, err
	}

	result, err := cli.List(context.Background(), e.ResourceGroup)
	if err != nil {
		return nil, err
	}

	items := make([]string, 0)

	for _, el := range result.Values() {
		if strings.Contains(*el.Name, e.ClusterID) {
			items = append(items, *el.Name)
		}
	}

	return items, err
}
