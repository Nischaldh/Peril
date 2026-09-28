package main

import (
	"fmt"

	"github.com/Nischaldh/Peril/internal/gamelogic"
	"github.com/Nischaldh/Peril/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState){
	return func(ps routing.PlayingState) {
		defer fmt.Print(">")
		gs.HandlePause(ps)
	}
}