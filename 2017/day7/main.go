package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	data "github.com/neilfenwick/advent-of-code/data_structures"
)

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic("Could not open input.txt")
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	discs := make(map[string]Disc, 0)
	for scanner.Scan() {
		text := scanner.Text()
		d := ParseDisc(text)
		discs[d.Name] = d
	}

	var t *data.GenericTree[Disc]
	nodeMap := make(map[string]*data.GenericTreeNode[Disc])
	for _, disc := range discs {
		var node *data.GenericTreeNode[Disc]
		if t == nil {
			t = data.NewGenericTree(disc)
			node = t.Root
		}
		if n, ok := nodeMap[disc.Name]; ok {
			node = n
		} else if node == nil {
			node = &data.GenericTreeNode[Disc]{Value: disc}
			nodeMap[disc.Name] = node
		}
		nodeMap[disc.Name] = node
		for _, childName := range disc.Children {
			child := discs[childName]
			childNode := nodeMap[child.Name]
			if childNode == nil {
				childNode = node.AddChild(child)
				nodeMap[child.Name] = childNode
			} else {
				childNode.Parent = node
				node.Children = append(node.Children, childNode)
			}
		}
	}

	root := t.Root
	for root.Parent != nil {
		root = root.Parent
	}
	name := root.Value.Name
	unbalancedName, weightDiff := GetUnbalanced(root)
	unbalancedNode := nodeMap[unbalancedName]
	unbalancedDisc := unbalancedNode.Value

	fmt.Printf("The bottom program is called: %s\n", name)
	fmt.Printf("Program '%+v' is unbalanced. Its weight is %d away from what it should be.\n", unbalancedDisc, weightDiff)
}

// GetUnbalanced recursively searches down the tree (depth-first), following branches that do not
// have Discs with a total weight equal to their siblings.
// Returns the name of the lowest-level unbalanced Disc, with the difference in its weight
func GetUnbalanced(root *data.GenericTreeNode[Disc]) (string, int) {
	name, weightDelta := searchForUnbalancedChildren(root)
	return name, weightDelta
}

func searchForUnbalancedChildren(startNode *data.GenericTreeNode[Disc]) (string, int) {
	type nodeWeights struct {
		names  []string
		weight int
	}

	children := startNode.Children
	var weights []nodeWeights
	weightMap := make(map[int]int)

	for _, child := range children {
		weight := getNodeWeight(child)
		if pos, ok := weightMap[weight]; ok {
			nodeWeight := weights[pos]
			nodeWeight.names = append(nodeWeight.names, child.Value.Name)
			weights[pos] = nodeWeight
		} else {
			var names []string
			names = append(names, child.Value.Name)
			nodeWeight := nodeWeights{names: names, weight: weight}
			weights = append(weights, nodeWeight)
			weightMap[weight] = len(weights) - 1
		}
	}

	// all weights equal, exit
	if len(weights) == 1 {
		return startNode.Value.Name, 0
	}

	sort.Slice(weights, func(i, j int) bool {
		return len(weights[i].names) < len(weights[j].names)
	})

	unbalancedName := weights[0].names[0]
	unbalancedDelta := weights[0].weight - weights[1].weight

	childNameSearch, childDelta := searchForUnbalancedChildren(findChildByName(children, unbalancedName))

	if childDelta == 0 {
		return unbalancedName, unbalancedDelta
	}

	return childNameSearch, childDelta
}

func findChildByName(children []*data.GenericTreeNode[Disc], name string) *data.GenericTreeNode[Disc] {
	for _, child := range children {
		if child.Value.Name == name {
			return child
		}
	}
	return nil
}

func getNodeWeight(node *data.GenericTreeNode[Disc]) int {
	var weight int
	for _, child := range node.Children {
		weight += getNodeWeight(child)
	}
	weight += node.Value.Weight
	return weight
}
