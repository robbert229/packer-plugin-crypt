# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "crypt-mkpasswd" "default" {
  plaintext = "foobar"
  algorithm = "sha512"
}

locals {
  hash = data.crypt-mkpasswd.default.result
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo foo: ${local.hash}",
    ]
  }
}
