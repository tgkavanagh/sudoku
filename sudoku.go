package main

import (
	"fmt"
	"log"
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

type Sudoku struct {
	board [9][9]uint16
}

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

var emptyData = [9][9]int8{
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

var rootData = [9][9]int8{
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

var easy_rootData = [9][9]int8{
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

func getQuad(row int, col int) uint8 {
	return (uint8)(((row / 3) * 3) + (col / 3) + 1)
}

func initBoard() Sudoku {
	var sud Sudoku

	for r, cols := range rootData {
		for c, val := range cols {
			if mask, ok := Val2Mask[val]; !ok {
				log.Fatal("Invalid value in table (%d, %d): %v\n", r, c, val)
			} else {
				sud.board[r][c] = mask
			}
		}
	}

	return sud
}

func (sud Sudoku) dumpBoard(convert bool) {
	for _, cols := range sud.board {
		for _, mask := range cols {
			if convert {
				val, _ := Mask2Val[mask]
				fmt.Printf("\t%v ", val)
			} else {
				fmt.Printf("\t%09b ", mask)
			}
		}
		fmt.Println("")
	}
	fmt.Printf("\n\n")
}

func (sud *Sudoku) checkBoard(row int, col int, mask uint16) bool {
	quad := getQuad(row, col)

	for r, cols := range sud.board {
		if r != row {
			// Not our row so we need to check the specified column or any
			// column with a matching quad
			for c, tmask := range cols {
				q := getQuad(r, c)

				if c == col || q == quad {
					if tmask == mask {
						return false
					}
				}
			}
		} else {
			// This is our row so check all columns except our
			for c, tmask := range cols {
				if c != col && tmask == mask {
					return false
				}
			}
		}
	}

	return sud.solve()
}

func (sud *Sudoku) updateBoard() {
	for row, cols := range sud.board {
		for col, mask := range cols {
			if _, ok := Mask2Val[mask]; ok {
				quad := getQuad(row, col)

				// Loop through the table again and mark off possible values
				for r, ucols := range sud.board {
					if r != row {
						// Update all entries with a match column or quad
						for c, tmask := range ucols {
							q := getQuad(r, c)
							if c == col || q == quad && tmask != mask {
								sud.board[r][c] = tmask &^ mask
							}
						}
					} else {
						// Update all cells in our row
						for c, tmask := range ucols {
							if tmask != mask {
								sud.board[r][c] = tmask &^ mask
							}
						}
					}
				}
			}
		}
	}
}

func (sud *Sudoku) solve() bool {
	for r, cols := range sud.board {
		for c, mask := range cols {
			if _, ok := Mask2Val[mask]; !ok {
				var b int8 = 1
				for b <= MAX_VAL {
					bmask, _ := Val2Mask[b]
					if mask&bmask == bmask {
						sud.board[r][c] = bmask
						if sud.checkBoard(r, c, sud.board[r][c]) {
							return true
						}

						sud.board[r][c] = mask
					}
					b++
				}
				return false
			}
		}
	}

	return true
}

func (sud Sudoku) Solve(update bool) {
	start := time.Now()

	if update {
		sud.updateBoard()
	}

	solved := sud.solve()
	elapsed := time.Since(start)

	if !solved {
		sud.dumpBoard(false)
	} else {
		fmt.Printf("Time to solve: %s\n\n", elapsed)
		sud.dumpBoard(true)
	}
}

func main() {
	sud := initBoard()
	sud.Solve(false)

	sud1 := initBoard()
	sud1.Solve(true)
}
