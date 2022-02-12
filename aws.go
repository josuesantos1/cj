package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"fmt"
)

func Instance(command string) {
	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("us-east-2"),
		},
	}))

	// Create new EC2 client
	svc := ec2.New(sess)

	if os.Args[1] == "START" { // START
		startInstance(os.Args[2], svc)
	} else if os.Args[1] == "STOP" { // Turn instances off
		instanceStop(os.Args[2], svc)
	} else if os.Args[1] == "CREATE" { // Create instance
		createInstance(os.Args[2], svc)
	} else if os.Args[1] == "DESTROY" { // Destroy instance
		destroyInstance(os.Args[2], svc)
	} else if os.Args[1] == "LIST" { // List instances
		getInstances(svc)
	} else {
		fmt.Println("Invalid command")
	}
}

func createInstance(ami string, svc *ec2.EC2) {
	// create instance with ami
	result, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(ami),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Success", result.Instances)
	}

	// get instance id
	instanceID := *result.Instances[0].InstanceId

	fmt.Println("Instance ID:", instanceID)
}

func destroyInstance(instanceID string, svc *ec2.EC2) {
	result, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", result.TerminatingInstances)
}

func startInstance(instanceID string, svc *ec2.EC2) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StartInstances(input)

	awsErr, ok := err.(awserr.Error)

	// If the error code is `DryRunOperation` it means we have the necessary
	// permissions to Start this instance
	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to start the instances
		input.DryRun = aws.Bool(false)
		result, err = svc.StartInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("Success", result.StartingInstances)
		}
	} else { // This could be due to a lack of permissions
		fmt.Println("Error", err)
	}
}

func instanceStop(instanceID string, svc *ec2.EC2) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(true),
	}
	result, err := svc.StartInstances(input)

	fmt.Println(result)

	awsErr, ok := err.(awserr.Error)

	// If the error code is `DryRunOperation` it means we have the necessary
	// permissions to Start this instance
	if ok && awsErr.Code() == "DryRunOperation" {
		// Let's now set dry run to be false. This will allow us to start the instances
		input.DryRun = aws.Bool(false)
		result, err = svc.StartInstances(input)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Println("Success", result.StartingInstances)
		}
	} else { // This could be due to a lack of permissions
		fmt.Println("Error", err)
	}
}

func getInstances(svc *ec2.EC2) {
	// List all instances in this region
	result, err := svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	instances := []string{}

	for _, instance := range result.Reservations {
		instances = append(instances, aws.StringValue(instance.Instances[0].InstanceId))
	}

	fmt.Println(instances)
}

// utils

func getSecurityGroup(name string, svc *ec2.EC2) (string, error) {
	// Create the input for DescribeSecurityGroups
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("group-name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}

	// Describe the security groups
	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}

	GroupId := aws.StringValue(result.SecurityGroups[0].GroupId)

	return GroupId, nil
}
