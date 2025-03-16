# User Service

## About
A sample user management service that demonstrates JWT based authentication mechanism when making an API request to the server. In addition, the service also demonstrates a Role-based Access Control (RBAC) mechanism during the request handling.

## Usage

### API
The OpenAPI spec is available at https://github.com/amarjeet000/user-mgmnt-service/blob/main/openapi.yml

The service exposes two endpoints
- GET `/api/users`: Returns a list of users. Requires JWT based authentication.
- GET `/api/token`: This is an optional endpoint, which returns a JWT access token, but not needed to run or test the service. If you wish to use this endpoint, check the details at the bottom under [Using token endpoint](#Using-token-endpoint) section.

### Run the service
There are two ways you can run the service.

Note that if you wish to change the port, you can do so by setting the ENV var `PORT` to your desired value. Alternatively, you can also provide a custom value in the `service_config/config.yml` file.

#### Running with Go
If you have `Go` installed (tested with `Go 1.23.1`), you can execute the following
- `cd service` (if you are within `user-service` directory)
- `go run main.go`

#### Running with Docker
Make sure that you have docker engine running before you execute any docker command.

**Fresh build:** You can run the service via Docker as `docker compose up --build`.

**With pre-built image** If you wish to build the image, you can execute `make docker-image` within the `service` directory. You can change the value of the image tag in the `Makefile` if you wish. This image can be used within the docker-compose file. You can comment out the `build` section and uncomment the `image` section. Make sure that the image name matches with the one that you provided while creating the docker image.

With the image, you can simply run `docker compose up`.

**Notes on docker volumes for config file and keys**
- The `service_config` directory has been committed to enable easy testing.
- Only `public.pem` file of the `keys` directory has been committed to enable service usage and testing. The public key is used to validate the already generated token for this sample service.

### Testing the service
In default mode, the service runs on port `3030`.

Since the `/api/users` endpoint requires JWT based authentication, an RSA signed token has been pre-generated with a validity of 30 days to enable testing of this sample service. The value of the token is:

```
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjbGllbnQxIiwiZXhwIjoxNzQ0NjQyNTM0LCJpYXQiOjE3NDIwNTA1MzQsImlzcyI6InBsYXRmb3JtL3VzZXItc2VydmljZSIsInN1YiI6ImNsaWVudF91c2VyIiwidXNlckNsYWltcyI6eyJ1c2VyX2lkIjoiY2xpZW50X3VzZXIifX0.OD4ftTZyXhvdb3g0TWuxhJ30ohmNmZiq5sL3Va7lGeJ9PRihH_x5hk3T6Z0Vz9jqUqsolzF1I-DYya6rUrpY6gYV-L1L94KPy_zGvyiYoYkbJrTMFlBGgsyBRwbFNeOxziAJdgiPkStelA7N_dsdR8j4AjN3VU_1t_gUZODC9G_cpFWUpK1TFqQffh5G2jdtuvvcdC8rhEn08heqtzbNKKaxbG3yA3pr6S0iSPb62Y9JAYwMTS04mynlXRGIkDSvTLlMPNTrRV0nxxQypm1f-3sy3WI-8zXWNmlnQDvCzRzun1EUy-rU3F3ImXd3NfNrYO_UC3tqYPVpyOxlVPv0rQ
```

**NOTE:** The supplied public key (`/keys/public.pem`) will be used to validate this token.

However, if you wish to generate a new token using the `/api/token` endpoint, check the details at the bottom under [Using token endpoint](#Using-token-endpoint) section.

There are a few ways that you can test the service easily:

#### Use the provided script in Makefile
Just execute `make get-users` within the `service` directory. You can check the script details in the `Makefile` and in the `get-users.sh`.

#### Use curl
This is basically same as the script file that is being provided. Use the value of the token in the following curl command.
```
curl -H "Authorization: Bearer $token" http://localhost:3030/api/users
```

#### Use an API client, such as Insomnia or Postman
The steps are same - first fetch the token, then use that token to fetch users. Ensure that you use `Bearer` (or `Bearer Token`) `Authorization` scheme.

## Implementation design and considerations

I have tried to add notes closer to the code in many cases to describe the implementation choices. In some cases, I have also hinted at better ways of implementing the features.

However, I will describe the design choice at the high level here, and will repeat some notes here as well.

### Authentication Token
- In default mode, the implementation uses `RSA` signing mechanism while generating tokens. This is advisable in a production scenario. This is also appropriate from the point of view of scalability. However, I have also added a configurable option to enable a `HMAC` signing mechanism. Some organizations use this mechanism in cases where scaleability is not a major concern and the secret can be kept encrypted.
- If you generate a token via (`/api/token`) endpoint, the validity of the token is hardcoded with 30 min for testing purposes. **However**, in production scenario, it should be less (around 10 min). The client can always refresh of get a new token re-issued.
- The token issuance endpoint (`/api/token`) is open. **However**, in production scenario, this step should be allowed only to registered clients. During registration, a client should be issued with a `client_id` and a `client_secret`. At this stage, the client should also provide the server with a `public-key` of its public/private key pair. The server will use this public-key to verify the signature of the client when issuing a token. The `client_id` and `client_secret` are also verified at this stage.
- The implementation does not deal with measures such as token re-refresh or security measures such as a deny list to protect agains the stolen tokens. **However**, such measures are desired in a production scenario.

### Authorization using Role-based Access Control (RBAC)
In the implementation, the `AccessRights` data structure specifies the schema to represent RBAC. An AccessRights policy indicates which `Role` has what kind of `Permission(s)` on what `Resource` under what `Conditions`.

I have kept the data structure relatively simple, except the `Conditions` part that offers some flexibility. Usually, a rich data structure is required depending up the complexity of the authorization policies. In my experience, such a policy schema depends heavily upon the usecase. It is also possible to keep both a rich policy schema as well as a simple RBAC schema side-by-side or in control of different services. These two types of schema work together in deciding the final authorization for a user.

#### Middleware for RBAC authorization check
The implementation primarily uses a middleware to check the RBAC. This means that the handler (`GetUsers`) can simply worry about performing the domain operation.

In a complex scenario, often there is a need to perform permission checks at the handler level as well. This happens especially when we are dealing with different categories of permissions - for example, global vs specific domain level. So, a global permission check is appropriate at the middleware level, but the specific permission checks might be performed within the handler. Such specific permission checks might happen only after the handler performs some initial operations.

### Datastore
The datastore aspect is not the focus of this sample service, so I have kept it extremely simple with some hardcoded data.

### Tests
In the interest of time, I have only written API tests (`handlers_test.go`). Ideally, tests should cover more ground at the package level.

### Others
- Hardcoded values have been used in areas that were not of importance for this sample service.
- A logger package is desired in a production scenario for structured logging.

### Using token endpoint
If you decide to use the `/api/token` endpoint, you have two options:

**OPTION 1: Your own RSA key-pair**

You will have to generate an RSA key-pair to make this work. Steps below:

1. Generate key-pair

```
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -outform PEM -pubout -out public.pem
```
Ensure that the filenames remain `private.pem` and `public.pem`. The filenames are hardcoded in the code logic.

2. Put these two files in the `keys` directory (inside `user-service` dir).

**OPTION 2: Use HMAC signed token**

Set the value of `signing-method` to `hmac` in the `service_config/config.yml` file. Alternatively, you can also set the value of ENV var `SIGNING_METHOD` to `hmac`.
