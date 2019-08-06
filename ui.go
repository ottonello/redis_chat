package main

import tm "github.com/buger/goterm"
import "bufio"
import "os"
import "strings"
import "fmt"

type ChatUI struct {
	chatBox *tm.Box
}

func NewChatUI() *ChatUI {
	c := new(ChatUI)
	chatBox := tm.NewBox(100|tm.PCT, 20, 0)
	c.chatBox = chatBox
	return c
}

func (c *ChatUI) AskForInput(prompt string) string {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Printf(prompt)
	tm.Flush()
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")
	return name
}

func (c *ChatUI) Reset() {
	tm.Print(tm.MoveTo(c.chatBox.String(), 1, 1))
	tm.MoveCursor(1, 21)
	tm.Print("> ")
	tm.MoveCursor(3, 21)
}

func (c *ChatUI) Printf(format string, args ...interface{}) {
	fmt.Fprintf(c.chatBox, "%s: %s", args...)
}

func (c *ChatUI) Flush() {
	tm.Flush()
}
