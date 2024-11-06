package eg

import (
	"fmt"
	"sync"

	"github.com/LukeEuler/dolly/log"

	"github.com/LukeEuler/eth-go/config"
	"github.com/LukeEuler/eth-go/key"
)

var luckyLen int

// 靓号生成器
func Lucky() {
	conf := config.Get()
	if !conf.Lucky {
		return
	}
	luckyLen = config.Get().Length

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
	for {
		k, err := key.NewKey()
		if err != nil {
			log.Entry.Fatal(err)
		}
		p, a := k.PrivateKey(), k.Address()
		l1 := luckyPreLen(a)
		l2 := luckySufLen(a)
		if l1 >= luckyLen || l2 >= luckyLen {
			fmt.Printf("%s\t%s\t%d %d\n", p, a, l1, l2)
		}
	}
}

func luckyPreLen(content string) int {
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

func luckySufLen(content string) int {
	if len(content) == 0 {
		return 0
	}
	ll := len(content)
	a := content[ll-1]
	for i := len(content) - 2; i >= 0; i-- {
		if a != content[i] {
			return ll - 1 - i
		}
	}
	return ll
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
