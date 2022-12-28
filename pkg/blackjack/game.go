package blackjack

import (
	"fmt"

	"github.com/JaydenTeoh/card-deck/pkg/deck"
)

type state int8

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 10
	}
	if opts.BlackjackPayout == 0.0 {
		opts.BlackjackPayout = 1.5
	}
	g.numDecks = opts.Decks
	g.numHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout
	return g
}

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Game struct {
	//configurations
	numDecks        int
	numHands        int
	blackjackPayout float64
	//game states
	state state
	deck  []deck.Card
	//player fields
	player    []deck.Card
	playerBet int
	balance   int
	//dealer fields
	dealer   []deck.Card
	dealerAI AI
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("It isn't currently any player's turn.")
	}
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
}

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.numDecks / 3 // arbitrary card number to signal that deck is running low

	//numHands rounds of Blackjack
	for i := 0; i < g.numHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.numDecks), deck.Shuffle) //create a shuffled n-decks game everytime cards run low (< min)
			shuffled = true
		}
		bet(g, ai, shuffled)
		deal(g)

		//check if dealer has hit Blackjack prematurely
		if Blackjack(g.dealer...) {
			endHand(g, ai)
			continue
		}

		//Player Turn
		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		//Dealer Turn
		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endHand(g, ai)
	}
	return g.balance
}

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func MoveStand(g *Game) {
	g.state++
}

// draw the top card of the deck
func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// Soft returns true if the score of a hand is a soft score; when ace is counted as 11 points
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

// Score will take in a hand of cards and return the best blackjack score
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	//cannot convert Ace from value 1 to 11 if total score of hand is already > 11
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			//change value of Ace card from 1 to 11
			return minScore + 10
		}
	}
	return minScore
}

// Returns true if player hits a blackjack
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	pBlackjack, dBlackjack := Blackjack(g.player...), Blackjack(g.dealer...)
	winnings := g.playerBet //amount that can be won per round (how much player bet)
	switch {
	case pBlackjack && dBlackjack:
		fmt.Println("Draw.")
		winnings = 0 //no winnings
	case dBlackjack:
		fmt.Println("You lose.")
		winnings = -winnings //lose bet
	case pBlackjack:
		fmt.Println("BLACKJACK!")
		winnings = int(float64(winnings) * g.blackjackPayout) //win bet * blackjack payout
	case pScore > 21:
		fmt.Println("You busted.")
		winnings = -winnings //lose bet
	case dScore > 21:
		fmt.Println("Dealer busted.") //win bet
	case pScore > dScore:
		fmt.Println("You win!") //win bet
	case dScore > pScore:
		fmt.Println("You lose.")
		winnings = -winnings //lose bet
	case dScore == pScore:
		fmt.Println("Draw.")
		winnings = 0 //no winnings
	}
	g.balance += winnings
	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}
