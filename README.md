# Cirrostratus Cloud - OAuth2 AWS

## Deploy

First create an dotenv file (`.env`), like this:

```shell
export AWS_DEFAULT_REGION="us-west-1"
export AWS_ACCESS_KEY_ID="your-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
```

```shell
. .env
task deploy
```