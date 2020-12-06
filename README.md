# user_data_swap

WARNING: do not deploy this to any account you (or anyone else) actively uses

User data swap is an example of an automated lambda function that runs that swaps out user data on RunInstance events. The original user data script is restored after ours is run.

This exists as an example of how an attacker could semi-covertly backdoor EC2 instances on creation. The API calls will stand out, but from the user's perspective the instance is simply taking longer to start up. This is a well known attack, only change to what I've seen elsewhere is adding Event Bridge and Lambda.

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
