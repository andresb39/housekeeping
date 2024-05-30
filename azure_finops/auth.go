package azure_finops

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

// GetAzureUsageDetailsClient creates a new UsageDetailsClient
func getAzureDetailsClient() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// ListSubscriptions lists all subscriptions in the tenant
func listSubscriptions(cred *azidentity.DefaultAzureCredential) ([]*armsubscriptions.Subscription, error) {
	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return nil, err
	}

	pager := client.NewListPager(nil)
	var subscriptions []*armsubscriptions.Subscription

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, page.SubscriptionListResult.Value...)
	}

	return subscriptions, nil
}

// Get client for Azure network
func getAzureNetworkClient(subscriptionID string, cred *azidentity.DefaultAzureCredential) (*armnetwork.PublicIPAddressesClient, *armnetwork.InterfacesClient, error) {
	ipClient, err := armnetwork.NewPublicIPAddressesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, nil, err
	}
	nicClient, err := armnetwork.NewInterfacesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, nil, err
	}
	return ipClient, nicClient, nil
}

// Get client for Azure compute
func getAzureComputeClient(subscriptionID string, cred *azidentity.DefaultAzureCredential) (*armcompute.VirtualMachinesClient, *armcompute.DisksClient, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, nil, err
	}
	diskClient, err := armcompute.NewDisksClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, nil, err
	}
	return vmClient, diskClient, nil
}
