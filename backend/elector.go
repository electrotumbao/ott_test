// in main package for simple compile
package main

import (
	"log"
	"time"

	"github.com/pborman/uuid"

	redis "gopkg.in/redis.v5"
)

const (
	LOCK_KEY = "generator_id"

	REFRESH_INTERVAL = 1 * time.Second
	EXPIRE_INTERVAL  = 3 * time.Second

	MODE_NONE = 0
	MODE_GEN  = 1
	MODE_PROC = 2
)

type Mode int8

type Elector struct {
	appID string
	mode  Mode
	c     *redis.Client
	gen   *Generator
	proc  *Processor
}

func NewElector(c *redis.Client, g *Generator, p *Processor) *Elector {
	return &Elector{
		appID: uuid.New(),
		c:     c,
		gen:   g,
		proc:  p,
		mode:  MODE_NONE,
	}
}

func (e *Elector) Run() {
	for range time.Tick(REFRESH_INTERVAL) {
		err := e.c.Eval(`
local id = redis.call("get",KEYS[1])
if id == false or id == ARGV[1]
then
    return redis.call("setex",KEYS[1],ARGV[2],ARGV[1])
else
    return false
end
`, []string{LOCK_KEY}, e.appID, time.Duration(EXPIRE_INTERVAL).Seconds()).Err()
		if err == nil {
			if e.mode != MODE_GEN {
				log.Println("Message generation mode")
				e.mode = MODE_GEN
				e.proc.Stop()
				e.gen.Start()
			}
		} else {
			if e.mode != MODE_PROC {
				log.Println("Message processing mode")
				e.mode = MODE_PROC
				e.gen.Stop()
				e.proc.Start()
			}
		}
	}
}
