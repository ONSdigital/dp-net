# awsauth package

## Round Tripper

The `NewAWSSignerRoundTripper` creates a http.Transport that can then be used instead of the default transport used in
this libraries http package to autosign http requests to AWS services, such as AWS Elasticsearch.

See [Signer](#signer) for details on various signing mechanisms.

### v2

Below is an example of how to setup the aws signer round tripper and attach to a new client using
the `NewClientWithTransport` method in http package. Example uses elasticsearch as the aws service denoted by `es`
and `eu-west-1` for region.

    ```go
    import (
        "github.com/ONSdigital/dp-net/v2/awsauth"
        dphttp "github.com/ONSdigital/dp-net/v2/http"
    )
    ...
    awsSignerRT, err := awsauth.NewAWSSignerRoundTripper("", "", "eu-west-1", "es")
	if err != nil {
		...
	}

	httpClient := dphttp.NewClientWithTransport(awsSignerRT)
    ```

If you are looking to connect a local instance of application to managed AWS Elasticsearch,
then you will need to create a tunnel onto the VPC that Elasticsearch is running on, and
implement the following in your application:

    ```go
    import (
        "github.com/ONSdigital/dp-net/v2/awsauth"
        dphttp "github.com/ONSdigital/dp-net/v2/http"
    )
    ...

    awsSignerRT, err := awsauth.NewAWSSignerRoundTripper("~/.aws/credentials", "default", "eu-west-1", "es",
        awsauth.Options{TlsInsecureSkipVerify: true})
	if err != nil {
		...
	}

	httpClient := dphttp.NewClientWithTransport(awsSignerRT)
    ```

The file location and profile need to be set as the first two variables in method signature
respectively; in the above example these are the default values expected across the industry.

:warning: setting `TlsInsecureSkipVerify` to `true` should only be used for developer testing. If used in an application
use a new environment variable to control whether this is on/off,
e.g. `awsauth.Options{TlsInsecureSkipVerify: cfg.TlsInsecureSkipVerify}` :warning:

### Updates from v2 to v3:

In v3, the NewAWSSignerRoundTripper function now accepts context.Context as an argument in the request to facilitate
better control over request-scoped values, like deadlines or cancellation signals. You can pass the context.Context to
this function to provide more flexibility when making requests in applications with varying lifetimes or environments.

```
import (
    "github.com/ONSdigital/dp-net/v3/awsauth"
    dphttp "github.com/ONSdigital/dp-net/v3/http"
    "context"
    "time"
)

...

// Create context for the request
ctx := context.Background()

// Setup AWS Signer Round Tripper with context
awsSignerRT, err := awsauth.NewAWSSignerRoundTripper(ctx, "", "", "eu-west-1", "es")
if err != nil {
    ...
}

httpClient := dphttp.NewClientWithTransport(awsSignerRT)

```

## Signer

Using AWS SDK library to create a signer function to successfully sign elasticsearch requests hosted in AWS.
The function adds multiple providers to the credentials chain that is used by
the [AWS SDK V4 signer method `Sign`](https://docs.aws.amazon.com/sdk-for-go/api/aws/signer/v4/#Signer.Sign).

### Updates from v2 to v3:

In v3, the context has been added as input to both the roundtripper and signer functions. When creating a signer, you now pass a context.Context to the NewAwsSigner function to provide better request management, especially in environments where request-specific context (e.g., cancellation, deadlines) is required.

1) **Environment Provider** will attempt to retrieve credentials from `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
   set on the environment.

   Requires `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` variables to be set/exported onto the environment.

    ```go
    import awsauth "github.com/ONSdigital/dp-net/v3/awsauth"
    ...
        signer, err := esauth.NewAwsSigner(ctx, "", "", "eu-west-1", "es")
        if err != nil {
            ... // Handle error
        }
    ...
    ```

2) **Shared Credentials Provider** will attempt to retrieve credentials from the absolute path to credentials file, the
   default value if the filename is set to empty string `""` will be `~/.aws/credentials` and the default profile will
   be `default` if set to empty string.

   Requires credentials file to exist in the location specified in NewAwsSigner func.
   File must contain the keys necessary under the matching Profile heading, see example below:

    ```
        [development]
        aws_access_key_id=<access key id>
        aws_secret_access_key=<secret access key>
        region=<region>
    ```

    ```go
    import esauth "github.com/ONSdigital/dp-net/v3/awsauth"
    ...
        signer, err := esauth.NewAwsSigner(ctx, "~/.aws/credentials", "development", "eu-west-1", "es")
        if err != nil {
            ...
        }
    ...
    ```

3) **EC2 Role Provider** will attempt to retrieve credentials using an EC2 metadata client (this is created using an AWS
   SDK session).

   Requires Code is run on EC2 instance.

    ```go
    import esauth "github.com/ONSdigital/dp-net/v3/awsauth"
    ...
        signer, err := esauth.NewAwsSigner(ctx, "", "", "eu-west-1", "es")
        if err != nil {
            ...
        }
    ...
    ```

For more information on Providers for obtaining
credentials, [see AWS documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials).

The signer object should be created once on startup of an application and reused for each request, otherwise you will
experience performance issues due to creating a session for every request.

To sign elasticsearch requests, one can use the signer like so:

```go
    ...

    var req *http.Request
   
    // TODO set request
   
    var bodyReader io.ReadSeeker
   
    if payload != <zero value of type> { // Check for a payload
        bodyReader = bytes.NewReader(<payload in []byte>)
        req, err = http.NewRequest(<method>, <path>, bodyReader)
    } else { // No payload (request body is empty)
        req, err = http.NewRequest(<method>, <path>, nil)
    }
   
    if err = signer.Sign(req, bodyReader, time.Now()); err != nil {
        ... // handle error
    }
   
   ...
```