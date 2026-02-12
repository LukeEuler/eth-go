package config

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/LukeEuler/dolly/log"

	"github.com/LukeEuler/eth-go/common"
)

var conf = new(config)

type config struct {
	Lucky     bool `toml:"lucky"`
	Goroutine int  `toml:"goroutine"`
	Length    int  `toml:"length"`
	Net       struct {
		ChainID int64  `toml:"chain_id"`
		URL     string `toml:"url"`
		Show    string `toml:"show"`
	} `toml:"net"`
	Keys struct {
		Enable bool `toml:"enable"`
		Number int  `toml:"number"`
	} `toml:"keys"`
	Derivation struct {
		Enable     bool   `toml:"enable"`
		PrivateKey string `toml:"private_key"`
		PublicKey  string `toml:"public_key"`
	} `toml:"derivation"`
	Balance struct {
		Enable  bool     `toml:"enable"`
		Address []string `toml:"address"`
		Token   struct {
			Enable   bool     `toml:"enable"`
			Contract []string `toml:"contract"`
		} `toml:"token"`
	} `toml:"balance"`
	KeyPair  map[string]string      `toml:"key_pair"`
	Transfer map[string][]*Transfer `toml:"transfer"`
}

type Transfer struct {
	Enable   bool   `toml:"enable"`
	From     string `toml:"from"`
	To       string `toml:"to"`
	GasLimit uint64 `toml:"gas_limit"`
	// MaxPriorityFeePerGas uint64 `toml:"max_priority_fee_per_gas"`
	// MaxFeePerGas         uint64 `toml:"max_fee_per_gas"`

	Contract string `toml:"contract"`
	Data     string `toml:"data"`
	Amount   string `toml:"amount"`

	// mpc
	RSV string `toml:"rsv"`
}

func (t *Transfer) format() {
	t.From = common.FormatHex(t.From)
	t.To = common.FormatHex(t.To)
	t.Contract = common.FormatHex(t.Contract)
}

func New(configPath string) {
	bs, err := os.ReadFile(configPath)
	if err != nil {
		log.Entry.Fatal(err)
	}
	_, err = toml.Decode(string(bs), conf)
	if err != nil {
		log.Entry.Fatal(err)
	}
	conf.format()
}

func Get() *config {
	return conf
}

func (c *config) format() {
	for i := range c.Balance.Address {
		c.Balance.Address[i] = common.FormatHex(c.Balance.Address[i])
	}
	for i := range c.Balance.Token.Contract {
		list := strings.SplitN(c.Balance.Token.Contract[i], " ", 2)
		temp := common.FormatHex(list[0])
		if len(list) > 1 {
			temp += " " + list[1]
		}
		c.Balance.Token.Contract[i] = temp
	}
	temp := make(map[string]string, len(c.KeyPair))
	for k, v := range c.KeyPair {
		temp[common.FormatHex(k)] = v
	}
	c.KeyPair = temp

	for _, tfs := range c.Transfer {
		for _, tf := range tfs {
			tf.format()
		}
	}
}
