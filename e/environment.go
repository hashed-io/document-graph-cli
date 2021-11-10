package e

import (
	"context"
	"fmt"
	"sync"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Environment struct {
	A        *eos.API
	X        context.Context
	AppName  string
	Contract eos.AccountName
	User     eos.AccountName
	Pause    time.Duration
}

var once sync.Once
var Env *Environment

func E() *Environment {
	onceBody := func() {

		Env = &Environment{
			A:        eos.New(viper.GetString("EosioEndpoint")),
			X:        context.Background(),
			AppName:  viper.GetString("AppName"),
			Contract: eos.AN(viper.GetString("Contract")),
			User:     eos.AN(viper.GetString("UserAccount")),
			Pause:    viper.GetDuration("Pause"),
		}

		keyBag := &eos.KeyBag{}
		Env.A.SetSigner(keyBag)

		zap.S().Debug("Configured Environment object with sync.Once.Do")
	}
	once.Do(onceBody)
	return Env
}

func Pause(seconds time.Duration, headline, prefix string) {
	if headline != "" {
		fmt.Println(headline)
	}
	time.Sleep(seconds)
	fmt.Println()
}

func DefaultPause(headline string) {
	Pause(E().Pause, headline, "")
}
