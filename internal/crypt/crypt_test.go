package crypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashing(t *testing.T) {
	hash, err := Hash("password", HashOptions{
		Algorithm: "sha512",
		Salt:      "salt",
	})
	require.NoError(t, err)
	require.Equal(t, "$6$salt$IxDD3jeSOb5eB1CX5LBsqZFVkJdido3OUILO5Ifz5iwMuTS4XMS130MTSuDDl3aCI6WouIL9AjRbLCelDCy.g.", hash)

	r, err := Hash("password", HashOptions{
		Algorithm: "sha512",
		Salt:      "",
	})
	require.NoError(t, err)
	require.NotEqual(t, hash, r)
}

func wantEq(value string) func(t *testing.T, actual string) {
	return func(t *testing.T, actual string) {
		require.Equal(t, value, actual)
	}
}

func wantNe(value string) func(t *testing.T, actual string) {
	return func(t *testing.T, actual string) {
		require.NotEqual(t, value, actual)
	}
}

func TestHash(t *testing.T) {
	type args struct {
		plaintext string
		opt       HashOptions
	}
	tests := []struct {
		name    string
		args    args
		wantFn  func(t *testing.T, actual string)
		wantErr bool
	}{
		{
			name: "sha512 with salt",
			args: args{
				plaintext: "password",
				opt: HashOptions{
					Algorithm: "sha512",
					Salt:      "salt",
				},
			},
			wantFn:  wantEq("$6$salt$IxDD3jeSOb5eB1CX5LBsqZFVkJdido3OUILO5Ifz5iwMuTS4XMS130MTSuDDl3aCI6WouIL9AjRbLCelDCy.g."),
			wantErr: false,
		},
		{
			name: "sha512 with generated salt",
			args: args{
				plaintext: "password",
				opt: HashOptions{
					Algorithm: "sha512",
					Salt:      "",
				},
			},
			wantFn:  wantNe("$6$salt$IxDD3jeSOb5eB1CX5LBsqZFVkJdido3OUILO5Ifz5iwMuTS4XMS130MTSuDDl3aCI6WouIL9AjRbLCelDCy.g."),
			wantErr: false,
		},
		{
			name: "unsupported algorithm",
			args: args{
				plaintext: "password",
				opt: HashOptions{
					Algorithm: "unsupported",
					Salt:      "salt",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.plaintext, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantFn != nil {
				tt.wantFn(t, got)
			}
		})
	}
}

func TestHashNoCollission(t *testing.T) {
	// Test that the hash function does not produce collisions
	hash1, err := Hash("password1", HashOptions{
		Algorithm: "sha512",
		Salt:      "",
	})
	require.NoError(t, err)

	hash2, err := Hash("password2", HashOptions{
		Algorithm: "sha512",
		Salt:      "",
	})
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2)
}
