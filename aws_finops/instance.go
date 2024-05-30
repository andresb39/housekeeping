package aws_finops

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getStoppedInstances(sess *session.Session, shutdownMonths int, deleteFlag bool) (int, error) {
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return 0, err
	}

	cutoffDate := time.Now().AddDate(0, -shutdownMonths, 0)
	var count int

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if *instance.State.Name == "stopped" && instance.StateTransitionReason != nil {
				launchTime := *instance.LaunchTime
				if launchTime.Before(cutoffDate) {
					count++
					fmt.Printf("Instance stopped for more than %d months found: %s\n", shutdownMonths, *instance.InstanceId)
					if deleteFlag {
						_, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
							InstanceIds: []*string{instance.InstanceId},
						})
						if err != nil {
							return 0, err
						}
						// Log termination
						fmt.Printf("Terminated instance: %s\n", *instance.InstanceId)
					}
				}
			}
		}
	}

	return count, nil
}
