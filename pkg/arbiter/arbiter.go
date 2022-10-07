package arbiter

import (
	"log"
)

type Arbiter struct {
}

func NewArbiter() Arbiter {
	return Arbiter{}
}

func (arb *Arbiter) Start(broker chan string) {
	log.Println("starting bazaar arbiter")

	for msg := range broker {
		if msg == "refresh" {
			arb.Search()
		}
	}
}

func (arb *Arbiter) Search() {
	log.Println("searching possible arbitration paths")
}
