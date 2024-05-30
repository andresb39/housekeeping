package aws_finops

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func RunAWS(shutdownMonths, snapshotRetentionMonths int, region string, deleteFlag bool) {
	// AWS Session
	awsSess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Error creating AWS session: %v", err)
	}

	housekeeping := housekeeping(awsSess, shutdownMonths, snapshotRetentionMonths, deleteFlag)
	if deleteFlag {
		fmt.Printf("Total orphan resources to be deleted: %d\n", housekeeping)
	} else {
		fmt.Printf("Total orphan resources deleted: %d\n", housekeeping)
	}
}

func housekeeping(sess *session.Session, shutdownMonths int, snapshotRetentionMonths int, deleteFlag bool) int {
	// AWS Snapshots
	snapshots, err := getOrphanedSnapshots(sess, snapshotRetentionMonths, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting AWS orphaned snapshots: %v", snapshots)
	}

	// AWS Volumes
	vols, err := getOrphanedVolumes(sess, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting AWS orphaned volumes: %v", vols)
	}

	// AWS Elastic IPs
	elasticIps, err := getUnassociatedElasticIPs(sess, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting AWS unassociated elastic IPs: %v", elasticIps)
	}

	// AWS Stopped Instances
	instances, err := getStoppedInstances(sess, shutdownMonths, deleteFlag)
	if err != nil {
		log.Fatalf("Error getting AWS instances stopped for more than %d months: %v", shutdownMonths, instances)
	}
	return snapshots + vols + elasticIps + instances
}
