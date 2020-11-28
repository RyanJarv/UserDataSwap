package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"time"
)

func StartInstance(ctx context.Context, instance string) {
	fmt.Printf("[INFO] Starting Instance: %s\n", instance)
	if _, err := client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []*string{&instance},
	}); err != nil {
		panic(err)
	}
}

func StopInstance(ctx context.Context, instance string) {
	fmt.Printf("[INFO] Stoppping Instance: %s\n", instance)
	if _, err := client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: []*string{&instance},
		Force:       aws.Bool(true),
	}); err != nil {
		panic(err)
	}
}

func WaitForInstance(ctx context.Context, instance string, instanceState types.InstanceStateName) {
	for {
		var status types.InstanceStateName
		if out, err := client.DescribeInstances(ctx, &(ec2.DescribeInstancesInput{
			InstanceIds: []*string{&instance},
		})); err != nil {
			panic(err)
		} else {
			status = out.Reservations[0].Instances[0].State.Name
		}

		fmt.Printf("[DEBUG] Status: %v, instanceState %v\n", status, instanceState)
		if status == instanceState {
			return
		} else {
			fmt.Printf("[DEBUG] Instance in %s, desired state is %s, sleeping and will try again\n", status, instanceState)
			time.Sleep(time.Second * 5)
		}
	}
}

func GetUserData(ctx context.Context, instanceId *string) (userData UserData) {
	out, err := client.DescribeInstanceAttribute(ctx, &ec2.DescribeInstanceAttributeInput{
		Attribute:  "userData",
		InstanceId: instanceId,
	})

	var encoded string
	if err != nil {
		panic(err)
	} else {
		if out.UserData != nil {
			encoded = *out.UserData.Value
		} else {
			encoded = ""
		}
	}

	if b, err := base64.StdEncoding.DecodeString(encoded); err != nil {
		fmt.Println("decode error:", err)
		return
	} else {
		userData.Value = string(b)
	}

	return userData
}

// ModifyUserData will modify user data.
// It will wait for running state, stop the instance, modify the user data, start it and wait for it to enter
// the running state.
func ModifyUserData(ctx context.Context, instance string, userData string) {
	var state string
	if out, err := client.DescribeInstances(ctx, &(ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instance},
	})); err != nil {
		panic(err)
	} else {
		state = string(out.Reservations[0].Instances[0].State.Name)
	}

	fmt.Printf("[DEBUG] Modifing instance data for %s %s to '%s'\n", state, instance, userData)


	if state == "pending" {
		fmt.Printf("[DEBUG] Waiting on instance %s to enter running state\n", instance)
		WaitForInstance(ctx, instance, types.InstanceStateNameRunning)
		StopInstance(ctx, instance)
	} else if state == "running" {
		fmt.Printf("[INFO] Stopping instance %s\n", instance)
		StopInstance(ctx, instance)
	} else if state == "stopping" {
		fmt.Printf("[INFO] Found instance %s in stopping state already\n", instance)
	}

	fmt.Printf("[DEBUG] Waiting on instance %s to enter stopped state\n", instance)
	WaitForInstance(ctx, instance, types.InstanceStateNameStopped)

	fmt.Printf("[INFO] Modifying instance %s user data\n", instance)
	fmt.Printf("[DEBUG] User data: %s \n", userData)
	if _, err := client.ModifyInstanceAttribute(ctx, &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instance),
		UserData: &types.BlobAttributeValue{
			Value: []byte(userData),
		},
	}); err != nil {
		panic(err)
	}

}
