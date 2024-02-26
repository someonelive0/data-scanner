package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

func opts() {
	servers := []string{"nats://192.168.30.59:4222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}
	opts := nats.GetDefaultOptions()
	opts.Url = strings.Join(servers, ",")
	opts.Verbose = true
	opts.Pedantic = true
}

func main() {
	servers := []string{"nats://192.168.30.59:4222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}

	nc, err := nats.Connect(
		strings.Join(servers, ","),
		nats.Name("API Name for client"),
		nats.UserInfo("idss", "natsidss"),
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
		log.Fatal(err)
	}
	defer nc.Close()

	// Do something with the connection
	mp := nc.MaxPayload()
	log.Printf("Maximum payload is %v bytes", mp)

	getStatusTxt := func(nc *nats.Conn) string {
		switch nc.Status() {
		case nats.CONNECTED:
			return "Connected"
		case nats.CLOSED:
			return "Closed"
		default:
			return "Other"
		}
	}
	log.Printf("The connection is %v\n", getStatusTxt(nc))

	// // Subscribe
	// sub, err := nc.SubscribeSync("updates")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Wait for a message
	// msg, err := sub.NextMsg(10 * time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := sub.Unsubscribe(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Use the response
	// log.Printf("Reply: %s", msg.Data)

	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe
	sub, err := nc.Subscribe("updates", func(m *nats.Msg) {
		wg.Done()
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := sub.Unsubscribe(); err != nil {
		log.Fatal(err)
	}

	// // Create a queue subscription on "updates" with queue name "workers"
	// if _, err := nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) {
	// 	wg.Done()
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	// Wait for a message to come in
	wg.Wait()

}
