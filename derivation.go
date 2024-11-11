package eg

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/LukeEuler/dolly/log"
	ecies "github.com/ecies/go/v2"

	"github.com/LukeEuler/eth-go/config"
	"github.com/LukeEuler/eth-go/key"
)

func Derivation() {
	conf := config.Get()
	if !conf.Derivation.Enable {
		return
	}
	if len(conf.Derivation.PrivateKey) != 0 {
		sk, err := key.NewKeyFromHex(conf.Derivation.PrivateKey)
		if err != nil {
			log.Entry.WithError(err).Fatal(err)
		}
		fmt.Printf("private key: %s\n", conf.Derivation.PrivateKey)
		fmt.Printf("public key: %s\n", sk.PublicKey(false))
		fmt.Printf("public key(compressed): %s\n", sk.PublicKey(true))
		fmt.Printf("address: %s\n", sk.Address())
		return
	}

	if len(conf.Derivation.PublicKey) == 0 {
		fmt.Println("There is no private or public key")
		return
	}

	pk, err := ecies.NewPublicKeyFromHex(conf.Derivation.PublicKey)
	if err != nil {
		log.Entry.Fatal(err)
	}

	stdPK := ecdsa.PublicKey{
		Curve: pk.Curve,
		X:     pk.X,
		Y:     pk.Y,
	}

	fmt.Printf("public key: %s\n", hex.EncodeToString(pk.Bytes(false)))
	fmt.Printf("public key(compressed): %s\n", hex.EncodeToString(pk.Bytes(true)))
	fmt.Printf("address: %s\n", key.GetAddrByPubkey(stdPK))
}
