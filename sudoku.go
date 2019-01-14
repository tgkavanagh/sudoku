package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_ROWS = 9
	MAX_COLS = 9
	MAX_QUAD = 9
	MAX_VAL  = 9
	MASK_1   = 0x01 << 0
	MASK_2   = 0x01 << 1
	MASK_3   = 0x01 << 2
	MASK_4   = 0x01 << 3
	MASK_5   = 0x01 << 4
	MASK_6   = 0x01 << 5
	MASK_7   = 0x01 << 6
	MASK_8   = 0x01 << 7
	MASK_9   = 0x01 << 8
	MASK_ALL = 511
)

var Mask2Val = map[uint16]int8{
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

var Val2Mask = map[int8]uint16{
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
	row  uint16
	col  uint16
	quad uint16
}

type CellData struct {
	row int
	col int
	quad int
	pvals int
	val int8
}

type Puzzle []CellData
type Board []Masks

//var puzzle Puzzle
//var board [MAX_ROWS]Masks

type ByPval Puzzle
func (pv ByPval) Len() int { return len(pv) }
func (pv ByPval) Less(i, j int) bool { return pv[i].pvals < pv[j].pvals }
func (pv ByPval) Swap(i, j int) { pv[i], pv[j] = pv[j], pv[i]}

type ByPos Puzzle
func (pv ByPos) Len() int { return len(pv) }
func (pv ByPos) Swap(i, j int) { pv[i], pv[j] = pv[j], pv[i] }
func (pv ByPos) Less(i, j int) bool {
	if pv[i].row == pv[j].row {
		return pv[i].col < pv[j].col
	}

	return pv[i].row < pv[j].row
}

func getQuad(row int, col int) int {
	return ((row / 3) * 3) + (col / 3)
}

func calcBitsSet(mask uint16) int {
	var i uint = 0
	count := 0
	for i < 16 {
		if (mask >> i) & 1 == 1 {
			count++
		}
		i++
	}

	return count
}

func (puz Puzzle)calcPvals(board Board) {
	for i, cell := range puz {
		row := cell.row
		col := cell.col
		quad := cell.quad

		if cell.val == 0 {
			mask := board[row].row & board[col].col & board[quad].quad
			puz[i].pvals = calcBitsSet(mask)
		} else {
			puz[i].pvals = 0
		}
	}
}

func (puz Puzzle)solve(idx int, board Board) (Board, bool) {
	count++
	for idx < (MAX_ROWS*MAX_COLS) {
		if puz[idx].val == 0 {
			row := puz[idx].row
			col := puz[idx].col
			quad := puz[idx].quad

			// Find possible values
			masks := board[row].row & board[col].col & board[quad].quad

			var b int8 = 1
			for b <= MAX_VAL {
				tmask, _ := Val2Mask[b]
				if masks&tmask == tmask {
					or := board[row].row
					oc := board[col].col
					oq := board[quad].quad

					board[row].row = board[row].row &^ tmask
					board[col].col = board[col].col &^ tmask
					board[quad].quad = board[quad].quad &^ tmask

					//fmt.Printf("Solving for idx %d val %d\n", idx, b)
					if tb, solved := puz.solve(idx+1, board); solved {
						puz[idx].val = b
						return tb, solved
					}

					board[row].row = or
					board[col].col = oc
					board[quad].quad = oq
				}

				b++
			}

			if idx == 17 {
				dumpPvals(board)
			}
			return board, false
		}

		idx++
	}

	return board, true
}

func initBoard(fn string) (Puzzle, Board) {
	var puzzle Puzzle
	var board Board

	// Initialize the Board
	i := 0
	for i < MAX_ROWS {
		newBcell := Masks{row: MASK_ALL, col: MASK_ALL, quad: MASK_ALL}
		board = append(board, newBcell)
		i++
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
		if len(vals) == MAX_COLS {
			for col, strc := range vals {
				val, _ := strconv.Atoi(strc)
				quad := getQuad(row, col)

				newCell := CellData{
					row: row,
					col: col,
					quad: quad,
					pvals: 0,
					val: (int8)(val),
				}

				puzzle = append(puzzle, newCell)

				if val > 0 {
					if mask, ok := Val2Mask[(int8)(val)]; ok {
						// Remove the value from the list of possible values for the corresponding
						// row,, col, quad
						board[row].row = board[row].row &^ mask
						board[col].col = board[col].col &^ mask
						board[quad].quad = board[quad].quad &^ mask
					}
				}
			}
		}
	}

	puzzle.calcPvals(board)

	return puzzle, board
}

func main() {
	filename := flag.String("f", "", "Puzzle file")
	flag.Parse()

	if *filename == "" {
		log.Fatal("Need to specify a puzzle file")
	}

	// Read in the puzzle and initialize the board
	puzzle, board := initBoard(*filename)
	sort.Sort(ByPval(puzzle))

	start := time.Now()
	if _, solved := puzzle.solve(0, board); solved {
		elapsed := time.Since(start)
		fmt.Printf("Time to solve: %s (recursive calls: %d)\n\n", elapsed, count)
		sort.Sort(ByPos(puzzle))
		puzzle.dumpBoard()
	} else {
		fmt.Printf("Failed to solve puzzle: %d\n", count)
	}
}

func (puz Puzzle)dumpBoard() {
	prevRow := 0
	for _, data := range puz {
		if data.row != prevRow {
			prevRow = data.row
			fmt.Println("")
		}

		fmt.Printf("\t%v ", data.val)
	}
	fmt.Printf("\n\n")
}

func (puz Puzzle)dumpTable(board Board) {
	for i, cell := range puz {
		row := puz[i].row
		col := puz[i].col
		quad := puz[i].quad

		// Find possible values
		masks := board[row].row & board[col].col & board[quad].quad

		fmt.Printf("%d: %v %09b\n", i, cell, masks)
	}
}

func dumpPvals(board Board) {
		r := 0

		for r < MAX_ROWS {
			c := 0
			for c < MAX_COLS {
				q := getQuad(r, c)

				masks := board[r].row & board[c].col & board[q].quad
				fmt.Printf("R %d C %d Q %d M %09b\n", r, c, q, masks)
				c++
			}
			r++
		}
}

func dumpMasks(board Board) {
	for i, data := range board {
		fmt.Printf("%2d %09b %09b %09b\n", i, data.row, data.col, data.quad)
	}
}
