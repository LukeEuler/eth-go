package key

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"math/big"

	ecies "github.com/ecies/go/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type KeyPair struct {
	priv *ecdsa.PrivateKey
}

func NewKey() (*KeyPair, error) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	return &KeyPair{key}, errors.WithStack(err)
}

func NewKeyFromBytes(bs []byte) (*KeyPair, error) {
	if len(bs) != 32 {
		return nil, errors.New("invalid bytes")
	}
	curve := crypto.S256()
	x, y := curve.ScalarBaseMult(bs)
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(bs),
	}
	return &KeyPair{priv}, nil
}

func NewKeyFromHex(hexKey string) (*KeyPair, error) {
	bs, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewKeyFromBytes(bs)
}

func (k *KeyPair) PrivateKey() string {
	return hex.EncodeToString(common.LeftPadBytes(k.priv.D.Bytes(), 32))
}

func (k *KeyPair) PublicKey(compressed bool) string {
	pk := ecies.PublicKey{
		Curve: k.priv.PublicKey.Curve,
		X:     k.priv.PublicKey.X,
		Y:     k.priv.PublicKey.Y,
	}
	return hex.EncodeToString(pk.Bytes(compressed))
}

func (k *KeyPair) Address() string {
	bs := PubkeyToEthAddressBytes(k.priv.PublicKey)
	return hex.EncodeToString(bs)
}

func (k *KeyPair) ToECDSA() *ecdsa.PrivateKey {
	return k.priv
}

func GetAddrByPubkey(publicKey ecdsa.PublicKey) string {
	bs := PubkeyToEthAddressBytes(publicKey)
	return hex.EncodeToString(bs)
}

func getAddrByPublicKeyHex(raw string) (string, error) {
	pk, err := ecies.NewPublicKeyFromHex(raw)
	if err != nil {
		return "", errors.WithStack(err)
	}
	stdPK := ecdsa.PublicKey{
		Curve: pk.Curve,
		X:     pk.X,
		Y:     pk.Y,
	}
	return GetAddrByPubkey(stdPK), nil
}
