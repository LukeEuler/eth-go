package eg

import (
	"sync"

	"github.com/LukeEuler/dolly/log"
	"github.com/LukeEuler/dolly/net/rpc"

	"github.com/LukeEuler/eth-go/config"
)

var (
	once sync.Once
	node *rpc.Client
)

func initClient() {
	f := func() {
		conf := config.Get()
		var err error
		node, err = rpc.DialInsecureSkipVerify(conf.Net.URL, "", "", rpc.JSONRPCVersion2)
		if err != nil {
			log.Fatal(err)
		}
	}
	once.Do(f)
}
