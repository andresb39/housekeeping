package main

import (
	"flag"
	"log"

	"github.com/housekeeping/aws_finops"
	"github.com/housekeeping/azure_finops"
)

func main() {
	// Command line flags
	provider := flag.String("provider", "", "Cloud provider: aws or azure")
	subscriptionID := flag.String("subscriptionID", "", "Azure subscription ID (optional, required for a single subscription)")
	shutdownMonths := flag.Int("shutdownMonths", 6, "Number of months before considering instances as orphaned")
	snapshotRetentionMonths := flag.Int("snapshotRetentionMonths", 6, "Number of months to retain snapshots")
	region := flag.String("region", "us-east-1", "AWS region")
	deleteFlag := flag.Bool("delete", false, "Delete identified orphaned resources")

	flag.Parse()

	if *provider == "" {
		log.Fatal("You must specify a cloud provider: aws or azure")
	}

	switch *provider {
	case "aws":
		aws_finops.RunAWS(*shutdownMonths, *snapshotRetentionMonths, *region, *deleteFlag)
	case "azure":
		azure_finops.RunAzure(*subscriptionID, *shutdownMonths, *snapshotRetentionMonths, *deleteFlag)
	default:
		log.Fatalf("Unsupported cloud provider: %s. Please specify aws or azure.", *provider)
	}
}
