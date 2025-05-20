// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mkpasswd

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/stretchr/testify/require"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

// Run with: PACKER_ACC=1 go test -count 1 -v ./datasource/scaffolding/data_acc_test.go  -timeout=120m
func TestAccScaffoldingDatasource(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "mkpasswd_datasource_basic_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceHCL2Basic,
		Type:     "crypt-mkpasswd",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := io.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			require.Contains(t, logsString, "hash: $6$salt$IxDD3jeSOb5eB1CX5LBsqZFVkJdido3OUILO5Ifz5iwMuTS4XMS130MTSuDDl3aCI6WouIL9AjRbLCelDCy.g.")
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccMissingPlaintextDatasource(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "mkpasswd_datasource_missingPlainTextShouldError_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: `
data "crypt-mkpasswd" "default" {
		# nothing
}
		`,
		Type: "crypt-mkpasswd",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			require.NotEqual(t, 0, buildCommand.ProcessState.ExitCode(), "")
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccMissingAlgorithmDatasource(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "mkpasswd_datasource_missingAlgorithmShouldFallbackToSha512_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: `
data "crypt-mkpasswd" "default" {
		# nothing
}
		`,
		Type: "crypt-mkpasswd",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			require.NotEqual(t, 0, buildCommand.ProcessState.ExitCode(), "")
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
