# Terraform Provider sda (Based on the Terraform Plugin Framework)

_This repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

This repository holds the source code for the creation of a [Terraform](https://www.terraform.io) provider for SDA SaaS, containing:

- A resource and a data source (`internal/provider/`),
- Terraform example templates fo deploying to SDA (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

Tutorials for creating Terraform providers can be found on the [HashiCorp Developer](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) platform. _Terraform Plugin Framework specific guides are titled accordingly._


Once the SDA provider is completed, we will want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider (ready for release)

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the make `release` command:

```shell
make release
```

This will:
- Build the provider for all supported OSs and Architectures
- Package (zip) the builds
- Compute the checksums for all builds and add them to a file
- Sign the checksum file

All these are required to deploy the provider to the Hashicorp repository. 

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider
```
terraform {
  required_providers {
    sda = {
      source  = "sda/sda"
      version = "0.1.0"
    }
  }
}

provider "sda" {
  host      = <sda_url>
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
