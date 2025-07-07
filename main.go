package main

import (
	"context"
	"strings"
	"time"

	"github.com/Lunarisnia/socrates"
	"github.com/Lunarisnia/stream-plays/internal/controller"
	"github.com/Lunarisnia/stream-plays/internal/keysim"
	"github.com/Lunarisnia/stream-plays/internal/lcservice"
)

func main() {
	ctx := context.Background()
	lcService, err := lcservice.NewService(ctx)
	if err != nil {
		panic(err)
	}
	ks, err := keysim.NewKeySim()
	if err != nil {
		panic(err)
	}

	socrates.Init("lunarisnia")
	commandPrefix := "!mock "
	for {
		controllerService := controller.NewController(ks)

		youtubeContainer, err := lcService.Poll()
		if err != nil {
			panic(err)
		}
		for _, chat := range youtubeContainer.FindPrefix(commandPrefix) {
			command := strings.TrimPrefix(chat.Content, commandPrefix)
			controllerService.CastVote(command)
		}

		tiktokContainer, err := socrates.Monitor(2*time.Second, true)
		if err != nil {
			panic(err)
		}
		for _, chat := range tiktokContainer.FindPrefix(commandPrefix) {
			command := strings.TrimPrefix(chat.Content, commandPrefix)
			controllerService.CastVote(command)
		}
		controllerService.Execute()

		// moveHistory := chatContainer.FindPrefix("!mock ")
		// for _, move := range moveHistory {
		// 	fmt.Println(move.Username, "said", move.Content)
		// }
	}

}
