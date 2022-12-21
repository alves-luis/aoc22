/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// day3Cmd represents the day3 command
var day3Cmd = &cobra.Command{
	Use:   "day3",
	Short: "Run AOC day 3 challenge",
	Long:  `Find the item type that appears in both compartments of each rucksack. What is the sum of the priorities of those item types?`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading input...")
		rucksacks := loadRucksackInput()
		ch1, ch2 := rucksacks.sumPriority()
		fmt.Printf("Challenge 1: %d\n", ch1)
		fmt.Printf("Challenge 2: %d\n", ch2)
	},
}

/*
for each line
split string /2
add map of occurrences of string 0
go through string 1 and check map of occurences of string 0
if already exists, it is common so sum its priority
*/

type Rucksack struct {
	FirstCompartment  string
	SecondCompartment string
	Priority          int
}

func getRunePriority(r rune) int {
	// a-z 97-122 1-26
	// A-z 65-90 27-52
	if int(r) > 96 {
		return int(r) - 96
	}
	return int(r) - 38
}

func (r Rucksack) getPriority() int {
	occurences := make(map[rune]bool)
	for _, c := range r.FirstCompartment {
		occurences[c] = true
	}
	for _, c := range r.SecondCompartment {
		if occurences[c] {
			prio := getRunePriority(c)
			r.Priority = prio
			return prio
		}
	}
	return -1 // should never happen
}

type Group struct {
	Rucksacks []string
	Priority  int
}

func (g Group) getPriority() int {
	occurencesA := make(map[rune]bool)
	for _, c := range g.Rucksacks[0] {
		occurencesA[c] = true
	}
	occurencesB := make(map[rune]bool)
	for _, c := range g.Rucksacks[1] {
		occurencesB[c] = true
	}
	for _, c := range g.Rucksacks[2] {
		if occurencesA[c] && occurencesB[c] {
			prio := getRunePriority(c)
			g.Priority = prio
			return prio
		}
	}
	return -1 // should never happen
}

type Rucksacks struct {
	Rucksacks   []Rucksack
	Groups      []Group
	RucksackSum int
	GroupSum    int
}

func (r Rucksacks) sumPriority() (int, int) {
	for _, ruck := range r.Rucksacks {
		r.RucksackSum += ruck.getPriority()
	}
	for _, group := range r.Groups {
		r.GroupSum += group.getPriority()
	}
	return r.RucksackSum, r.GroupSum
}

func loadRucksackInput() Rucksacks {
	file, err := os.Open("input/day3.txt")
	defer file.Close()

	if err != nil {
		log.Fatal("Could not open file!")
	}
	result := Rucksacks{}
	currentGroup := Group{}

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		scanResult := scan.Text()
		length := len(scanResult)
		rucksack := Rucksack{FirstCompartment: scanResult[0 : length/2], SecondCompartment: scanResult[length/2:]}
		result.Rucksacks = append(result.Rucksacks, rucksack)
		currentGroup.Rucksacks = append(currentGroup.Rucksacks, scanResult)
		if len(currentGroup.Rucksacks) == 3 {
			result.Groups = append(result.Groups, currentGroup)
			currentGroup = Group{}
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(day3Cmd)
}
