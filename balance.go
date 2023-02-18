package eg

import (
	"fmt"
	"math/big"
	"strings"

	dcommon "github.com/LukeEuler/dolly/common"
	"github.com/LukeEuler/dolly/log"
	"github.com/pkg/errors"

	"github.com/LukeEuler/eth-go/common"
	"github.com/LukeEuler/eth-go/config"
)

func GetBalance() {
	conf := config.Get()
	if !conf.Balance.Enable {
		return
	}

	for _, address := range conf.Balance.Address {
		getBalance(address)
		fmt.Printf("%s/address/0x%s\n\n", conf.Net.Show, address)
	}
}

func getBalance(address string) {
	initClient()

	var res string
	err := node.SyncCall(&res, "eth_getBalance", "0x"+address, "latest")
	if err != nil {
		log.Fatal(errors.Wrapf(err, "address: %s", address))
	}

	amount, ok := big.NewInt(0).SetString(common.FormatHex(res), 16)
	if !ok {
		log.Entry.Fatalf("invalid value: %s", res)
	}
	fmt.Printf("%s balance:\n", address)
	ra, err := dcommon.Cut(amount.String(), 18, 8)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  %s %s ETH\n", amount, ra)

	conf := config.Get()
	if conf.Balance.Token.Enable {
		for _, item := range conf.Balance.Token.Contract {
			list := strings.SplitN(item, " ", 2)
			contract := list[0]
			name := "$"
			if len(list) > 1 {
				name = list[1]
			}
			data := balanceFunc + "000000000000000000000000" + address

			err = node.SyncCall(&res, "eth_call", buildBaseEthCallParams("0x"+contract, data), "latest")
			if err != nil {
				log.Fatal(err)
			}
			raw := common.FormatHex(res)
			amount = big.NewInt(0)
			if len(raw) > 0 {
				_, ok = amount.SetString(raw, 16)
				if !ok {
					log.Entry.Fatalf("invalid value: %s", res)
				}
			}
			fmt.Printf("  %s: %s %s\n", contract, amount, name)
		}
	}
}
