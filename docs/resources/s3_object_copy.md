---
subcategory: "S3 (Simple Storage)"
layout: "aws"
page_title: "aws_s3_object_copy"
description: |-
  Manages copying for an S3 object.
---

[canned-acl]: https://docs.k2.cloud/en/api/s3/acl.html#cannedacl
[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[RFC3339 format]: https://tools.ietf.org/html/rfc3339#section-5.8
[w3c cache_control]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9
[w3c content_disposition]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1
[w3c content_encoding]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11

# Resource: aws_s3_object_copy

Manages copying for an S3 object.

## Example Usage

```terraform
resource "aws_s3_object_copy" "example" {
  bucket = "tf-example"
  key    = "destination_key"
  source = "source_bucket/source_key"

  acl = "authenticated-read"
}
```

## Argument Reference

The following arguments are required:

* `bucket` - (Required) Name of the bucket to put the file in.
* `key` - (Required) Name of the object once it is in the bucket.
* `source` - (Required) Specifies the source object for the copy operation. The source object consists of the name of the source bucket and the key of the source object, separated by a slash (`/`). For example, `testbucket/test1.json`.

The following arguments are optional:

* `acl` - (Optional, Conflicts with `grant`) [Canned ACL][canned-acl] to apply.
    * _Valid values:_ `private`, `public-read`, `public-read-write`, `authenticated-read`
    * _Default value:_ `private`
* `content_type` - (Optional) Standard MIME type describing the format of the object data, e.g., `application/octet-stream`. All valid MIME types are valid for this input.
* `grant` - (Optional, Conflicts with `acl`) Configuration block for header grants [documented below](#grant).

### grant

This configuration block has the following required arguments:

* `permissions` - (Required) List of permissions to grant to grantee.
    * _Valid values:_ `READ`, `READ_ACP`, `WRITE_ACP`, `FULL_CONTROL`
* `type` - (Required) - Type of grantee.
    * _Valid values:_ `CanonicalUser`, `Group`, or `AmazonCustomerByEmail`

This configuration block has the following optional arguments (one of the three is required):

* `email` - (Optional) Email address of the grantee (S3 project email). Used only when `type` is `AmazonCustomerByEmail`.
* `id` - (Optional) The canonical user ID of the grantee (S3 User ID). Used only when `type` is `CanonicalUser`.
* `uri` - (Optional) URI of the grantee group. Supported groups are `http://acs.amazonaws.com/groups/global/AllUsers` and `http://acs.amazonaws.com/groups/global/AuthenticatedUsers`. Used only when `type` is `Group`.

-> **Note** Terraform ignores all leading `/`s in the object's `key` and treats multiple `/`s in the rest of the object's `key` as a single `/`, so values of `/index.html` and `index.html` correspond to the same S3 object as do `first//second///third//` and `first/second/third/`.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `etag` - ETag generated for the object (an MD5 sum of the object content). For plaintext objects the hash is an MD5 digest of the object data. For objects created by either the multipart upload or part copy operation, the hash is not an MD5 digest, regardless of the method of encryption.
* `id` - The `key` of the resource supplied above.
* `last_modified` - Returns the date that the object was last modified, in [RFC3339 format].
* `version_id` - Version ID of the newly created copy.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`cache_control`, `content_disposition`, `content_encoding`, `content_language`, `copy_if_match`, `copy_if_modified_since`, `copy_if_none_match`, `copy_if_unmodified_since`, `customer_algorithm`, `customer_key`, `customer_key_md5`, `expiration`, `expected_bucket_owner`, `expected_source_bucket_owner`, `expires`, `force_destroy`, `kms_encryption_context`, `kms_key_id`, `metadata`, `metadata_directive`, `object_lock_legal_hold_status`, `object_lock_mode`, `object_lock_retain_until_date`, `request_charged`, `request_payer`, `server_side_encryption`, `source_customer_algorithm`, `source_customer_key`, `source_customer_key_md5`, `source_version_id`, `storage_class`, `tagging_directive`, `tags`, `tags_all`, `website_redirect`.
