/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// day5Cmd represents the day5 command
var day5Cmd = &cobra.Command{
	Use:   "day5",
	Short: "Run AOC day 5 challenge",
	Long:  `After the rearrangement procedure completes, what crate ends up on top of each stack?`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading input...")
		stacks := parseCrateStacks()
		stacks.ApplyMoves()
		ch1, ch2 := stacks.GetTops()
		fmt.Printf("Challenge 1: %s\n", ch1)
		fmt.Printf("Challenge 2: %s\n", ch2)
	},
}

type Stack []string

func (s *Stack) Push(e string) {
	*s = append(*s, e)
}

func (s *Stack) Size() int {
	return len(*s)
}

func (s *Stack) Pop() string {
	n := s.Size() - 1
	elem := (*s)[n]
	*s = (*s)[:n]
	return elem
}

func (s *Stack) PushMultiple(elems []string) {
	for _, e := range elems {
		*s = append(*s, e)
	}
}

func (s *Stack) PopMultiple(n int) []string {
	// 1-2-3-4-5
	// pop 2
	// s = 1-2-3
	// pop = 4-5
	// l = 5
	// result = [3,4]
	l := s.Size() // 5
	result := (*s)[l-n : l]
	// 3-4
	*s = (*s)[0 : l-n]
	return result
}

func (s *Stack) Invert() {
	for left, right := 0, len(*s)-1; left < right; left, right = left+1, right-1 {
		(*s)[left], (*s)[right] = (*s)[right], (*s)[left]
	}
}

type Move struct {
	Quantity  int
	FromIndex int
	ToIndex   int
}

type CrateStacks struct {
	stacks         []Stack
	stacksMultiple []Stack
	moves          []Move
}

func (c *CrateStacks) InvertStacks() {
	for _, stack := range (*c).stacks {
		stack.Invert()
	}
}

func (c *CrateStacks) CloneStack() {
	for _, stack := range (*c).stacks {
		clone := make(Stack, len(stack))
		copy(clone, stack)
		c.stacksMultiple = append(c.stacksMultiple, clone)
	}
}

func (c *CrateStacks) ApplyMove(move Move) {
	elems := c.stacksMultiple[move.FromIndex].PopMultiple(move.Quantity)
	c.stacksMultiple[move.ToIndex].PushMultiple(elems)
	for ; move.Quantity > 0; move.Quantity-- {
		elem := c.stacks[move.FromIndex].Pop()
		c.stacks[move.ToIndex].Push(elem)
	}
}

func (c *CrateStacks) ApplyMoves() {
	for _, move := range c.moves {
		c.ApplyMove(move)
	}
}

func (c *CrateStacks) GetTops() ([]string, []string) {
	result := []string{}
	resultMultiple := []string{}
	for _, stack := range c.stacks {
		result = append(result, stack[len(stack)-1])
	}
	for _, stack := range c.stacksMultiple {
		resultMultiple = append(resultMultiple, stack[len(stack)-1])
	}
	return result, resultMultiple
}

func parseStack(stackLine string, currentStacks *[]Stack) []Stack {
	// length = s*3 + (s - 1) = s*4 - 1
	// stacks = (len + 1) / 4
	// its either nums or chars (we ignore nums)
	stackLineRune := []rune(stackLine)
	stackIndex := 1
	currentStack := 0
	// if no stacks, init them
	if len(*currentStacks) == 0 {
		*currentStacks = make([]Stack, (len(stackLine)+1)/4)
	}
	for stackIndex < len(stackLine) {
		value := string(stackLineRune[stackIndex])
		if !strings.ContainsAny(value, " 123456789") {
			(*currentStacks)[currentStack].Push(value)
		}
		stackIndex += 4
		currentStack++
	}
	return *currentStacks
}

func parseMove(move string) Move {
	// move 1 from 4 to 1
	// move (\d) from (\d) to (\d)
	r := regexp.MustCompile("move (?P<Move>\\d+) from (?P<From>\\d) to (?P<To>\\d)")
	matches := r.FindStringSubmatch(move)
	moveIndex := r.SubexpIndex("Move")
	fromIndex := r.SubexpIndex("From")
	toIndex := r.SubexpIndex("To")
	quantity, _ := strconv.Atoi(matches[moveIndex])
	from, _ := strconv.Atoi(matches[fromIndex])
	to, _ := strconv.Atoi(matches[toIndex])

	return Move{Quantity: quantity, FromIndex: from - 1, ToIndex: to - 1}
}

func parseCrateStacks() CrateStacks {
	file, err := os.Open("input/day5.txt")
	defer file.Close()

	if err != nil {
		log.Fatal("Could not open file!")
	}
	result := CrateStacks{}

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)
	stacks := true // determine if parsing stacks or moves

	for scan.Scan() {
		scanResult := scan.Text()
		if scanResult == "" {
			stacks = false
			continue
		}
		if stacks {
			result.stacks = parseStack(scanResult, &result.stacks)
		} else {
			result.moves = append(result.moves, parseMove(scanResult))
		}
	}
	result.InvertStacks() // we invert each stack to append afterwards
	result.CloneStack()
	return result
}

func init() {
	rootCmd.AddCommand(day5Cmd)
}
