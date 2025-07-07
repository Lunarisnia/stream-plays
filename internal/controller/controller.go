package controller

import (
	"fmt"

	"github.com/Lunarisnia/stream-plays/internal/keysim"
	"github.com/micmonay/keybd_event"
)

type Controller interface {
	CastVote(choice string)
	String() string
	Execute()
}

type Command struct {
	VoteCount int
	KeyCode   int
}

func NewCommand(keyCode int) *Command {
	return &Command{
		VoteCount: 0,
		KeyCode:   keyCode,
	}
}

type controllerImpl struct {
	keyService keysim.KeySim

	commandMap map[string]*Command
}

func NewController(keyService keysim.KeySim) Controller {
	return &controllerImpl{
		keyService: keyService,
		commandMap: map[string]*Command{
			"up":    NewCommand(keybd_event.VK_W),
			"down":  NewCommand(keybd_event.VK_S),
			"left":  NewCommand(keybd_event.VK_A),
			"right": NewCommand(keybd_event.VK_D),
			// TODO: Add command for the other action buttons
		},
	}
}

func (c *controllerImpl) String() string {
	str := ""
	for k, v := range c.commandMap {
		str += fmt.Sprintf("%v = %v\n", k, v.VoteCount)
	}
	return str
}

func (c *controllerImpl) CastVote(choice string) {
	if command, exist := c.commandMap[choice]; exist {
		command.VoteCount++
	}
}

func (c *controllerImpl) Execute() {
	finalCommand := &Command{}
	prevHighest := 0
	keyStr := ""
	for key, val := range c.commandMap {
		if val.VoteCount > prevHighest {
			prevHighest = val.VoteCount
			finalCommand = val
			keyStr = key
		}
	}
	if prevHighest == 0 {
		fmt.Println("Nothing is pressed")
		return
	}

	fmt.Printf("Pressing %v (%v votes)\n", keyStr, prevHighest)
	c.keyService.Press(finalCommand.KeyCode)
}
