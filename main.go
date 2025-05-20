// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/robbert229/packer-plugin-crypt/datasource/mkpasswd"
	"github.com/robbert229/packer-plugin-crypt/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("mkpasswd", new(mkpasswd.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
