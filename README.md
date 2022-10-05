# Vex (Vexillology)

A service to manage feature flags and remote configuration.

Accounts can contain multiple projects which consist of a set of configuration flags.


![diagram](vex.png)

## Accounts
`curl -X GET -H 'Authorization: Bearer <token here>' /accounts/{accountId}`

```json
{
  "data": {
    "id": "cb6049d9-7720-4442-89be-f9500c72a73b",
    "name": "production",
    "description": "the production account",
    "created_on": "2022-08-15T02:57:56.753876Z",
    "modified_on": "2022-08-15T02:57:56.753876Z"
  },
  "success": true,
  "errors": []
}
```

## Tokens
A token is used to manage account objects (read_only = false) or can be set to read only for just reading the project flags
from the CDN (read_only = true).

`curl -X POST -H 'Authorization Bearer <token here>` /accounts/{accountId}/tokens?readOnly=true`
```json
{
    "data": {
        "id": "8eafc880-9493-4d00-b9e7-389e9ce989fd",
        "accountId": "cb6049d9-7720-4442-89be-f9500c72a73b",
        "token": "0be2784d2c16943be7295b8dedc4561b",  <--- only shown once
        "read_only": true,
        "created_on": "2022-09-22T03:26:25.841193Z",
        "modified_on": "2022-09-22T03:26:25.841193Z"
    },
    "success": true,
    "errors": []
}
```
`curl -X GET -H 'Authorization: Bearer <token here>` /accounts/{acountId}/tokens`
```json
{
    "data": [
        {
            "id": "8eafc880-9493-4d00-b9e7-389e9ce989fd",
            "accountId": "cb6049d9-7720-4442-89be-f9500c72a73b",
            "read_only": true,
            "created_on": "2022-09-22T03:26:25.841193Z",
            "modified_on": "2022-09-22T03:26:25.841193Z"
        }
    ],
    "success": true,
    "errors": []
}
```

Reroll an existing token to get a new value without creating a new token.
`curl -X PUT -H 'Authorization: Bearer <token here>' /accounts/{accountId/tokens/{tokenId}`
```json
{
    "data": {
        "id": "8eafc880-9493-4d00-b9e7-389e9ce989fd",
        "accountId": "cb6049d9-7720-4442-89be-f9500c72a73b",
        "token": "32deaa8426569e594150329ca26b1dd2",   <--- rerolled token value
        "read_only": true,
        "created_on": "2022-09-22T03:26:25.841193Z",
        "modified_on": "2022-09-22T03:26:25.841193Z"
    },
    "success": true,
    "errors": []
}
```

## Projects

A project is a set of configuration flags.

`curl -X GET -H 'Authorization: Bearer <token here>' /accounts/{accountId}/projects/{projectId}`
```json
{
  "data": [
    {
      "id": "ed7f9f1c-4416-4f2f-8ff1-cfe10c8d14e0",
      "account_id": "cb6049d9-7720-4442-89be-f9500c72a73b",
      "name": "proj 1",
      "description": "project one",
      "created_on": "2022-08-15T02:58:38.846618Z",
      "modified_on": "2022-08-15T02:58:38.846618Z"
    },
    {
      "id": "e746f8d8-3f46-462d-8bad-c980e9c01152",
      "account_id": "cb6049d9-7720-4442-89be-f9500c72a73b",
      "name": "proj 2",
      "description": "project two",
      "created_on": "2022-09-12T22:22:19.257571Z",
      "modified_on": "2022-09-12T22:22:19.257571Z"
    }
  ],
  "success": true,
  "errors": []
}
```

## Flags
Flags hold the configuration values for a project. They can be of types `BOOLEAN`, `NUMBER`, and `STRING`.

Flags store their raw value as strings with an enum that specifies their type. Each SDK can decide
how to parse the flag value in their own language.

`curl -X GET -H 'Authorization: Bearer <token here>' /accounts/{accountId}/projects/{projectId}/flags/{flagId}`
```json
{
    "data": [
        {
            "id": "00489c7e-0bf1-4636-865e-294079234658",
            "project_id": "ed7f9f1c-4416-4f2f-8ff1-cfe10c8d14e0",
            "account_id": "cb6049d9-7720-4442-89be-f9500c72a73b",
            "created_on": "2022-08-15T03:00:20.973395Z",
            "modified_on": "2022-08-15T03:00:20.973395Z",
            "key": "feature1",
            "type": "NUMBER",
            "value": "123.45"
        },
        {
            "id": "53280952-2048-4961-bbb1-3c388973b667",
            "project_id": "ed7f9f1c-4416-4f2f-8ff1-cfe10c8d14e0",
            "account_id": "cb6049d9-7720-4442-89be-f9500c72a73b",
            "created_on": "2022-08-15T03:00:41.128669Z",
            "modified_on": "2022-08-15T03:00:41.128669Z",
            "key": "feature2",
            "type": "BOOLEAN",
            "value": "true"
        },
        {
            "id": "78ac98bb-afbc-4c72-8c58-e6d4f1df276b",
            "project_id": "ed7f9f1c-4416-4f2f-8ff1-cfe10c8d14e0",
            "account_id": "cb6049d9-7720-4442-89be-f9500c72a73b",
            "created_on": "2022-09-09T01:32:55.941958Z",
            "modified_on": "2022-09-09T01:32:55.941958Z",
            "key": "feature3",
            "type": "STRING",
            "value": "string value"
        }
    ],
    "success": true,
    "errors": []
}
```

## CDN 

When projects are modified the configuration is rendered and provisioned in the Cloudflare CDN Worker.

This worker allows for fast and distributed access to the project configuration.

`curl -X GET -H 'Authorization: Bearer <token here>' /{projectId}`
```json
{
  "feature1": {
    "type": "BOOLEAN",
    "value": "true"
  },
  "feature2": {
    "type": "NUMBER",
    "value": "123.45"
  },
  "feature3": {
    "type": "STRING",
    "value": "some text"
  }
}
```

### OpenAPI 3 
An OpenAPI spec that describes all endpoints is located at `./openapi/openapi.yaml`

## Client Libraries

Go - https://github.com/broswen/vex-go

## Terraform

terraform-provider-vex - https://github.com/broswen/terraform-provider-vex

### Todo
- [x] handle postgres errors and wrap in custom store errors (ongoing)
- [x] provision account tokens to cloudflare kv
  - [x] implement worker token authentication
- [x] handle docker-compose initialize local postgres with schema
  - still need to handle multiple migrations
- [x] handle local provisioning for dockerfile, flag to skip api calls?
- [ ] add mocks and tests with testify
- [x] add created_on and modified_on fields to all resources
- [ ] incremental config builds
  - store prerendered config in postgres, parse and insert/update flags as needed
