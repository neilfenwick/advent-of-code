package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	var (
		file *os.File
		err  error
	)

	priorityStrategy = calculatePrioritiesPart1

	switch len(os.Args) {
	case 1:
		file = os.Stdin
	case 3:
		scoringStrategy, _ := strconv.Atoi(os.Args[2])
		if scoringStrategy > 0 {
			priorityStrategy = calculatePrioritiesPart2
		}
		fallthrough
	case 2:
		file, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Error opening file: %s", os.Args[1])
		}
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	ruckSacks := priorityStrategy(file)

	sum := 0
	for _, ruckSack := range ruckSacks {
		sum += ruckSack.getPriorityScore()
	}

	fmt.Printf("Priority sum: %d\n", sum)
}

type ruckSack struct {
	compartment1 []rune
	compartment2 []rune
	compartment3 []rune
	priorityItem rune
}

var priorityStrategy strategy

type strategy func(file io.Reader) []ruckSack

func (r *ruckSack) calculatePriorityItem() {
	/*
	   Start by sorting the items in each compartment and then comparing each
	   compartment by walking them in order, similar to merge sort, to find the first
	   item that exists in both compartments
	*/
	sort.Slice(r.compartment1, func(i, j int) bool { return r.compartment1[i] < r.compartment1[j] })
	sort.Slice(r.compartment2, func(i, j int) bool { return r.compartment2[i] < r.compartment2[j] })
	sort.Slice(r.compartment3, func(i, j int) bool { return r.compartment3[i] < r.compartment3[j] })

	var currentRune rune
	i, j, k := 0, 0, 0

	for i < len(r.compartment1) {
		currentRune = r.compartment1[i]

		for j < len(r.compartment2) {
			testRune := r.compartment2[j]
			if currentRune == testRune {

				if len(r.compartment3) > 0 {
					for k < len(r.compartment3) {
						testRuneInner := r.compartment3[k]
						if testRune == testRuneInner {
							goto foundPriorityItem
						} else if testRune > testRuneInner {
							k++
							continue
						} else if testRune < testRuneInner {
							break
						}
					}
					break
				} else {
					goto foundPriorityItem
				}
			} else if currentRune > testRune {
				j++
				continue
			} else if currentRune < testRune {
				break
			}
		}
		i++
	}

foundPriorityItem:
	r.priorityItem = currentRune
}

func (r *ruckSack) getPriorityScore() int {
	/*
	  Go runes for 'a-z' run from 97-122
	  Go runes for 'A-Z' run from 65-90
	  This function needs to return 1-26 for 'a-z' and 27-52 for 'A-Z'
	*/
	r.calculatePriorityItem()
	priorityScore := int(r.priorityItem)

	if 97 <= priorityScore && priorityScore <= 122 {
		return priorityScore - 96
	}

	return priorityScore - 64 + 26
}

func (r *ruckSack) String() string {
	return fmt.Sprintf(
		"RuckSack:\nCompartment1: %v\nCompartment2: %v\nCompartment3: %v\nPriorityItem: %s\n",
		string(r.compartment1),
		string(r.compartment2),
		string(r.compartment3),
		string(r.priorityItem))
}

func calculatePrioritiesPart1(file io.Reader) []ruckSack {
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)
	ruckSacks := []ruckSack{}

	for s.Scan() {
		pack := ruckSack{}
		items := s.Text()
		pack.compartment1 = []rune(items[0 : len(items)/2])
		pack.compartment2 = []rune(items[len(items)/2:])
		ruckSacks = append(ruckSacks, pack)
	}

	return ruckSacks
}

func calculatePrioritiesPart2(file io.Reader) []ruckSack {
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)
	ruckSacks := []ruckSack{}
	lines := make([]string, 3)
	i := 1

	for s.Scan() {

		lines[i-1] = s.Text()
		if i%3 == 0 {
			pack := ruckSack{}
			pack.compartment1 = []rune(lines[0])
			pack.compartment2 = []rune(lines[1])
			pack.compartment3 = []rune(lines[2])
			ruckSacks = append(ruckSacks, pack)
			i = 0
		}
		i++
	}
	return ruckSacks
}
