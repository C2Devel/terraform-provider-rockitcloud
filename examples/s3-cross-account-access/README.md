# S3 bucket cross-account access example

This example demonstrates how to create an S3 bucket in one account and grant access to that bucket to the user from another account using the bucket ACL.
To access infrastructure in different accounts, two providers are configured with different aliases.

The example takes credentials from two `c2rc.sh` files, one per account.
Get the file for each project and place them in this directory as `c2rc-first.sh` and `c2rc-second.sh`.

Running the example:

```shell
$ terraform init
$ terraform apply
```

Destroying the example:

```shell
$ terraform destroy
```
