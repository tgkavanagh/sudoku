package main

import (
	"fmt"
	"sort"
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
	pval int
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

var emptyData = [9][9]int8 {
//   0  1  2  3  4  5  6  7  8
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //0
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //1
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //2
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //3
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //4
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //5
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //6
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //7
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //8
}

var rootData = [9][9]int8 {
//   0  1  2  3  4  5  6  7  8
	{0, 0, 0, 5, 0, 0, 2, 0, 0}, //0
	{0, 8, 0, 0, 0, 0, 0, 0, 0}, //1
	{0, 0, 0, 1, 0, 0, 0, 0, 0}, //2
	{0, 0, 0, 0, 7, 2, 0, 0, 3}, //3
	{5, 0, 1, 0, 0, 0, 0, 4, 0}, //4
	{6, 0, 0, 0, 0, 0, 0, 0, 0}, //5
	{0, 0, 0, 0, 0, 7, 0, 5, 0}, //6
	{0, 2, 0, 0, 3, 0, 0, 0, 0}, //7
	{4, 0, 0, 0, 0, 0, 1, 0, 0}, //8
}

var expert_rootData = [9][9]int8 {
//   0  1  2  3  4  5  6  7  8
	{0, 7, 2, 0, 3, 0, 0, 8, 9}, //0
	{5, 0, 0, 8, 0, 0, 0, 0, 0}, //1
	{8, 0, 0, 0, 4, 0, 0, 0, 0}, //2
	{0, 0, 0, 9, 0, 2, 8, 7, 0}, //3
	{0, 0, 3, 0, 0, 5, 0, 2, 0}, //4
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, //5
	{3, 0, 0, 0, 1, 0, 5, 0, 0}, //6
	{0, 0, 7, 0, 0, 0, 0, 0, 0}, //7
	{0, 0, 0, 0, 0, 0, 6, 9, 1}, //8
}

var easy_rootData = [9][9]int8 {
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

type ByPval []*CellData
func (pv ByPval) Len() int { return len(pv) }
func (pv ByPval) Less(i, j int) bool { return pv[i].pval < pv[j].pval }
func (pv ByPval) Swap(i, j int) { pv[i], pv[j] = pv[j], pv[i] }

type ByPos []*CellData
func (pv ByPos) Len() int { return len(pv) }
func (pv ByPos) Swap(i, j int) { pv[i], pv[j] = pv[j], pv[i] }
func (pv ByPos) Less(i, j int) bool {
	if pv[i].row == pv[j].row {
		return pv[i].col < pv[j].col
	}

	return pv[i].row < pv[j].row
}

var count int = 0

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

			pval := 1
			mask, ok := Val2Mask[rootData[row][col]]
			if !ok {
				fmt.Printf("Invalid value: %d\n")
				mask = 511
			}

			if mask == 511 {
				pval = 9
			}

			newCell := &CellData {
				row: row,
				col: col,
				quad: quad,
				mask: mask,
				pval: pval,
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
				count++
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

func (sud Sudoku)update() {
	for _, cell := range sud.board {
		if _, ok := Mask2Val[cell.mask]; ok {
			// Loop through the entire board and mark off the correct cells
			for _, tmp := range sud.board {
				if cell.row == tmp.row && cell.col == tmp.col {
					continue
				}

				if _, ok := Mask2Val[tmp.mask]; !ok {
					if cell.row == tmp.row || cell.col == tmp.col || cell.quad == tmp.quad {
						tmask := tmp.mask &^ cell.mask
						if tmask != tmp.mask {
							tmp.mask = tmask
							tmp.pval--
						}
					}
				}
			}
		}
	}
}

func (sud Sudoku)Solve(update bool, sortboard bool) {
	//Start brute force
	start := time.Now()

	if update {
		sud.update()
	}

	if sortboard {
		sort.Sort(ByPval(sud.board))
/*
		for _, data := range sud.board {
			fmt.Printf("%v\n", data)
		}
		return
*/
	}

	solved := sud.solve()
	elapsed := time.Since(start)

	if sortboard {
		sort.Sort(ByPos(sud.board))
	}

	if !solved {
		sud.dumpBoard(false)
	} else {
		fmt.Printf("Time to solve: %s\n\n", elapsed)
		sud.dumpBoard(true)
	}
}

func main() {
	/*Board()
	sud.Solve(false, false)
	fmt.Printf("Count: %d\n", count)
	count = 0

	sud1 := initBoard()
	sud1.Solve(true, false)
	fmt.Printf("Count: %d\n", count)
	count = 0
*/
	sud2 := initBoard()
	sud2.Solve(true, true)
	fmt.Printf("Count: %d\n", count)
}
