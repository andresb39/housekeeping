package aws_finops

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetOrphanedSnapshotsAWS calculates the number of orphaned snapshots and optionally deletes them
func getOrphanedSnapshots(sess *session.Session, snapshotRetentionMonths int, deleteFlag bool) (int, error) {
	svc := ec2.New(sess)

	// Get only snapshots owned by me
	input := &ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{aws.String("self")},
	}

	snapshots, err := svc.DescribeSnapshots(input)
	if err != nil {
		return 0, err
	}

	retentionTime := time.Now().AddDate(0, -snapshotRetentionMonths, 0)
	var count int

	for _, snapshot := range snapshots.Snapshots {
		// Calculate the age of the snapshot
		if snapshot.StartTime.Before(retentionTime) {
			if isSnapshotOrphaned(svc, snapshot) {
				// Calculate the cost
				if snapshot.VolumeSize != nil {
					count++
					fmt.Printf("Old orphaned snapshot for more than %d months found: %s, Size: %d GB\n", snapshotRetentionMonths, *snapshot.SnapshotId, *snapshot.VolumeSize)
				}
				// Delete the snapshot if the delete flag is set
				if deleteFlag {
					_, err := svc.DeleteSnapshot(&ec2.DeleteSnapshotInput{
						SnapshotId: snapshot.SnapshotId,
					})
					if err != nil {
						log.Printf("Error deleting snapshot %s: %v", *snapshot.SnapshotId, err)
						continue
					}
					// Log deletion
					fmt.Printf("Deleted snapshot: %s\n", *snapshot.SnapshotId)
				}
			}
		}
	}

	return count, nil
}

// isSnapshotOrphaned checks if a snapshot is orphaned by ensuring it is not associated with any AMIs
func isSnapshotOrphaned(svc *ec2.EC2, snapshot *ec2.Snapshot) bool {
	// Describe images to check if the snapshot is associated with any AMIs
	imageInput := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("block-device-mapping.snapshot-id"),
				Values: []*string{snapshot.SnapshotId},
			},
		},
	}

	images, err := svc.DescribeImages(imageInput)
	if err != nil {
		log.Printf("Error describing images for snapshot %s: %v", *snapshot.SnapshotId, err)
		return false
	}
	if len(images.Images) > 0 {
		return false // Snapshot is associated with an AMI
	}

	return true // Snapshot is orphaned
}
