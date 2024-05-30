package azure_finops

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// extractResourceGroup extracts the resource group from the given Azure resource ID
func extractResourceGroup(resourceID string) string {
	parts := strings.Split(resourceID, "/")
	for i, part := range parts {
		if strings.EqualFold(part, "resourceGroups") && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// extractResourceName extracts the resource name from the given Azure resource ID
func extractResourceName(resourceID string) string {
	parts := strings.Split(resourceID, "/")
	return parts[len(parts)-1]
}

// extractResourceNameAndGroup extracts the resource name and group from the given Azure resource ID
func extractResourceNameAndGroup(resourceID string) (string, string) {
	parts := strings.Split(resourceID, "/")
	return parts[len(parts)-1], extractResourceGroup(resourceID)
}

// deleteDisk deletes a managed disk
func deleteDisk(ctx context.Context, diskClient *armcompute.DisksClient, diskID string) error {
	diskName, resourceGroup := extractResourceNameAndGroup(diskID)
	fmt.Printf("Deleting disk: %s\n", diskName)
	poller, err := diskClient.BeginDelete(ctx, resourceGroup, diskName, nil)
	if err != nil {
		return fmt.Errorf("failed to begin deleting disk: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to poll disk deletion: %v", err)
	}
	return nil
}

// DeletePublicIP deletes a public IP address
func deletePublicIP(ctx context.Context, ipClient *armnetwork.PublicIPAddressesClient, publicIPID string) error {
	publicIPName, resourceGroup := extractResourceNameAndGroup(publicIPID)
	fmt.Printf("Deleting Public IP: %s\n", publicIPName)
	ipPoller, err := ipClient.BeginDelete(ctx, resourceGroup, publicIPName, nil)
	if err != nil {
		return fmt.Errorf("failed to begin deleting Public IP: %v", err)
	}
	_, err = ipPoller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to poll Public IP deletion: %v", err)
	}
	return nil
}
