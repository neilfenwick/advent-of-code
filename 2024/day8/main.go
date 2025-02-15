package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening file: %s", os.Args[1])
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	result := processFile(file)
	result.populateAntinodeMapPart1()
	result.trimMapToBounds(result.size.x, result.size.y)
	antiNodeCount := result.countUniqueAntinodes()

	fmt.Printf("Number of anti-nodes: %d\n", antiNodeCount)
}

func processFile(file io.Reader) *antiNodeMap {
	anm := createAntiNodeMap()
	maxXBound, maxYBound := 0, 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maxXBound = len(line) - 1
		antennae := parseLine(line)
		for _, antenna := range antennae {
			anm.addAntenna(antenna.name, antenna.x, maxYBound)
		}
		maxYBound++
	}
	anm.size = coordinate{maxXBound, maxYBound - 1}
	return anm
}

func parseLine(line string) []struct {
	name rune
	x    int
} {
	antennae := make([]struct {
		name rune
		x    int
	}, 0)

	for idx, char := range line {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			antennae = append(antennae, struct {
				name rune
				x    int
			}{name: char, x: idx})
		}
	}
	return antennae
}

type coordinate struct {
	x int
	y int
}

type antiNodeMap struct {
	size        coordinate
	antennaMap  map[rune][]coordinate
	antiNodeMap map[rune][]coordinate
}

func createAntiNodeMap() *antiNodeMap {
	return &antiNodeMap{
		antennaMap:  make(map[rune][]coordinate),
		antiNodeMap: make(map[rune][]coordinate),
	}
}

func (anm *antiNodeMap) addAntenna(name rune, x, y int) {
	anm.antennaMap[name] = append(anm.antennaMap[name], coordinate{x, y})
}

func (anm *antiNodeMap) populateAntinodeMapPart1() {
	for name := range anm.antennaMap {
		antennae := anm.antennaMap[name]
		if len(antennae) < 2 {
			continue
		}
		for idx, coord := range antennae {
			for _, otherAntenna := range antennae[idx+1:] {
				deltaX := coord.x - otherAntenna.x
				deltaY := coord.y - otherAntenna.y

				antiNode1Pos := coordinate{coord.x + deltaX, coord.y + deltaY}
				antiNode2Pos := coordinate{otherAntenna.x - deltaX, otherAntenna.y - deltaY}
				antiNodeSlice := anm.antiNodeMap[name]
				antiNodeSlice = append(antiNodeSlice, antiNode1Pos, antiNode2Pos)
				anm.antiNodeMap[name] = antiNodeSlice
			}
		}
	}
}

func (anm *antiNodeMap) trimMapToBounds(x, y int) {
	for name, antiNodes := range anm.antiNodeMap {
		indexesToRemove := make([]int, 0)
		for idx, antiNode := range antiNodes {
			if antiNode.x < 0 || antiNode.x > x || antiNode.y < 0 || antiNode.y > y {
				indexesToRemove = append(indexesToRemove, idx)
			}
		}

		// Remove the antiNodes in reverse order to avoid index shifting issues
		for i := len(indexesToRemove) - 1; i >= 0; i-- {
			idx := indexesToRemove[i]
			antiNodes[idx] = antiNodes[len(antiNodes)-1]
			antiNodes = antiNodes[:len(antiNodes)-1]
		}

		anm.antiNodeMap[name] = antiNodes
	}
}

func (anm *antiNodeMap) countUniqueAntinodes() int {
	uniqueNodes := make(map[coordinate]bool)
	for _, nodeCoords := range anm.antiNodeMap {
		for _, node := range nodeCoords {
			uniqueNodes[node] = true
		}
	}
	return len(uniqueNodes)
}

func (anm *antiNodeMap) printMap() {
	for name, antiNodes := range anm.antiNodeMap {
		fmt.Printf("Antenna %c\n", name)
		for _, node := range antiNodes {
			fmt.Printf("x: %d, y: %d\n", node.x, node.y)
		}
	}
}
