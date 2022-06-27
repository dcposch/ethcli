package act

import (
	"log"
	"time"

	"dcposch.eth/cli/eth"
)

var (
	client   *eth.Client
	state    State
	renderer func(*State)
	queue    chan Action
)

func Init(_client *eth.Client, _renderer func(*State)) {
	client = _client
	renderer = _renderer
	queue = make(chan Action)

	go run()
}

func Dispatch(a Action) {
	log.Printf("action %#v", a)
	queue <- a
}

func run() {
	reloadChainState()
	render()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case a := <-queue:
			a.Run()
		case t := <-ticker.C:
			if t.Unix()%10 == 0 {
				reloadChainState()
			} else {
				reloadTabState()
			}
			render()
		}
	}
}
