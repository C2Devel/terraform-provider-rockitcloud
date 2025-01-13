---
layout: "aws"
page_title: "Provider: Rockit Cloud"
description: |-
  Use the Terraform Rockit Cloud Provider to manage resources in public and private clouds based on K2 Cloud and Rockit Cloud platforms. You must configure the provider with the proper credentials before you can use it.
---

[hashicorp-tutorials]: https://learn.hashicorp.com/tutorials/terraform/infrastructure-as-code?in=terraform/aws-get-started&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS
[aws-tutorials]: https://learn.hashicorp.com/collections/terraform/aws?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS
[c2-tutorials]: https://github.com/C2Devel/terraform-examples
[aws-configure-files]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
[terraform]: https://docs.cloud.croc.ru/en/api/tools/terraform.html

# Rockit Cloud Provider

The Rockit Cloud Provider is based on the AWS Provider.
It is designed for managing resources in public and private clouds based on K2 Cloud and Rockit Cloud platforms.
The provider must be configured with the proper credentials before you can use it.

Use the navigation on the left to read about the available resources.

~> **Note**
Resource names in the navigation bar have an automatically generated prefix that matches the *rockitcloud* name.
For compatibility with AWS provider configurations, we retained the ``aws`` prefix in resource description and usage examples.

To learn the basics of Terraform using this provider, follow the
hands-on [get started tutorials][hashicorp-tutorials] on HashiCorp's Learn platform.

Examples of using Rockit Cloud services with Terraform can be found in [reference test suite on GitHub][c2-tutorials].

Rockit Cloud API is based on AWS API so you can also see examples of using AWS services with Terraform: [AWS services tutorials][aws-tutorials].

## Example Usage

For Terraform 0.13 and later:

```terraform
terraform {
  required_providers {
    aws = {
      source  = "c2devel/rockitcloud"
      version = "24.1.0"
    }
  }
}

# Configure the rockitcloud provider.
# The section is named `aws` for backward compatibility.
provider "aws" {
  region = "region-1"
}

# Create a VPC.
resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}
```

## Authentication and Configuration

Configuration for the Rockit Cloud Provider can be derived from several sources,
which are applied in the following order:

1. Parameters in the provider configuration.
2. Environment variables.
3. Shared configuration and credentials files.

### Provider Configuration

!> **Warning:** Hard-coded credentials are not recommended in any Terraform
configuration because they run the risk of secret leakage should this file ever be committed to a
public version control system.

Credentials can be provided by adding `access_key` and `secret_key` to the `aws` provider block.

Usage:

```terraform
provider "aws" {
  region     = "region-1"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment Variables

Credentials can also be provided by using the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
The region can be set using the `AWS_REGION` or `AWS_DEFAULT_REGION` environment variables.

For example:

```terraform
provider "aws" {}
```

```sh
$ export AWS_ACCESS_KEY_ID="my-access-key"
$ export AWS_SECRET_ACCESS_KEY="my-secret-key"
$ export AWS_REGION="region-1"
$ terraform plan
```

### Shared Configuration and Credentials Files

Rockit Cloud Provider can use [AWS shared configuration and credentials files][aws-configure-files] and source credentials and other settings from them.
By default, these files are located at `$HOME/.aws/config` and `$HOME/.aws/credentials` on Linux and macOS,
and `"%USERPROFILE%\.aws\config"` and `"%USERPROFILE%\.aws\credentials"` on Windows.

If no named profile is specified, the `default` profile is used.
Use the `profile` parameter or `AWS_PROFILE`, `AWS_DEFAULT_PROFILE` environment variables to specify a named profile.

The locations of the shared configuration and credentials files can be configured using either
the parameters `shared_config_files` and `shared_credentials_files`
or the environment variables `AWS_CONFIG_FILE` and `AWS_SHARED_CREDENTIALS_FILE`.

For example:

```terraform
provider "aws" {
  shared_config_files      = ["/Users/tf_user/.aws/conf"]
  shared_credentials_files = ["/Users/tf_user/.aws/creds"]
  profile                  = "customprofile"
}
```

## Rockit Cloud Provider Full Configuration

For more information about the Rockit Cloud Provider configuration, see the documentation on [using Terraform][terraform].
