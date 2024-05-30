package azure_finops

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getUnassociatedPublicIPs(cred *azidentity.DefaultAzureCredential, subscriptionID string, deleteFlag bool) (int, error) {
	client, _, err := getAzureNetworkClient(subscriptionID, cred)
	if err != nil {
		return 0, err
	}

	pager := client.NewListAllPager(nil)
	ctx := context.Background()

	var count int

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to list public IP addresses: %v", err)
		}

		for _, ip := range page.PublicIPAddressListResult.Value {
			if ip.Properties == nil || ip.Properties.IPConfiguration == nil {
				count++
				fmt.Printf("Unassociated public IP found: %s\n", *ip.Name)
				// Delete the Public IP if the delete flag is set
				if deleteFlag {
					err := deletePublicIP(ctx, client, *ip.ID)
					if err != nil {
						return 0, err
					}
					// Log deletion
					fmt.Printf("Deleted Public IP: %s\n", *ip.Name)
				}
			}
		}
	}
	return count, nil
}
