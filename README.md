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


------------

## To test in local 

Do not use `terraform init` to install the provider, as it will attempt to download the provider from the Terraform Registry. Use `terraform plan` or `terraform apply` instead, which will use the provider binary built locally.

### 1. MAC / Linux:

Find or create `.terraformrc` file in the root and add the following content to it:

```
provider_installation {
  dev_overrides {
      "sda/sda" = "/Users/<username>/.terraform.d/plugins/registry.terraform.io/sda/sda/0.1.0/darwin_amd64"
  }
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

### 2. Windows:
Find or create `terraform.rc` file in the `C:\Users\<username>\AppData\Roaming` and add the following content to it:

```provider_installation {
  dev_overrides {
    "sda/sda" = "C:/Users/<username>/AppData/Roaming/terraform.d/plugins/registry.terraform.io/sda/sda/0.1.0"
  }
  # For all other providers, install them from the official registry
  # or other sources as needed
}
```
**Note: The path ends at the folder, not the binary itself. You will need to copy the provider binary to the specified path.**


##############################################################################################################################################################################################################################
################################################################################ CONNECT FROM LOCAL ##########################################################################################################################
##############################################################################################################################################################################################################################

## Terraform Provider SDA — Local Setup Guide

## Prerequisites
Tool                Version     Download
Go                  1.21+       https://go.dev/dl/
Terraform           1.5+        https://developer.hashicorp.com/terraform/downloads
Git                 Latest      https://git-scm.com/downloads

## Verify installation:
bashgo version
terraform version
git --version

## Step 1: Clone the Repository
bashgit clone https://github.com/your-org/terraform-provider-sda.git
cd terraform-provider-sda

## Step 2: Install Dependencies & Build
bashgo mod download
go mod tidy
go install .

--go install compiles the provider binary and places it in your GOBIN directory automatically.

## Find Your GOBIN Path
bashgo env GOBIN
go env GOPATH

--If GOBIN is empty, the default binary location is:

OS                  Default Path
Windows             C:\Users\<username>\go\bin\terraform-provider-sda.exe
macOS               /Users/<username>/go/bin/terraform-provider-sda

## Verify Binary Exists

## Windows (PowerShell):
powershellTest-Path "$env:USERPROFILE\go\bin\terraform-provider-sda.exe"

## macOS:
bashls ~/go/bin/terraform-provider-sda

## Step 3: Configure Terraform for Local Provider
Create a Terraform CLI config file with dev_overrides so Terraform uses your local binary instead of downloading from the registry.

⚠️ Do NOT run terraform init when using dev_overrides. It will attempt to download the provider from the Terraform Registry and fail. Skip init and go directly to terraform plan or terraform apply.


## Windows

1. Find your APPDATA path:
powershell$env:APPDATA

This typically returns C:\Users\<username>\AppData\Roaming.

2. Create terraform.rc in that directory:
File path: C:\Users\<username>\AppData\Roaming\terraform.rc

hclprovider_installation {
  dev_overrides {
    "sda/sda" = "C:/Users/<username>/go/bin"
  }
  direct {}
}

--Replace <username> with your actual Windows username. Use forward slashes (/) in the path.

## Or create it via PowerShell (one command):
@"
provider_installation {
  dev_overrides {
    "sda/sda" = "C:/Users/$env:USERNAME/go/bin"
  }
  direct {}
}
"@ | Out-File -FilePath "$env:APPDATA\terraform.rc" -Encoding utf8

3. Verify:
powershellcat "$env:APPDATA\terraform.rc"

## macOS / Linux

1. Create or edit ~/.terraformrc:

cat > ~/.terraformrc << 'EOF'
provider_installation {
  dev_overrides {
    "sda/sda" = "/Users/<username>/go/bin"
  }
  direct {}
}
EOF

Replace <username> with your actual macOS username, or use the output of go env GOPATH + /bin.

2. Verify:

cat ~/.terraformrc

Step 4: Create a Test Configuration
Create a working directory for testing:

## Windows:
mkdir C:\Users\<username>\projects\test-sda
cd C:\Users\<username>\projects\test-sda

## macOS:
mkdir -p ~/projects/test-sda
cd ~/projects/test-sda

## Create main.tf
hclterraform {
  required_providers {
    sda = {
      source = "sda/sda"
    }
  }
}

provider "sda" {
  host     = "https://api.sdaconsole.io"
  username = var.sda_username
  password = var.sda_password
}


## Create variables.tf
hclvariable "sda_username" {
  type        = string
  description = "SDA account username"
}

variable "sda_password" {
  type        = string
  sensitive   = true
  description = "SDA account password"
}


## Create terraform.tfvars
hclsda_username = "your-username"
sda_password = "your-password"

--🔒 Security: Never commit terraform.tfvars to Git. Add it to .gitignore.

## Alternative: Use Environment Variables Instead of tfvars

## Windows (PowerShell):
-$env:SDA_HOST = "https://api.sdaconsole.io"
-$env:SDA_USERNAME = "your-username"
-$env:SDA_PASSWORD = "your-password"

## macOS:
export SDA_HOST="https://api.sdaconsole.io"
export SDA_USERNAME="your-username"
export SDA_PASSWORD="your-password"

When using environment variables, the provider block can be empty:
hclprovider "sda" {}

Step 5: Run Terraform
From your test directory:

terraform plan
terraform apply
terraform destroy

## Remember: No terraform init. Go straight to terraform plan.

Expected output on first run:
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - sda/sda in C:\Users\<username>\go\bin
This warning is normal and confirms Terraform is using your local provider.

## Development Workflow

Every time you make code changes to the provider:

# 1. Go to provider project
cd /path/to/terraform-provider-sda

# 2. Rebuild and install
go install .

# 3. Go to test directory
cd /path/to/test-sda

# 4. Test
terraform plan
terraform apply

Run Tests
# All tests
go test ./internal/...

# Specific resource tests
go test ./internal/provider/device/... -v

# Acceptance tests (calls real API)
make testacc

Troubleshooting
Problem                                                                         Solution
go: command not found                                                           Add C:\Program Files\Go\bin (Windows) or /usr/local/go/bin (Mac) to your system PATH. Restart terminal.
terraform: command not foundAdd                                                 the  folder containing the terraform binary to your system PATH.Error: Failed to query available provider packagesYou ran terraform init. Don't do that — skip directly to terraform plan.Unable to Create SDA API ClientCheck that your host URL, username, and password are correct. Verify you can reach the API host (curl https://api.sdaconsole.io). You may need VPN.invalid character '<' looking for beginning of valueThe host URL is wrong — it is returning an HTML page instead of JSON. Verify the correct API base URL with your team.Build fails with import errorsRun go mod tidy to fix dependencies.Binary not found after go installCheck go env GOPATH and look in the bin subdirectory.Changes not reflected after rebuildMake sure you ran go install . again and are pointing to the correct GOBIN path in your terraform.rc / .terraformrc.