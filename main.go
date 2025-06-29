package main

import (
	"fmt"
	"strings"

	"github.com/Lunarisnia/socrates"
	"github.com/Lunarisnia/stream-plays/internal/keysim"
	"github.com/micmonay/keybd_event"
)

func main() {
	ks, err := keysim.NewKeySim()
	if err != nil {
		panic(err)
	}
	fmt.Println("Keysim initiated")
	socrates.Init("lunarisnia")
	for {
		chatContainer, err := socrates.MonitorThreshold(1, "!move ", true)
		if err != nil {
			panic(err)
		}
		moveHistory := chatContainer.FindPrefix("!move ")

		decisionCounter := map[string]int{
			"up":    0,
			"down":  0,
			"left":  0,
			"right": 0,
		}
		for _, chat := range moveHistory {
			command := strings.ToLower(strings.TrimPrefix(chat.Content, "!move "))
			switch command {
			case "up":
				decisionCounter["up"]++
			case "down":
				decisionCounter["down"]++
			case "left":
				decisionCounter["left"]++
			case "right":
				decisionCounter["right"]++
			default:
				fmt.Println("Ngetik yang bener woy!")
			}
		}

		finalCommand := ""
		prevHighest := 0
		for key, val := range decisionCounter {
			if val > prevHighest {
				prevHighest = val
				finalCommand = key
			}
		}

		keyMap := map[string]int{
			"up":    keybd_event.VK_W,
			"down":  keybd_event.VK_S,
			"left":  keybd_event.VK_A,
			"right": keybd_event.VK_D,
		}
		fmt.Printf("Pressing %v (%v votes)\n", finalCommand, prevHighest)
		ks.Press(keyMap[finalCommand])

	}
}
