# EC2 transit gateway cross-account VPC attachment example

This example demonstrates how to create a transit gateway in one account and share it with a second account to attach it to a VPC.
To access infrastructure in different accounts, two providers are configured with different aliases.

The example takes credentials from two `c2rc.sh` files, one per account.
Get the file for each project and place them in this directory as `c2rc-first.sh` and `c2rc-second.sh`.
The id of the second account is derived from `c2rc-second.sh`, so no variables have to be specified.

Running the example:

```shell
$ terraform init
$ terraform apply
```

Destroying the example:

```shell
$ terraform destroy
```
