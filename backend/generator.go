// in main package for simple compile
package main

import (
	"fmt"
	"time"

	redis "gopkg.in/redis.v5"
)

const (
	GENERATION_INTERVAL = 500 * time.Millisecond
)

type Generator struct {
	q  string
	c  *redis.Client
	on bool
}

func NewGenerator(q string, c *redis.Client) *Generator {
	g := &Generator{
		c: c,
		q: q,
	}
	go g.generation()
	return g
}

func (g *Generator) Start() {
	g.on = true
}

func (g *Generator) Stop() {
	g.on = false
}

func (g *Generator) generation() {
	for range time.Tick(GENERATION_INTERVAL) {
		if !g.on {
			continue
		}
		g.c.RPush(g.q, fmt.Sprint(time.Now().UnixNano()))
	}
}
