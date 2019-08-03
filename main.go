package main

import "fmt"
import "bufio"
import "os"
import "strings"
import "github.com/go-redis/redis"
import "github.com/satori/go.uuid"
import "encoding/json"
import tm "github.com/buger/goterm"

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

func createBox() *tm.Box {
	chatBox := tm.NewBox(100|tm.PCT, 20, 0)
	return chatBox
}

func main() {
	clientId := uuid.Must(uuid.NewV4())

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:30379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pubsub := client.Subscribe("mychannel")

	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Printf("Type in your name: ")
	tm.Flush()
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")

	c := make(chan Message)
	history := make([]Message, 0)

	tm.Clear()
	chatBox := createBox()
	tm.Print(tm.MoveTo(chatBox.String(), 1, 1))
	tm.MoveCursor(1, 21)
	tm.Print("> ")
	tm.MoveCursor(3, 21)
	tm.Flush()

	for true {
		tm.Clear()
		chatBox := createBox()
		tm.MoveCursor(1, 21)
		tm.Print("> ")
		tm.MoveCursor(3, 21)

		go read(client, name, clientId)
		go read_redis(c, pubsub)
		message := <-c
		if message.ClientId != clientId {
		}

		history = append(history, message)

		for index := range history {
			message := history[index]
			fmt.Fprintf(chatBox, "%s: %s", message.Name, message.Message)
		}
		if len(history) >= 18 {
			history = history[1:]
		}
		tm.Print(tm.MoveTo(chatBox.String(), 1, 1))
		tm.MoveCursor(3, 21)
		tm.Flush()
	}
}
