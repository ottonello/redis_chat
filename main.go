package main

import "bufio"
import "os"
import "github.com/go-redis/redis"
import "github.com/satori/go.uuid"
import "encoding/json"

type Message struct {
	Name     string    `json:"name"`
	ClientId uuid.UUID `json:"clientId"`
	Message  string    `json:"message"`
}

func read(client *redis.Client, name string, clientId uuid.UUID) {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	message := Message{Name: name, Message: str, ClientId: clientId}
	serialized, _ := json.Marshal(message)
	client.Publish("mychannel", string(serialized))
}

func read_redis(c chan Message, pubsub *redis.PubSub) {
	msg, err := pubsub.ReceiveMessage()
	if err != nil {
		panic(err)
	}
	var message Message
	bytes := []byte(msg.Payload)
	json.Unmarshal(bytes, &message)

	c <- message
}

func main() {
	clientId := uuid.Must(uuid.NewV4())

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:30379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pubsub := client.Subscribe("mychannel")

	ui := NewChatUI()
	name := ui.AskForInput("Type in your name: ")

	c := make(chan Message)
	history := make([]Message, 0)

	ui.Reset()
	ui.Flush()
	for true {
		ui.Reset()
		go read(client, name, clientId)
		go read_redis(c, pubsub)
		message := <-c
		if message.ClientId != clientId {
		}

		history = append(history, message)

		for index := range history {
			message := history[index]
			ui.Printf("%s: %s", message.Name, message.Message)
		}
		if len(history) >= 18 {
			history = history[1:]
		}
		ui.Flush()
	}
}
