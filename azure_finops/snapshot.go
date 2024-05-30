package azure_finops

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

func getSnapshotSize(snapshot armcompute.Snapshot) float64 {
	if snapshot.Properties != nil && snapshot.Properties.DiskSizeGB != nil {
		return float64(*snapshot.Properties.DiskSizeGB)
	}
	return 0
}

func getOrphanedSnapshots(cred *azidentity.DefaultAzureCredential, subscriptionID string, snapshotRetentionMonths int, deleteFlag bool) (int, error) {
	client, err := armcompute.NewSnapshotsClient(subscriptionID, cred, nil)
	if err != nil {
		return 0, err
	}

	pager := client.NewListPager(nil)
	cutoffDate := time.Now().AddDate(0, -snapshotRetentionMonths, 0)

	var count int
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return 0, fmt.Errorf("failed to list snapshots: %v", err)
		}

		for _, snapshot := range page.Value {
			if snapshot.ManagedBy == nil || *snapshot.ManagedBy == "" {
				if snapshot.Properties.TimeCreated != nil && snapshot.Properties.TimeCreated.Before(cutoffDate) {
					count++
					size := getSnapshotSize(*snapshot)
					fmt.Printf("Old snapshot for more than %d months found: %s, Size: %.2f GB\n", snapshotRetentionMonths, *snapshot.Name, size)
					// Extract resource group from snapshot ID
					resourceGroup := extractResourceGroup(*snapshot.ID)
					// Delete the snapshot if the delete flag is set
					if deleteFlag && resourceGroup != "" {
						poller, err := client.BeginDelete(context.Background(), resourceGroup, *snapshot.Name, nil)
						if err != nil {
							return 0, err
						}
						_, err = poller.PollUntilDone(context.Background(), nil)
						if err != nil {
							return 0, err
						}
						// Log deletion
						fmt.Printf("Deleted snapshot: %s\n", *snapshot.Name)
					}
				}
			}
		}
	}
	return count, nil
}
