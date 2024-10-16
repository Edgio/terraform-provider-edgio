## Running examples

Before you run example, you need to set environment variables, as all examples
are requiring them. You can set them in your shell or in a file (e.g. `env.sh`, `.env`).

```shell
export TF_VAR_client_id=f8c....
export TF_VAR_client_secret=veBWK...
export TF_VAR_organization_id=6b1...
```

Then you can run the example (e.g. property):

```shell
cd examples
cd property
terraform init
terraform apply
```

You can tweak values in terraform example files to see how it affects the resources.

