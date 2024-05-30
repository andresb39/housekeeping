package azure_finops

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func RunAzure(subscriptionID string, shutdownMonths, snapshotRetentionMonths int, deleteFlag bool) {
	// Azure Authentication
	azureClient, err := getAzureDetailsClient()
	if err != nil {
		log.Fatalf("Error creating Azure authorizer: %v", err)
	}

	if subscriptionID != "" {
		housekeeping := housekeeping(azureClient, subscriptionID, shutdownMonths, snapshotRetentionMonths, deleteFlag)
		if deleteFlag {
			fmt.Printf("Total orphan resources to be deleted: %d\n", housekeeping)
		} else {
			fmt.Printf("Total orphan resources deleted: %d\n", housekeeping)
		}
	} else {
		subscriptions, err := listSubscriptions(azureClient)
		if err != nil {
			log.Fatalf("Error listing subscriptions: %v", err)
		}

		for _, sub := range subscriptions {
			housekeeping := housekeeping(azureClient, *sub.SubscriptionID, shutdownMonths, snapshotRetentionMonths, deleteFlag)
			if deleteFlag {
				fmt.Printf("Total orphan resources to be deleted for subscription %s is: %d\n", *sub.DisplayName, housekeeping)
			} else {
				fmt.Printf("Total orphan resources deleted for subscription %s is: %d\n", *sub.DisplayName, housekeeping)
			}
		}
	}
}

func housekeeping(cred *azidentity.DefaultAzureCredential, subscriptionID string, shutdownMonths int, snapshotRetentionMonths int, deleteFlag bool) int {
	// Azure Snapshots
	snapshots, err := getOrphanedSnapshots(cred, subscriptionID, snapshotRetentionMonths, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting Azure orphaned snapshots: %v", snapshots)
	}

	// Azure Managed Disks
	disk, err := getOrphanedDisks(cred, subscriptionID, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting Azure orphaned managed disks: %v", disk)
	}

	// Azure Public IPs
	publicIP, err := getUnassociatedPublicIPs(cred, subscriptionID, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting Azure unassociated public IPs: %v", publicIP)
	}

	// Azure Stopped VMs
	vms, err := getStoppedVMs(cred, subscriptionID, shutdownMonths, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting Azure VMs stopped for more than %d months: %v", shutdownMonths, vms)
	}

	return disk + publicIP + snapshots + vms
}
