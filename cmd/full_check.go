package main

import (
	"rabbitMQhelloworld/cmd/internal/server"
	"rabbitMQhelloworld/cmd/internal/user1"
	"rabbitMQhelloworld/cmd/internal/user2"
	"time"
)

func main() {
	server.ServerProduce()
	user1.UserProduce1()
	user2.UserProduce2()
	server.ServerConsume()

	//Press CTRL+C to cancel the server consume
	server.ServerProduce()
	time.Sleep(50)
	user1.UserConsume()
	user2.UserConsume()

}
