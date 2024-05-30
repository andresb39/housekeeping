package azure_finops

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func getStoppedVMs(cred *azidentity.DefaultAzureCredential, subscriptionID string, shutdownMonths int, deleteFlag bool) (int, error) {

	vmClient, diskClient, err := getAzureComputeClient(subscriptionID, cred)
	if err != nil {
		return 0, err
	}

	ipClient, nicClient, err := getAzureNetworkClient(subscriptionID, cred)
	if err != nil {
		return 0, err
	}

	pager := vmClient.NewListAllPager(nil)
	cutoffDate := time.Now().AddDate(0, -shutdownMonths, 0)
	ctx := context.Background()

	var count int

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to list virtual machines: %v", err)
		}

		for _, vm := range page.VirtualMachineListResult.Value {
			if vm.ID != nil {
				resourceGroup := extractResourceGroup(*vm.ID)
				instanceView, err := vmClient.InstanceView(ctx, resourceGroup, *vm.Name, nil)
				if err != nil {
					fmt.Printf("failed to get instance view for VM: %s, error: %v\n", *vm.Name, err)
					continue
				}

				for _, status := range instanceView.Statuses {
					if status.Code != nil && (strings.HasSuffix(*status.Code, "stopped") || strings.HasSuffix(*status.Code, "deallocated")) {
						if shutdownMonths == 0 || (vm.Properties.TimeCreated != nil && vm.Properties.TimeCreated.Before(cutoffDate)) {
							count++
							fmt.Printf("Stopped VM for more than %d months found: %s\n", shutdownMonths, *vm.Name)
							// Delete the VM and associated resources if the delete flag is set
							if deleteFlag && resourceGroup != "" {
								err = deleteVMAndResources(ctx, vmClient, diskClient, nicClient, ipClient, resourceGroup, *vm.Name, *vm)
								if err != nil {
									return 0, err
								}
								// Log deletion
								fmt.Printf("Deleted VM and associated resources: %s\n", *vm.Name)
							}
						}
					}
				}
			}
		}
	}

	return count, nil
}

// deleteVMAndResources deletes a VM and its associated resources
func deleteVMAndResources(ctx context.Context, vmClient *armcompute.VirtualMachinesClient, diskClient *armcompute.DisksClient, nicClient *armnetwork.InterfacesClient, ipClient *armnetwork.PublicIPAddressesClient, resourceGroup, vmName string, vm armcompute.VirtualMachine) error {
	// Delete the VM
	fmt.Printf("Deleting VM: %s\n", vmName)
	poller, err := vmClient.BeginDelete(ctx, resourceGroup, vmName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete VM: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to poll VM deletion: %v", err)
	}

	// Delete associated disks
	if vm.Properties.StorageProfile != nil {
		if osDisk := vm.Properties.StorageProfile.OSDisk; osDisk != nil && osDisk.ManagedDisk != nil && osDisk.ManagedDisk.ID != nil {
			err = deleteDisk(ctx, diskClient, *osDisk.ManagedDisk.ID)
			if err != nil {
				return fmt.Errorf("failed to delete OS disk: %v", err)
			}
		}
		for _, dataDisk := range vm.Properties.StorageProfile.DataDisks {
			if dataDisk.ManagedDisk != nil && dataDisk.ManagedDisk.ID != nil {
				err = deleteDisk(ctx, diskClient, *dataDisk.ManagedDisk.ID)
				if err != nil {
					return fmt.Errorf("failed to delete data disk: %v", err)
				}
			}
		}
	}

	// Delete associated network interfaces and public IPs
	if vm.Properties.NetworkProfile != nil {
		for _, nic := range vm.Properties.NetworkProfile.NetworkInterfaces {
			if nic.ID != nil {
				err = deleteNICAndPublicIP(ctx, nicClient, ipClient, *nic.ID)
				if err != nil {
					return fmt.Errorf("failed to delete NIC and Public IP: %v", err)
				}
			}
		}
	}

	return nil
}

// / deleteNICAndPublicIP deletes a network interface and its associated public IP address
func deleteNICAndPublicIP(ctx context.Context, nicClient *armnetwork.InterfacesClient, ipClient *armnetwork.PublicIPAddressesClient, nicID string) error {
	nicName, resourceGroup := extractResourceNameAndGroup(nicID)
	fmt.Printf("Deleting NIC: %s\n", nicName)

	// Get NIC details to find associated Public IP
	nic, err := nicClient.Get(ctx, resourceGroup, nicName, nil)
	if err != nil {
		return fmt.Errorf("failed to get NIC details: %v", err)
	}

	// Delete the NIC
	poller, err := nicClient.BeginDelete(ctx, resourceGroup, nicName, nil)
	if err != nil {
		return fmt.Errorf("failed to begin deleting NIC: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to poll NIC deletion: %v", err)
	}

	// Delete associated Public IP addresses
	if nic.Properties.IPConfigurations != nil {
		for _, ipConfig := range nic.Properties.IPConfigurations {
			if ipConfig.Properties.PublicIPAddress != nil && ipConfig.Properties.PublicIPAddress.ID != nil {
				err := deletePublicIP(ctx, ipClient, *ipConfig.Properties.PublicIPAddress.ID)
				if err != nil {
					return fmt.Errorf("failed to delete Public IP: %v", err)
				}
			}
		}
	}

	return nil
}
