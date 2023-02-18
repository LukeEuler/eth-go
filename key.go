package eg

import (
	"fmt"

	"github.com/LukeEuler/dolly/log"

	"github.com/LukeEuler/eth-go/config"
	"github.com/LukeEuler/eth-go/key"
)

// 靓号生成器
func Lucky() {
	conf := config.Get()
	if !conf.Lucky {
		return
	}
	for {
		k, err := key.NewKey()
		if err != nil {
			log.Fatal(err)
		}
		p, a := k.PrivateKey(), k.Address()
		if a[0] == a[1] && a[0] == a[2] && a[0] == a[3] && a[0] == a[4] && a[0] == a[5] {
			fmt.Printf("%s %s\n", p, a)
		}
	}
}

func NewKeys() {
	conf := config.Get()
	if !conf.Keys.Enable || conf.Keys.Number <= 0 {
		return
	}
	for i := 0; i < conf.Keys.Number; i++ {
		k, err := key.NewKey()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s\n", k.PrivateKey(), k.Address())
	}
}
