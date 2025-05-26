package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening file: %s", os.Args[1])
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Load file map
	diskMap := LoadDiskMap(file)

	// Iterate through the file map and defrag right-to-left
	Defrag(diskMap)

	// Compute checksum
}

type Block struct {
	FileId int64
	Start  int64  // Start offset
	Length int    // Length of the block
	IsFree bool   // True if block is free
	Prev   *Block // Previous block
	Next   *Block // Next block
}

type DiskMap struct {
	Head *Block
	Tail *Block
}

func LoadDiskMap(file io.Reader) *DiskMap {
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanRunes)

	var index int64 = 0
	var blockPosition int64 = 0
	var prev *Block = nil
	var diskMap *DiskMap = nil

	for s.Scan() {
		fileSize, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatalf("Unexpected error trying to read integer value from diskmap %v", err)
		}

		// Deal with the case of the first block
		if prev == nil {
			prev = &Block{FileId: index, Start: blockPosition, Length: fileSize}
			diskMap = &DiskMap{Head: prev}

			index++
			blockPosition = blockPosition + int64(fileSize)
			continue
		}

		// Append to the tail of the diskmap
		tail := &Block{
			FileId: index,
			Start:  blockPosition,
			Length: fileSize,
			IsFree: index%2 == 1,
			Prev:   prev,
		}

		prev.Next = tail
		diskMap.Tail = tail

		index++
		blockPosition = blockPosition + int64(fileSize)
	}
	return diskMap
}

func Defrag(diskMap *DiskMap) {
	// Find the first free block
	seekBlock := diskMap.Head.Next
	for !seekBlock.IsFree {
		seekBlock = seekBlock.Next
	}

	// Find the last non-free block
	var lastNonFreeBlock *Block = diskMap.Tail
	for lastNonFreeBlock.IsFree {
		lastNonFreeBlock = diskMap.Tail.Prev
	}

	defragInternal(diskMap, seekBlock, lastNonFreeBlock)
}

func defragInternal(diskMap *DiskMap, firstFree *Block, lastNonFree *Block) {
	// If the defrag has iterated to the end of the diskmap, we are done
	if firstFree.Start >= diskMap.Tail.Start {
		return
	}

	if firstFree.Length == lastNonFree.Length {
		// if the last file fits exactly into the free space, just move the block

		// Unlink lastNonFree from previous predecessor and short-circuit it to lastNonFree's succesor
		lastNonFree.Prev.Next = lastNonFree.Next

		// Move lastNonFree to new position, slicing out old free block
		lastNonFree.Prev = firstFree.Prev
		lastNonFree.Next = firstFree.Next
		firstFree.Prev.Next = lastNonFree
		firstFree.Next.Prev = lastNonFree

	} else if firstFree.Length > lastNonFree.Length {
		// if the last file fits entirely into the free space, split the free block and insert the file

	} else {
		// if the last file does not fit entirely into the free space, split the last file and move the first part
	}
}
