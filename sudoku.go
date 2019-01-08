package main

import (
	"fmt"
	"log"
	"time"
)

var count int = 0

type Masks struct {
	row uint16
	col uint16
	quad uint16
}

var board [MAX_ROWS]*Masks

func getQuad(row int, col int) int {
	return ((row/3) * 3) + (col/3)
}

func initBoard() {
	// Initialize all board masks to all values
	for i, _ := range board {
		board[i] = &Masks{
			row: MASK_ALL,
			col: MASK_ALL,
			quad: MASK_ALL,
		}
	}

	// Cycle through the raw data table and update the Masks
	for r, cols := range rootData {
		for c, val := range cols {
			if val > 0 {
				q := getQuad(r, c)
				mask, ok := Val2Mask[val]
				if !ok {
					log.Fatal("Invalid table")
				}

				// Remove the value from the list of possible values for the corresponding
				// row,, col, quad
				board[r].row = board[r].row &^ mask
				board[c].col = board[c].col &^ mask
				board[q].quad = board[q].quad &^ mask
			}
		}
	}
}

func solve() bool {
	for r, cols := range rootData {
		for c, val := range cols {
			if val == 0 {
				q := getQuad(r, c)

				// Find possible values
				masks := board[r].row & board[c].col & board[q].quad

				var b int8 = 1
				for b <= MAX_VAL {
					tmask, _ := Val2Mask[b]
					if masks & tmask == tmask {
						ormask := board[r].row
						ocmask := board[c].col
						oqmask := board[q].quad

						board[r].row = board[r].row &^ tmask
						board[c].col = board[c].col &^ tmask
						board[q].quad = board[q].quad &^ tmask

						rootData[r][c] = b

						if solve() {
							return true
						}

						board[r].row = ormask
						board[c].col = ocmask
						board[q].quad = oqmask
						rootData[r][c] = 0
					}

					b++
				}

				return false
			}
		}
	}

	return true
}

func main() {
	initBoard()

	start := time.Now()
	solved := solve()
	elapsed := time.Since(start)

	if !solved {
		dumpMasks()
	} else {
		fmt.Printf("Time to solve: %s\n\n", elapsed)
		dumpBoard()
	}
}

func dumpBoard() {
	for _, cols := range rootData {
		for _, val := range cols {
			fmt.Printf("\t%v ", val)
		}
		fmt.Println("")
	}
	fmt.Printf("\n\n")
}

func dumpMasks() {
	for _, data := range board {
		fmt.Printf("%09b %09b %09b\n", data.row, data.col, data.quad)
	}
}
