package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"

	redis "gopkg.in/redis.v5"
)

func main() {
	// expvar
	go http.ListenAndServe(":1234", nil)

	// cli args
	er := flag.Bool("getErrors", false, "show errors")
	ad := flag.String("redisAddress", "localhost:6379", "specify redis node address")
	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr: *ad,
	})
	if err := client.Ping().Err(); err != nil {
		log.Fatalln("Redis connect fail:", err)
	}

	g := NewGenerator("messages", client)
	p := NewProcessor("messages", "errors", client)
	e := NewElector(client, g, p)

	if *er {
		for _, m := range p.Errors() {
			fmt.Println(m)
		}
	} else {
		e.Run()
	}
}
