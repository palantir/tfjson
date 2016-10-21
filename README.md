tfjson
======

Utility to read in a Terraform plan file and dump it out in JSON. Standalone
version of [Terraform PR #3170](https://github.com/hashicorp/terraform/pull/3170).

## Installation

```
$ go get github.com/palantir/tfjson
```

## Usage

Given the following Terraform resources:

```hcl
provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

module "inner" {
  source = "./inner"
}

// in inner module:

resource "aws_vpc" "inner" {
  cidr_block = "10.0.0.0/8"
}
```

Running `terraform plan -out=terraform.tfplan` produces a Terraform plan file.
The JSON representation produced by `tfjson` looks like:

```json
$ tfjson terraform.tfplan
{
    "aws_vpc.main": {
        "cidr_block": "10.0.0.0/16",
        "default_network_acl_id": "",
        "default_route_table_id": "",
        "default_security_group_id": "",
        "destroy": false,
        "destroy_tainted": false,
        "dhcp_options_id": "",
        "enable_classiclink": "",
        "enable_dns_hostnames": "",
        "enable_dns_support": "",
        "id": "",
        "instance_tenancy": "",
        "main_route_table_id": ""
    },
    "destroy": false,
    "inner": {
        "aws_vpc.inner": {
            "cidr_block": "10.0.0.0/8",
            "default_network_acl_id": "",
            "default_route_table_id": "",
            "default_security_group_id": "",
            "destroy": false,
            "destroy_tainted": false,
            "dhcp_options_id": "",
            "enable_classiclink": "",
            "enable_dns_hostnames": "",
            "enable_dns_support": "",
            "id": "",
            "instance_tenancy": "",
            "main_route_table_id": ""
        },
        "destroy": false
    }
}
```

## License

This project is made available under the [MIT License](http://opensource.org/licenses/MIT).
