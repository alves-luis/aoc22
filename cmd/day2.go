/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// day2Cmd represents the day2 command
var day2Cmd = &cobra.Command{
	Use:   "day2",
	Short: "Run AOC day 2 challenge",
	Long:  `What would your total score be if everything goes exactly according to your strategy guide?`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading input...")
		rounds := loadRPSInput()
		fmt.Printf("Challenge 1: %d\n", rounds.MyScore())
		fmt.Printf("Challenge 2: %d\n", rounds.MyScoreV2())
	},
}

type Play string

const (
	ROCK     Play = "ROCK"
	PAPER    Play = "PAPER"
	SCISSORS Play = "SCISSORS"
)

func (p Play) Points() int {
	switch p {
	case ROCK:
		return 1
	case PAPER:
		return 2
	case SCISSORS:
		return 3
	default:
		return 0
	}
}

func (p Play) WinsAgainst() Play {
	switch p {
	case ROCK:
		return SCISSORS
	case PAPER:
		return ROCK
	case SCISSORS:
		return PAPER
	}
	return ROCK // never happens
}

func (p Play) LosesAgainst() Play {
	switch p {
	case ROCK:
		return PAPER
	case PAPER:
		return SCISSORS
	case SCISSORS:
		return ROCK
	}
	return ROCK // never happens
}

func (p Play) AchieveOutcome(outcome Outcome) Play {
	switch outcome {
	case DRAW:
		return p
	case LOSS:
		return p.WinsAgainst()
	case WIN:
		return p.LosesAgainst()
	}
	return ROCK // never happens
}

type Outcome string

const (
	WIN  Outcome = "WIN"
	DRAW Outcome = "DRAW"
	LOSS Outcome = "LOSS"
)

type Round struct {
	Opponent       Play
	Me             Play
	DesiredOutcome Outcome
}

func (r Round) Score() int {
	if r.Me.WinsAgainst() == r.Opponent {
		return r.Me.Points() + 6
	}
	if r.Me == r.Opponent {
		return r.Me.Points() + 3
	}
	return r.Me.Points()
}

func (r Round) ScoreV2() int {
	myPlay := r.Opponent.AchieveOutcome(r.DesiredOutcome)
	if myPlay.WinsAgainst() == r.Opponent {
		return myPlay.Points() + 6
	}
	if myPlay == r.Opponent {
		return myPlay.Points() + 3
	}
	return myPlay.Points()
}

type Rounds []Round

func (r Rounds) MyScore() int {
	result := 0
	for _, round := range r {
		result += round.Score()
	}
	return result
}

func (r Rounds) MyScoreV2() int {
	result := 0
	for _, round := range r {
		result += round.ScoreV2()
	}
	return result
}

func getPlay(parsedPlay string) Play {
	switch parsedPlay {
	case "A", "X":
		return ROCK
	case "B", "Y":
		return PAPER
	case "C", "Z":
		return SCISSORS
	default:
		return ROCK // never happens
	}
}

func getOutcome(parsedOutcome string) Outcome {
	switch parsedOutcome {
	case "X":
		return LOSS
	case "Y":
		return DRAW
	case "Z":
		return WIN
	default:
		return WIN // never happens
	}
}

func loadRPSInput() *Rounds {
	file, err := os.Open("input/day2.txt")
	defer file.Close()

	if err != nil {
		log.Fatal("Could not open file!")
	}
	var result Rounds

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		plays := strings.Fields(scan.Text())
		round := Round{Me: getPlay(plays[1]), Opponent: getPlay(plays[0]), DesiredOutcome: getOutcome(plays[1])}

		result = append(result, round)
	}

	return &result
}

func init() {
	rootCmd.AddCommand(day2Cmd)
}
