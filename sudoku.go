package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_ROWS = 9
	MAX_COLS = 9
	MAX_QUAD = 9
	MAX_VAL = 9
	MASK_1 = 0x01 << 0
	MASK_2 = 0x01 << 1
	MASK_3 = 0x01 << 2
	MASK_4 = 0x01 << 3
	MASK_5 = 0x01 << 4
	MASK_6 = 0x01 << 5
	MASK_7 = 0x01 << 6
	MASK_8 = 0x01 << 7
	MASK_9 = 0x01 << 8
	MASK_ALL = 511
)

var Mask2Val = map[uint16]int8 {
	MASK_1: 1,
	MASK_2: 2,
	MASK_3: 3,
	MASK_4: 4,
	MASK_5: 5,
	MASK_6: 6,
	MASK_7: 7,
	MASK_8: 8,
	MASK_9: 9,
}

var Val2Mask = map[int8]uint16 {
	0: MASK_ALL,
	1: MASK_1,
	2: MASK_2,
	3: MASK_3,
	4: MASK_4,
	5: MASK_5,
	6: MASK_6,
	7: MASK_7,
	8: MASK_8,
	9: MASK_9,
}

var count int = 0

type Masks struct {
	row uint16
	col uint16
	quad uint16
}

type Puzzle [MAX_ROWS][MAX_COLS]int8
type Board [MAX_ROWS]Masks

func getQuad(row int, col int) int {
	return ((row/3) * 3) + (col/3)
}

func initBoard(fn string) (Puzzle, Board) {
	var puzzle Puzzle
	var tboard Board

	// Initialize the Board
	for i, _ := range tboard {
		tboard[i] = Masks{row: MASK_ALL, col: MASK_ALL, quad: MASK_ALL,}
	}

	// Read the puzzle file
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal("Failed to read %s: %v\n", fn, err)
	}

	// Split the file data into a slice of strings (one per line)
	lines := strings.Split(string(data), "\n")
	if len(lines) > MAX_ROWS {
		lines = lines[:MAX_ROWS]
	}

	// Convert the individual lines into individual values
	for row, line := range lines {
		vals := strings.Split(line, " ")
		if len(vals) == 9 {
			for col, strc := range vals {
				val, _ := strconv.Atoi(strc)
				q := getQuad(row, col)
				puzzle[row][col] = (int8)(val)
				if val > 0 {
					if mask, ok := Val2Mask[(int8)(val)]; ok {
						// Remove the value from the list of possible values for the corresponding
						// row,, col, quad
						tboard[row].row = tboard[row].row &^ mask
						tboard[col].col = tboard[col].col &^ mask
						tboard[q].quad = tboard[q].quad &^ mask
					}
				}
			}
		}
	}

	return puzzle, tboard
}

func solve(table Puzzle, tboard Board, row int, col int, val int8) (Puzzle, Board, bool) {
	if val > 0 {
		q := getQuad(row, col)

		tmask, _ := Val2Mask[val]
		tboard[row].row = tboard[row].row &^ tmask
		tboard[col].col = tboard[col].col &^ tmask
		tboard[q].quad = tboard[q].quad &^ tmask
		table[row][col] = val
	}

	for row < MAX_ROWS {
		for col < MAX_COLS {

			if table[row][col] == 0 {
				q := getQuad(row, col)

				// Find possible values
				masks := tboard[row].row & tboard[col].col & tboard[q].quad

				var b int8 = 1
				for b <= MAX_VAL {
					tmask, _ := Val2Mask[b]
					if masks & tmask == tmask {
						if tp, tb, solved := solve(table, tboard, row, col, b); solved {
							return tp, tb, solved
						}
					}

					b++
				}

				return table, tboard, false
			}
			col++
		}
		// Reset cols back to 0 to start
		col = 0
		row++
	}

	return table, tboard, true
}

func main() {
	filename := flag.String("f", "", "Puzzle file")
	flag.Parse()

	if *filename == "" {
		log.Fatal("Need to specify a puzzle file")
	}

	solved := false

	// Read in the puzzle and initialize the board
	puzzle, board := initBoard(*filename)

	start := time.Now()
	if puzzle, board, solved = solve(puzzle, board, 0, 0, 0); solved {
		elapsed := time.Since(start)
		fmt.Printf("Time to solve: %s\n\n", elapsed)
		dumpBoard(puzzle)
	}
}

func dumpBoard(puzzle Puzzle) {
	for _, cols := range puzzle {
		for _, val := range cols {
			fmt.Printf("\t%v ", val)
		}
		fmt.Println("")
	}
	fmt.Printf("\n\n")
}

func dumpMasks(tboard Board) {
	for _, data := range tboard {
		fmt.Printf("%09b %09b %09b\n", data.row, data.col, data.quad)
	}
}
