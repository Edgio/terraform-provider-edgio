# Edgio Terraform Provider

This is a Terraform provider for Edgio API. This provider is based on Terraform Plugin Framework, for more information you can check the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Examples

There are examples in the `examples` directory. To run an example, navigate to the example directory and run `terraform init` and `terraform apply`.
For more information check Readme.md in the examples directory.

## Debugging

To enable terraform debugging you need Dlv:

```shell
go get github.com/go-delve/delve/cmd/dlv
```

Before runninf Dlv, you need to build the provider with the following command:

```shell
go build -gcflags "all=-N -l" -o terraform-provider-edgio
```

Then you can run the following command to debug the provider:

```shell
dlv exec --headless --listen=:2345 --api-version=2 ./terraform-provider-edgio
```

Once Dlv is running, attach to the process with your IDE, once Dlv detects the connection, it will output the following message:

```shell
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

        TF_REATTACH_PROVIDERS='{"hashicorp.com/edu/edgio":{"Protocol":"grpc","ProtocolVersion":6,"Pid":62355,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/hc/mcfd08xn3k55gxznww512l_80000gn/T/plugin4155344464"}}}'
```

Copy the `TF_REATTACH_PROVIDERS` environment variable and set it in your terminal, then you can run the terraform command you want to debug.

