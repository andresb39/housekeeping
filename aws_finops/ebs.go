package aws_finops

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getOrphanedVolumes(sess *session.Session, deleteFlag bool) (int, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeVolumesInput{}

	volumes, err := svc.DescribeVolumes(input)
	if err != nil {
		return 0, err
	}

	var count int

	for _, volume := range volumes.Volumes {
		if volume.State != nil && *volume.State == "available" {
			count++
			fmt.Printf("Orphaned volume found: %s, Size: %d GB\n", *volume.VolumeId, *volume.Size)
			// Delete the volume if the delete flag is set
			if deleteFlag {
				_, err := svc.DeleteVolume(&ec2.DeleteVolumeInput{
					VolumeId: volume.VolumeId,
				})
				if err != nil {
					return 0, err
				}
				// Log deletion
				fmt.Printf("Deleted volume: %s\n", *volume.VolumeId)
			}
		}
	}
	return count, nil
}
