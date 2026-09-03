# Rockit Cloud provider examples

This directory contains a set of examples of using various Rockit Cloud services with Terraform.
The examples have their own README containing more details on what the example does.

The examples take credentials from a `c2rc.sh` file.
Get the file for your project and place it in the example's directory before running the example.
The cross-account examples need one file per account, see their README for details.

To run any example, clone the repository and run `terraform apply` within the example's own directory.

For example:

```shell
$ git clone https://github.com/C2Devel/terraform-provider-rockitcloud
$ cd terraform-provider-rockitcloud/examples/count
$ cp /path/to/c2rc.sh .
$ terraform apply
...
```
