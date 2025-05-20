// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package mkpasswd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/robbert229/packer-plugin-crypt/internal/crypt"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Plaintext string `mapstructure:"plaintext" required:"true"`
	Algorithm string `mapstructure:"algorithm" required:"false"`
	Salt      string `mapstructure:"salt" required:"false"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Result string `mapstructure:"result"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	if d.config.Plaintext == "" {
		return fmt.Errorf("you must provide a plaintext value")
	}

	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	salt := d.config.Salt
	if salt == "" {
		saltBytes := make([]byte, 8)
		_, err := rand.Read(saltBytes)
		if err != nil {
			return cty.NilVal, err
		}

		salt = base64.StdEncoding.EncodeToString(saltBytes)
	}

	hash, err := crypt.Hash(d.config.Plaintext, crypt.HashOptions{
		Algorithm: d.config.Algorithm,
		Salt:      d.config.Salt,
	})
	if err != nil {
		return cty.NilVal, err
	}

	output := DatasourceOutput{
		Result: hash,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
