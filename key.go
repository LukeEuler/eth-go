package eg

import (
	"fmt"
	"sync"

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

	num := conf.Goroutine
	if num < 1 {
		num = 1
	}
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			generateLuckyAddr()
		}()
	}

	wg.Wait()
}

func generateLuckyAddr() {
	length := config.Get().Length
	for {
		k, err := key.NewKey()
		if err != nil {
			log.Entry.Fatal(err)
		}
		p, a := k.PrivateKey(), k.Address()
		l := luckyLength(a)
		if l >= length {
			fmt.Printf("%d %s %s\n", l, p, a)
		}
	}
}

func luckyLength(content string) int {
	if len(content) == 0 {
		return 0
	}
	a := content[0]
	for i := range content {
		if content[i] != a {
			return i
		}
	}
	return len(content)
}

func NewKeys() {
	conf := config.Get()
	if !conf.Keys.Enable || conf.Keys.Number <= 0 {
		return
	}
	for i := 0; i < conf.Keys.Number; i++ {
		k, err := key.NewKey()
		if err != nil {
			log.Entry.Fatal(err)
		}
		fmt.Printf("%s %s\n", k.PrivateKey(), k.Address())
	}
}
