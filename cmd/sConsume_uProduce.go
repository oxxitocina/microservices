package main

import (
	"rabbitMQhelloworld/cmd/internal/server"
	"rabbitMQhelloworld/cmd/internal/user1"
	"rabbitMQhelloworld/cmd/internal/user2"
	"time"
)

func main() {
	//to check users sent message and server consume it right away
	/*user1.UserProduce1()
	user2.UserProduce2()
	go server.ServerConsume()*/

	//Press CTRL+C to cancel the server consume
	server.ServerProduce() // Run server producer in a goroutine
	time.Sleep(50)
	//here the message to user1 can't reach because they are chained
	user1.UserConsume() // Run user1 consumer in a goroutine
	user2.UserConsume()

	time.Sleep(50)
}
