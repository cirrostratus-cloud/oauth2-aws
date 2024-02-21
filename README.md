# Cirrostratus Cloud - OAuth2 AWS

## Requirements

First create an dotenv file (`.env`), like this:

```shell
AWS_STAGE=prod
LOG_LEVEL=INFO
CIRROSTRATUS_OAUTH2_MODULE_NAME=cirrostratus-oauth2
CIRROSTRATUS_OUTH2_USER_TABLE=users
AWS_DEFAULT_REGION=us-west-1
AWS_REGION=us-west-1
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
USER_MIN_PASSWORD_LENGTH=8
USER_UPPER_CASE_REQUIRED=true
USER_LOWER_CASE_REQUIRED=true
USER_NUMBER_REQUIRED=true
USER_SPECIAL_CHARACTER_REQUIRED=true
```

## Run locally

```shell
task client:serve
```

```shell
task user:serve
```

## Deploy

```shell
. .env
task deploy
```