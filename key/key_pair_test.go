package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyFromHex(t *testing.T) {
	tests := []struct {
		name        string
		hexKey      string
		wantAddress string
		wantErr     bool
	}{
		{
			"test 1",
			"00eb4b724feda53fbd1cab17f7009808e22eba716504beb651df7d782521c289",
			"4c6194eb29424c17a3b6c669c36fb9c31de4e9b3",
			false,
		},
		{
			"test 2",
			"88f8f9e8f83f25eb1b5d1b554ff98ffb3741595b84e051f9cf357c400af9f811",
			"0361515c23a264d13c897a2ae74ab46211781b0e",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKeyFromHex(tt.hexKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.hexKey, got.PrivateKey())
			assert.Equal(t, tt.wantAddress, got.Address())
		})
	}
}

func Test_GetPublicKeyByPrivHex(t *testing.T) {

	tests := []struct {
		name       string
		hexKey     string
		compressed bool
		want       string
		wantErr    bool
	}{
		{
			"test 1",
			"88f8f9e8f83f25eb1b5d1b554ff98ffb3741595b84e051f9cf357c400af9f811",
			true,
			"0382e51eae6be5b247d7467727d0c6bf7bbe58e84cb6f913a8db28082c1c4f093d",
			false,
		},
		{
			"test 1",
			"88f8f9e8f83f25eb1b5d1b554ff98ffb3741595b84e051f9cf357c400af9f811",
			false,
			"0482e51eae6be5b247d7467727d0c6bf7bbe58e84cb6f913a8db28082c1c4f093df034be14690fada529a443c1bf648dcd18f621b8fa23acf7fec843d891d6978b",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKeyFromHex(tt.hexKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			pk := got.PublicKey(tt.compressed)
			assert.Equal(t, tt.want, pk)
		})
	}
}

func Test_getAddrByPublicKeyHex(t *testing.T) {

	tests := []struct {
		name          string
		hexPubblicKey string
		want          string
		wantErr       bool
	}{
		{
			"test 1",
			"0382e51eae6be5b247d7467727d0c6bf7bbe58e84cb6f913a8db28082c1c4f093d",
			"0361515c23a264d13c897a2ae74ab46211781b0e",
			false,
		},
		{
			"test 1",
			"0482e51eae6be5b247d7467727d0c6bf7bbe58e84cb6f913a8db28082c1c4f093df034be14690fada529a443c1bf648dcd18f621b8fa23acf7fec843d891d6978b",
			"0361515c23a264d13c897a2ae74ab46211781b0e",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAddrByPublicKeyHex(tt.hexPubblicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAddrByPublicKeyHexV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAddrByPublicKeyHexV2() = %v, want %v", got, tt.want)
			}
		})
	}
}
