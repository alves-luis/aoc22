/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// day4Cmd represents the day4 command
var day4Cmd = &cobra.Command{
	Use:   "day4",
	Short: "Run AOC day 4 challenge",
	Long:  `In how many assignment pairs does one range fully contain the other?`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading input...")
		assignments := loadAssignmentPairs()
		ch1 := assignments.countContains()
		ch2 := assignments.countOverlaps()
		fmt.Printf("Challenge 1: %d\n", ch1)
		fmt.Printf("Challenge 2: %d\n", ch2)
	},
}

type Pair struct {
	Upper int
	Lower int
}

func (p *Pair) contains(pair Pair) bool {
	// pair A fully contains pair B if pairA.lower <= pairB.lower and pairA.upper >= pairB.upper
	return p.Lower <= pair.Lower && p.Upper >= pair.Upper
}

func (p *Pair) overlaps(pair Pair) bool {
	// 4-6 3-4 YES
	// 4-6 6-8 YES
	// 4-6 5-5 YES
	// 4-6 1-3 NO
	// 4-6 7-8 NO
	// 4-6 1-8 YES
	// pair A overlaps pair B if
	// Al <= Bu and Au >= Bl
	result := (p.Lower <= pair.Upper) && (p.Upper >= pair.Lower)
	return result
}

type Assignment struct {
	PairA Pair
	PairB Pair
}

func (a *Assignment) contained() bool {
	return a.PairA.contains(a.PairB) || a.PairB.contains(a.PairA)
}

func (a *Assignment) overlaped() bool {
	return a.PairA.overlaps(a.PairB)
}

type AssignmentPairs struct {
	Assignments    []Assignment
	FullyContained int
	Overlaped      int
}

func (a *AssignmentPairs) countContains() int {
	for _, assigns := range a.Assignments {
		if assigns.contained() {
			a.FullyContained++
		}
	}
	return a.FullyContained
}

func (a *AssignmentPairs) countOverlaps() int {
	for _, assigns := range a.Assignments {
		if assigns.overlaped() {
			a.Overlaped++
		}
	}
	return a.Overlaped
}

func convertLineToPair(line string) Pair {
	pairsBoundaries := strings.Split(line, "-")
	upper, _ := strconv.Atoi(pairsBoundaries[1])
	lower, _ := strconv.Atoi(pairsBoundaries[0])
	return Pair{Upper: upper, Lower: lower}
}

func convertLineToAssignment(line string) Assignment {
	pairs := strings.Split(line, ",")
	return Assignment{PairA: convertLineToPair(pairs[0]), PairB: convertLineToPair(pairs[1])}
}

func loadAssignmentPairs() AssignmentPairs {
	file, err := os.Open("input/day4.txt")
	defer file.Close()

	if err != nil {
		log.Fatal("Could not open file!")
	}
	result := AssignmentPairs{}

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		scanResult := scan.Text()
		result.Assignments = append(result.Assignments, convertLineToAssignment(scanResult))
	}
	return result
}

func init() {
	rootCmd.AddCommand(day4Cmd)
}
