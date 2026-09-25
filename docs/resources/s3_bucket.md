---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_bucket"
description: |-
  Manages an S3 bucket.
---

[bucket-naming]: https://docs.k2.cloud/en/services/object_storage/operations.html#s3bucketnaming
[canned-acl]: https://docs.k2.cloud/en/api/s3/acl.html#cannedacl
[cors]: https://docs.k2.cloud/en/services/object_storage/operations.html#cors
[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[hosting-website]: https://docs.k2.cloud/en/services/object_storage/operations.html#objectstoragestaticwebsitesmanual
[lifecycle-management]: https://docs.k2.cloud/en/services/object_storage/operations.html#id28
[s3-versioning]: https://docs.k2.cloud/en/services/object_storage/operations.html#s3versioningmanual
[website-redirect-rules]: https://docs.k2.cloud/en/services/object_storage/instructions.html#s3setredirectiontowebsite

# Resource: aws_s3_bucket

Manages an S3 bucket.

~> **Note on S3 Bucket canned ACL Configuration:** S3 bucket canned ACL can be configured in either the standalone resource [`aws_s3_bucket_acl`](s3_bucket_acl.md)
or with the deprecated parameter `acl` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

~> **Note on S3 Bucket ACL Grants Configuration:** S3 bucket grants can be configured in either the standalone resource [`aws_s3_bucket_acl`](s3_bucket_acl.md)
or with the deprecated parameter `grant` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

~> **Note on S3 Bucket CORS Configuration:** S3 bucket CORS can be configured in either the standalone resource [`aws_s3_bucket_cors_configuration`](s3_bucket_cors_configuration.md)
or with the deprecated parameter `cors_rule` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

~> **Note on S3 Bucket Lifecycle Configuration:** S3 bucket lifecycle can be configured in either the standalone resource [`aws_s3_bucket_lifecycle_configuration`](s3_bucket_lifecycle_configuration.md)
or with the deprecated parameter `lifecycle_rule` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

~> **Note on S3 Bucket Policy Configuration:** S3 bucket policy can be configured in either the standalone resource [`aws_s3_bucket_policy`](s3_bucket_policy.md)
or with the deprecated parameter `policy` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

~> **Note on S3 Bucket Versioning Configuration:** S3 bucket versioning can be configured in either the standalone resource [`aws_s3_bucket_versioning`](s3_bucket_versioning.md)
or with the deprecated parameter `versioning` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

~> **Note on S3 Bucket Website Configuration:** S3 bucket website can be configured in either the standalone resource [`aws_s3_bucket_website_configuration`](s3_bucket_website_configuration.md)
or with the deprecated parameter `website` in the resource `aws_s3_bucket`.
Configuring with both will cause inconsistencies and may overwrite configuration.

## Example Usage

### Private Bucket With Tags

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  tags = {
    Name        = "tf-example"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}
```

### Static Website Hosting

-> **Note** The parameter `website` is deprecated.
Use the resource [`aws_s3_bucket_website_configuration`](s3_bucket_website_configuration.md) instead.

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"

    routing_rules = <<EOF
[{
    "Condition": {
        "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
        "ReplaceKeyPrefixWith": "documents/"
    }
}]
EOF
  }
}
```

### Using CORS

-> **Note** The parameter `cors_rule` is deprecated.
Use the resource [`aws_s3_bucket_cors_configuration`](s3_bucket_cors_configuration.md) instead.

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"
  acl    = "public-read"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
```

### Using versioning

-> **Note** The parameter `versioning` is deprecated.
Use the resource [`aws_s3_bucket_versioning`](s3_bucket_versioning.md) instead.

```terraform
resource "aws_s3_bucket" "example" {
  bucket = "tf-example"
  acl    = "private"

  versioning {
    enabled = true
  }
}
```

### Using object lifecycle

-> **Note** The parameter `lifecycle_rule` is deprecated.
Use the resource [`aws_s3_bucket_lifecycle_configuration`](s3_bucket_lifecycle_configuration.md) instead.

```terraform
resource "aws_s3_bucket" "bucket" {
  bucket = "tf-example"
  acl    = "private"

  lifecycle_rule {
    id      = "log"
    enabled = true

    prefix = "log/"

    tags = {
      rule      = "log"
      autoclean = "true"
    }

    expiration {
      days = 90
    }
  }

  lifecycle_rule {
    id      = "tmp"
    prefix  = "tmp/"
    enabled = true

    expiration {
      date = "2016-01-12"
    }
  }
}

resource "aws_s3_bucket" "versioning_bucket" {
  bucket = "tf-example"
  acl    = "private"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    prefix  = "config/"
    enabled = true

    noncurrent_version_expiration {
      days = 90
    }
  }
}
```

### Using ACL policy grants

-> **Note** The parameters `acl` and `grant` are deprecated.
Use the resource [`aws_s3_bucket_acl`](s3_bucket_acl.md) instead.

```terraform
data "aws_canonical_user_id" "current_user" {}

resource "aws_s3_bucket" "example" {
  bucket = "tf-example"

  grant {
    id          = data.aws_canonical_user_id.current_user.id
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
  }

  grant {
    type        = "Group"
    permissions = ["READ_ACP", "WRITE"]
    uri         = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional, **Deprecated**, Conflicts with `grant`) The [canned ACL][canned-acl] to apply.  Terraform will only perform drift detection if a configuration value is provided. Use the resource [`aws_s3_bucket_acl`](s3_bucket_acl.md) instead.
    * _Valid values:_  `authenticated-read`, `private`, `public-read`, `public-read-write`
    * _Default value:_ `private`
* `bucket` - (Optional, Forces new resource) The name of the bucket. If omitted, Terraform will assign a random, unique name. Must be lowercase and less than or equal to 63 characters in length. A full list of bucket naming rules may be found in [user documentation][bucket-naming].
* `bucket_prefix` - (Optional, Conflicts with `bucket`, Forces new resource) Creates a unique bucket name beginning with the specified prefix. Must be lowercase and less than or equal to 37 characters in length. A full list of bucket naming rules may be found in [user documentation][bucket-naming].
* `cors_rule` - (Optional, **Deprecated**) A rule of [Cross-Origin Resource Sharing][cors]. See [CORS rule](#cors-rule) below for details. Terraform will only perform drift detection if a configuration value is provided. Use the resource [`aws_s3_bucket_cors_configuration`](s3_bucket_cors_configuration.md) instead.
* `grant` - (Optional, **Deprecated**, Conflicts with `acl`) An ACL policy grant. See [Grant](#grant) below for details. Terraform will only perform drift detection if a configuration value is provided. Use the resource [`aws_s3_bucket_acl`](s3_bucket_acl.md) instead.
* `force_destroy` - (Optional) A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error.
    * _Default value:_ `false`
* `lifecycle_rule` - (Optional, **Deprecated**) A configuration of [object lifecycle management][lifecycle-management]. See [Lifecycle Rule](#lifecycle-rule) below for details. Terraform will only perform drift detection if a configuration value is provided.
  Use the resource [`aws_s3_bucket_lifecycle_configuration`](s3_bucket_lifecycle_configuration.md) instead.
* `policy` - (Optional, **Deprecated**) A valid bucket policy JSON document. Note that if the policy document is not specific enough (but still valid), Terraform may view the policy as constantly changing in a `terraform plan`. In this case, please make sure you use the verbose/specific version of the policy.
  Terraform will only perform drift detection if a configuration value is provided.
  Use the resource [`aws_s3_bucket_policy`](s3_bucket_policy.md) instead.
* `versioning` - (Optional, **Deprecated**) A configuration of the [S3 bucket versioning state][s3-versioning]. See [Versioning](#versioning) below for details. Terraform will only perform drift detection if a configuration value is provided. Use the resource [`aws_s3_bucket_versioning`](s3_bucket_versioning.md) instead.
* `website` - (Optional, **Deprecated**) A configuration of the [S3 bucket website][hosting-website]. See [Website](#website) below for details. Terraform will only perform drift detection if a configuration value is provided.
  Use the resource [`aws_s3_bucket_website_configuration`](s3_bucket_website_configuration.md) instead.
* `tags` - (Optional) Map of tags to assign to the bucket. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

### CORS Rule

~> **Note** Currently, changes to the `cors_rule` configuration of _existing_ resources cannot be automatically detected by Terraform. To manage changes of CORS rules to an S3 bucket, use the `aws_s3_bucket_cors_configuration` resource instead. If you use `cors_rule` on an `aws_s3_bucket`, Terraform will assume management over the full set of CORS rules for the S3 bucket, treating additional CORS rules as drift. For this reason, `cors_rule` cannot be mixed with the external `aws_s3_bucket_cors_configuration` resource for a given S3 bucket.

The `cors_rule` configuration block supports the following arguments:

* `allowed_headers` - (Optional) List of headers allowed.
* `allowed_methods` - (Required) One or more HTTP methods that you allow the origin to execute. Can be `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.
* `allowed_origins` - (Required) One or more origins you want customers to be able to access the bucket from.
* `expose_headers` - (Optional) One or more headers in the response that you want customers to be able to access from their applications (for example, from a JavaScript `XMLHttpRequest` object).
* `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

### Grant

~> **Note** Currently, changes to the `grant` configuration of _existing_ resources cannot be automatically detected by Terraform. To manage changes of ACL grants to an S3 bucket, use the `aws_s3_bucket_acl` resource instead. If you use `grant` on an `aws_s3_bucket`, Terraform will assume management over the full set of ACL grants for the S3 bucket, treating additional ACL grants as drift. For this reason, `grant` cannot be mixed with the external `aws_s3_bucket_acl` resource for a given S3 bucket.

The `grant` configuration block supports the following arguments:

* `id` - (Optional) Canonical user ID to grant for (S3 User ID). Used only when `type` is `CanonicalUser`.
* `type` - (Required) Type of grantee to apply for.  `AmazonCustomerByEmail` is not supported.
    * _Valid values:_ `CanonicalUser` and `Group`
* `permissions` - (Required) List of permissions to apply for grantee.
    * _Valid values:_  `FULL_CONTROL`, `READ_ACP`, `READ`, `WRITE`, `WRITE_ACP`
* `uri` - (Optional) URI address to grant for. Supported groups are `http://acs.amazonaws.com/groups/global/AllUsers` and `http://acs.amazonaws.com/groups/global/AuthenticatedUsers`. Used only when `type` is `Group`.

### Lifecycle Rule

~> **Note** Currently, changes to the `lifecycle_rule` configuration of _existing_ resources cannot be automatically detected by Terraform. To manage changes of lLifecycle rules to an S3 bucket, use the `aws_s3_bucket_lifecycle_configuration` resource instead. If you use `lifecycle_rule` on an `aws_s3_bucket`, Terraform will assume management over the full set of lifecycle rules for the S3 bucket, treating additional lifecycle rules as drift. For this reason, `lifecycle_rule` cannot be mixed with the external `aws_s3_bucket_lifecycle_configuration` resource for a given S3 bucket.

~> **Note** At least one of `abort_incomplete_multipart_upload_days`, `expiration`, `noncurrent_version_expiration`, must be specified.

The `lifecycle_rule` configuration block supports the following arguments:

* `id` - (Optional) Unique identifier for the rule. Must be less than or equal to 255 characters in length.
* `tags` - (Optional) Specifies object tags key and value.
* `enabled` - (Required) Specifies lifecycle rule status.
* `expiration` - (Optional) Specifies a period in the object's expire. See [Expiration](#expiration) below for details.
* `noncurrent_version_expiration` - (Optional) Specifies when noncurrent object versions expire. See [Noncurrent Version Expiration](#noncurrent-version-expiration) below for details.

#### Expiration

The `expiration` configuration block supports the following arguments:

* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect.
* `days` - (Optional) Specifies the number of days after object creation when the specific rule action takes effect.
* `expired_object_delete_marker` - (Optional) On a versioned bucket (versioning-enabled or versioning-suspended bucket), you can add this element in the lifecycle configuration to direct S3 to delete expired object delete markers. This cannot be specified with days or date in a lifecycle expiration policy.

#### Noncurrent Version Expiration

The `noncurrent_version_expiration` configuration block supports the following arguments:

* `days` - (Required) Specifies the number of days noncurrent object versions expire.

### Versioning

~> **Note** Currently, changes to the `versioning` configuration of _existing_ resources cannot be automatically detected by Terraform. To manage changes of versioning state to an S3 bucket, use the `aws_s3_bucket_versioning` resource instead. If you use `versioning` on an `aws_s3_bucket`, Terraform will assume management over the versioning state of the S3 bucket, treating additional versioning state changes as drift. For this reason, `versioning` cannot be mixed with the external `aws_s3_bucket_versioning` resource for a given S3 bucket.

The `versioning` configuration block supports the following arguments:

* `enabled` - (Optional) Enable versioning. Once you version-enable a bucket, it can never return to an unversioned state. You can, however, suspend versioning on that bucket.

### Website

~> **Note** Currently, changes to the `website` configuration of _existing_ resources cannot be automatically detected by Terraform. To manage changes to the website configuration of an S3 bucket, use the `aws_s3_bucket_website_configuration` resource instead. If you use `website` on an `aws_s3_bucket`, Terraform will assume management over the configuration of the website of the S3 bucket, treating additional website configuration changes as drift. For this reason, `website` cannot be mixed with the external `aws_s3_bucket_website_configuration` resource for a given S3 bucket.

The `website` configuration block supports the following arguments:

* `index_document` - (Required, unless using `redirect_all_requests_to`) S3 returns this index document when requests are made to the root domain or any of the subfolders.
* `error_document` - (Optional) An absolute path to the document to return in case of a 4XX error.
* `redirect_all_requests_to` - (Optional) A hostname to redirect all website requests for this bucket to. Hostname can optionally be prefixed with a protocol (`http://` or `https://`) to use when redirecting requests. The default is the protocol that is used in the original request.
* `routing_rules` - (Optional) A json array containing [routing rules][website-redirect-rules] describing redirect behavior and when redirects are applied.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.
* `arn` - The Amazon Resource Name (ARN) of the bucket. Will be of format `arn:aws:s3:::bucketname`.
* `region` - The region this bucket resides in.
* `tags_all` - Map of tags assigned to the bucket, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`acceleration_status`, `bucket_domain_name`, `bucket_regional_domain_name`, `hosted_zone_id`, `lifecycle_rule.abort_incomplete_multipart_upload_days`, `lifecycle_rule.prefix`, `lifecycle_rule.noncurrent_version_transition`, `lifecycle_rule.transition`, `logging`, `object_lock_configuration`, `object_lock_enabled`, `replication_configuration`, `request_payer`, `server_side_encryption_configuration`, `versioning.mfa_delete`, `website_domain`, `website_endpoint`.

## Import

S3 bucket can be imported using the `bucket`, e.g.,

```
$ terraform import aws_s3_bucket.bucket bucket-name
```
