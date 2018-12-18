package main

import (
	"fmt"
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

type CellData struct {
	row uint8
	col uint8
	quad uint8
	mask uint16
}

type Sudoku struct {
	board []*CellData
}

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

var rootData = [9][9]int8 {
	{6, 4, 0, 0, 3, 0, 0, 0, 7}, //1
	{5, 0, 1, 0, 7, 0, 9, 0, 0}, //2
	{0, 0, 0, 0, 0, 0, 0, 1, 0}, //3

	{0, 0, 4, 9, 0, 8, 0, 6, 0}, //4
	{0, 8, 0, 0, 0, 3, 0, 2, 0}, //5
	{0, 0, 0, 4, 0, 0, 0, 0, 0}, //6

	{4, 0, 0, 1, 5, 7, 0, 3, 0}, //7
	{2, 0, 8, 3, 0, 0, 0, 4, 0}, //8
	{7, 5, 0, 0, 0, 0, 0, 9, 6}, //9
}

func getQuad(row uint8, col uint8) uint8 {
	return (uint8)(((row/3) * 3) + (col/3) + 1)
}

func initBoard() Sudoku {
	var sud Sudoku
	var row uint8 = 0

	for row <  MAX_ROWS {
		var col uint8 = 0

		for col < MAX_COLS {
			quad := getQuad(row, col)

			mask, ok := Val2Mask[rootData[row][col]]
			if !ok {
				fmt.Printf("Invalid value: %d\n")
				mask = 511
			}

			newCell := &CellData {
				row: row,
				col: col,
				quad: quad,
				mask: mask,
			}

			sud.board = append(sud.board, newCell)
			col++
		}
		row++
	}

	return sud
}

func (sud Sudoku) dumpBoard(convert bool) {
	var prevRow uint8 = 0
	for _, data := range sud.board {
		if data.row != prevRow {
			prevRow = data.row
			fmt.Println("")
		}

		val, ok := Mask2Val[data.mask]
		if !ok {
			val = 0
		}

		if convert {
			fmt.Printf("\t%v ", val)
		} else {
			fmt.Printf("\t%09b ", data.mask)
		}
	}
	fmt.Printf("\n\n")
}

func (sud Sudoku) checkBoard(idx int) bool {
	cell := sud.board[idx]

	_, ok := Mask2Val[cell.mask]
	if !ok {
		return false
	}

	// Ensure the value is not already used in the corresponding row, col or getQuad
	for _, tmp := range sud.board {
		if tmp.row == cell.row && tmp.col == cell.col {
			continue
		}

		if tmp.row == cell.row || tmp.col == cell.col || tmp.quad == cell.quad {
			if tmp.mask == cell.mask {
				return false
			}
		}
	}

	return sud.solve()
}

func (sud Sudoku)solve() bool {
	// Find a cell that does not have a set value
	for i, cell := range sud.board {
		_, ok := Mask2Val[cell.mask]
		if !ok {
			// We have a winner
			var b int8 = 1
			for b <= MAX_VAL {
				if mask, ok := Val2Mask[b]; ok {
					oMask := cell.mask
					cell.mask = cell.mask & mask
					if _, ok := Mask2Val[cell.mask]; ok {
						if sud.checkBoard(i) {
							return true
						}
					}
					cell.mask = oMask
				}
				b++
			}
			return false
		}
	}

	return true
}

func main() {
	start := time.Now()

	sud := initBoard()

	solved := sud.solve()

	elapsed := time.Since(start)

	fmt.Printf("Time to solve: %s\n\n", elapsed)
	if !solved {
		sud.dumpBoard(false)
	} else {
		sud.dumpBoard(true)
	}

}