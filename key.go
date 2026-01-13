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

	num := max(1, conf.Goroutine)
	var wg sync.WaitGroup

	for range num {
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
		l := luckyMaxLen(a)
		if l >= luckyLen {
			fmt.Printf("%s\t%s\t%d\n", p, a, l)
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

func luckyMaxLen(content string) int {
	if len(content) == 0 {
		return 0
	}
	a := content[0]
	lm := 1
	tempLm := lm
	for i := 1; i < len(content); i++ {
		if content[i] == a {
			tempLm++
			continue
		}
		// content[i] != a
		a = content[i]
		if tempLm > lm {
			lm = tempLm
		}
		tempLm = 1
	}
	if tempLm > lm {
		lm = tempLm
	}
	return lm
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
