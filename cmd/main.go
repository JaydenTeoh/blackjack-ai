package main

import (
	"fmt"

	"github.com/JaydenTeoh/blackjack-ai/pkg/blackjack"
)

func main() {
	game := blackjack.New()
	winnings := game.Play(blackjack.PlayerAI())
	fmt.Println(winnings)
}
