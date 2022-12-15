/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

// day1Cmd represents the day1 command
var day1Cmd = &cobra.Command{
	Use:   "day1",
	Short: "Run AOC day 1 challenge",
	Long:  `Find the Elf carrying the most Calories. How many total Calories is that Elf carrying?`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading input...")
		input := loadInput()
		fmt.Println("Loading elfs...")
		elfs, max := loadElfs(input)
		sort.Sort(elfs)
		challenge2 := max
		for i := 1; i < 3; i++ {
			challenge2 += (*elfs)[i].TotalCalories
		}
		fmt.Printf("Challenge 1: %d\n", max)
		fmt.Printf("Challenge 2: %d\n", challenge2)
	},
}

type Elf struct {
	TotalCalories int
	calories      []int
}

type Elfs []Elf

const ELFSEPARATOR int = -100

func (e Elfs) Len() int           { return len(e) }
func (e Elfs) Less(i, j int) bool { return e[i].TotalCalories > e[j].TotalCalories }
func (e Elfs) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func loadInput() *[]int {
	file, err := os.Open("input/day1.txt")
	defer file.Close()

	if err != nil {
		log.Fatal("Could not open file!")
	}
	var result []int

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		scanResult, err := strconv.Atoi(scan.Text())
		if err != nil {
			scanResult = ELFSEPARATOR
		}
		result = append(result, scanResult)
	}

	return &result
}

func loadElfs(elfsAsInt *[]int) (*Elfs, int) {
	elfs := make(Elfs, 0)
	currentElf := Elf{TotalCalories: 0, calories: make([]int, 0)}
	max := 0
	for _, calories := range *elfsAsInt {
		if calories == ELFSEPARATOR {
			if currentElf.TotalCalories > max {
				max = currentElf.TotalCalories
			}
			elfs = append(elfs, currentElf)
			currentElf = Elf{}
			continue
		}
		currentElf.TotalCalories = currentElf.TotalCalories + calories
		currentElf.calories = append(currentElf.calories, calories)
	}
	return &elfs, max
}

func init() {
	rootCmd.AddCommand(day1Cmd)
}
