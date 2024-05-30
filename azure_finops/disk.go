package azure_finops

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getOrphanedDisks(cred *azidentity.DefaultAzureCredential, subscriptionID string, deleteFlag bool) (int, error) {

	_, diskClient, err := getAzureComputeClient(subscriptionID, cred)
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	pager := diskClient.NewListPager(nil)

	var count int

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to list managed disks: %v", err)
		}

		for _, disk := range page.Value {
			if disk.ManagedBy == nil || *disk.ManagedBy == "" {
				if disk.Properties.DiskSizeGB != nil {
					count++
					fmt.Printf("Orphaned disk found: %s, Size: %d GB\n", *disk.Name, *disk.Properties.DiskSizeGB)
					// Delete the disk if the delete flag is set
					if deleteFlag {
						err := deleteDisk(ctx, diskClient, *disk.ID)
						if err != nil {
							return count, err
						}
						// Log deletion
						fmt.Printf("Deleted disk: %s\n", *disk.Name)
					}
				}
			}
		}
	}

	return count, nil
}
