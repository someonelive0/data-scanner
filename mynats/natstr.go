package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var arg_host = flag.String("host", "localhost", "nats server hostname or ip")
var arg_user = flag.String("user", "", "nats user")
var arg_password = flag.String("password", "", "")

func init() {
	flag.Parse()
}

// func opts() {
// 	servers := []string{"nats://" + *arg_host + ":4222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}
// 	opts := nats.GetDefaultOptions()
// 	opts.Url = strings.Join(servers, ",")
// 	opts.Verbose = true
// 	opts.Pedantic = true
// }

func main() {
	servers := []string{"nats://" + *arg_host + ":4222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}

	nc, err := nats.Connect(
		strings.Join(servers, ","),
		nats.Name("natscli of mime 0.1.0"),
		nats.UserInfo(*arg_user, *arg_password),
		nats.Timeout(10*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(5),
		nats.NoEcho(),
		nats.ReconnectWait(10*time.Second),
		nats.ReconnectBufSize(5*1024*1024),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			// handle disconnect error event
			log.Printf("DisconnectErrHandler client disconnected: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			// handle reconnect event
			log.Println("ReconnectHandler client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("ClosedHandler client closed")
		}),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Printf("DiscoveredServersHandler client discover")
			log.Printf("Known servers: %v\n", nc.Servers())
			log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Printf("ErrorHandler Error: %v", err)
		}), // logSlowConsumer
	)
	if err != nil {
		log.Fatal("Connect failed: ", err)
	}
	defer nc.Close()

	// Do something with the connection
	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "kvs-tmp",
		Subjects: []string{"foo-kvs"},
	})

	js.Publish("foo-kvs", []byte("Hello JS!"))

	// Publish messages asynchronously.
	for i := 0; i < 500; i++ {
		js.PublishAsync("foo-kvs", []byte("Hello JS Async!"))
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}

	fmt.Println("END")
}
