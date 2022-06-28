package act

import (
	"crypto/ecdsa"
	"log"
	"time"

	"dcposch.eth/cli/eth"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	client   *eth.Client
	state    State
	renderer func(*State)
	queue    chan Action
)

func Init(_client *eth.Client, _privateKey *ecdsa.PrivateKey, _renderer func(*State)) {
	client = _client
	renderer = _renderer
	queue = make(chan Action, 1)
	setPrivateKey(_privateKey)

	go run()
}

func setPrivateKey(prv *ecdsa.PrivateKey) {
	state.Chain.PrivateKey = prv
	if prv != nil {
		pub := crypto.PubkeyToAddress(prv.PublicKey)
		state.Chain.Account.Addr = pub
		log.Printf("recovering address from privkey %s %s", pub, state.Chain.Account.Addr)
	}
}

func Dispatch(a Action) {
	select {
	case queue <- a:
		log.Printf("action %#v", a)
	default:
		log.Printf("DROPPING action %#v", a)
	}
}

func run() {
	reloadChainState()

	tickChainState := time.NewTicker(time.Second * 10)
	tickTxState := time.NewTicker(time.Second * 2)
	for {
		select {
		case a := <-queue:
			a.Run()
		case <-tickChainState.C:
			reloadChainState()
		case <-tickTxState.C:
			reloadTxState()
		}
	}
}
