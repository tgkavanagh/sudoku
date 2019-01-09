package main

import (
	"fmt"
)

func printSubStr(str string, low int, max int) {
	i := 0

	fmt.Printf("Longest palindrome substring is: ");

	for i < max {
		fmt.Printf("%c", str[low+i])
		i++
	}
}

func printTable(table [][]bool) {
	for _, entries := range table {
		fmt.Printf("%v\n", entries)
	}
}

func findLongest(table [][]bool, str string) int {
	start := 0
	max := 0

	for i, entries := range table {
		for j, pal := range entries {
			if pal {
				if (j - i + 1) > max {
					start = i
					max = j - i + 1
				}
			}
		}
	}

	printSubStr(str, start, max)

	return max
}

func longestpalSubstr(str string) int {
	n := len(str) // get length of input string

    // table[i][j] will be false if substring str[i..j]
    // is not palindrome.
    // Else table[i][j] will be true
    table := make([][]bool, n)
	i := 0
	for i < n {
		table[i] = make([]bool, n)
		i++
	}

	i = 0
	for i < n {
		j := 0
		for j < n {
			table[i][j] = false
			j++
		}
		i++
	}

	// All Subsstrings of 1 are palindromes
	i = 0
	for i < n  {
		table[i][i] = true
		i++
	}

    // check for sub-string of length 2.
	i = 0
	for i < n-1 {
		if str[i] == str[i+1] {
			table[i][i+1] = true
		}
		i++
	}

	printTable(table)

	// Check for lengths greater than 2. k is length
	// of substring
	k := 3
	for k <= n {
		// Fix the starting index
		i = 0
		for i < n-k+1 {
			// Get the ending index of substring from
			// starting index i and length k
			j := i + k - 1

			// checking for sub-string from ith index to
			// jth index iff str[i+1] to str[j-1] is a
			// palindrome
			if table[i+1][j-1] && str[i] == str[j] {
				table[i][j] = true
			}
			i++
		}

		k++
	}

//	fmt.Printf("Longest palindrome substring is: ");
//	printSubStr( str, start, start + maxLength - 1 );
	fmt.Println("\n\n")
	printTable(table)

	return findLongest(table, str) // return length of LPS
}

func main() {
	test := "forgeeksskeegfor"

	fmt.Printf("\nLength is: %d\n", longestpalSubstr(test))
}
