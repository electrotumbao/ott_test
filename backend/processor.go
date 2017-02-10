// in main package for simple compile
package main

import (
	"log"
	"math/rand"
	"time"

	redis "gopkg.in/redis.v5"
)

const (
	TRIGGER_CHECK_INTERVAL = 500 * time.Millisecond
)

type Processor struct {
	mq, eq string
	c      *redis.Client
	on     bool
	rnd    *rand.Rand
}

func NewProcessor(mq, eq string, c *redis.Client) *Processor {
	p := &Processor{
		mq:  mq,
		eq:  eq,
		c:   c,
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	go p.processing()
	return p
}

func (p *Processor) Start() {
	p.on = true
}

func (p *Processor) Stop() {
	p.on = false
}

func (p *Processor) Errors() []string {
	size := p.c.LLen(p.eq).Val()
	mm := p.c.LRange(p.eq, 0, size).Val()
	p.c.LTrim(p.eq, 1, 0)
	return mm
}

func (p *Processor) processing() {
	for {
		if !p.on {
			time.Sleep(TRIGGER_CHECK_INTERVAL)
			continue
		}

		res, err := p.c.BLPop(0, p.mq).Result()

		if err != nil {
			log.Println("Error on fetch message:", err)
			continue
		}
		msg := res[1]

		// 5% random error emulation
		i := p.rnd.Intn(100)
		if i < 5 {
			p.c.RPush(p.eq, msg)
			log.Println("Message corrupted. Skip.")
		} else {
			log.Println("Message processed:", msg)
		}
	}
}
