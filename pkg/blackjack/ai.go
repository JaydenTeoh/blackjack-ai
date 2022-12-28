package blackjack

import (
	"fmt"

	"github.com/JaydenTeoh/card-deck/pkg/deck"
)

type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hand [][]deck.Card, dealer []deck.Card)
}

type dealerAI struct{}

func (ai dealerAI) Bet(shuffled bool) int {
	//noop
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	// If dealer score <= 16, dealer hit
	// If dealer has a soft 17 (Ace + Ten/Jack/Queen/King), dealer hit
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	//noop
}

func PlayerAI() AI {
	return playerAI{}
}

type playerAI struct{}

func (ai playerAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled.")
	}
	fmt.Println("How much would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai playerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer, "**HIDDEN**")
		fmt.Println("What will you do? h - hit, s - stand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option - Enter 'h' to hit or 's' to stand.")
		}
	}
}

func (ai playerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}
