package aws_finops

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getUnassociatedElasticIPs(sess *session.Session, deleteFlag bool) (int, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeAddressesInput{}

	addresses, err := svc.DescribeAddresses(input)
	if err != nil {
		return 0, err
	}

	var count int

	for _, address := range addresses.Addresses {
		if address.AssociationId == nil {
			count++
			fmt.Printf("Unassociated elastic IP found: %s\n", *address.PublicIp)
			// Release the elastic IP if the delete flag is set
			if deleteFlag {
				_, err := svc.ReleaseAddress(&ec2.ReleaseAddressInput{
					AllocationId: address.AllocationId,
				})
				if err != nil {
					return 0, err
				}
				// Log release
				fmt.Printf("Released elastic IP: %s\n", *address.PublicIp)
			}
		}
	}
	return count, nil
}
