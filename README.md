# user_data_swap

WARNING: do not deploy this to any account you (or anyone else) actively uses

User data swap is an example of an automated lambda function that runs that swaps out user data on RunInstance events. The original user data script is restored after ours is run.

This exists as an example of how an attacker could semi-covertly backdoor EC2 instances on creation.

For more info you can see my post on [Backdooring user data](https://blog.ryanjarv.sh/2020/11/27/backdooring-user-data.html)

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
