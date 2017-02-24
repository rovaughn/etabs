package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

var columnRe = regexp.MustCompile(`((^[ \t]*)?([^ \t].*?))(( [\t ]|\t)[ \t]*|$)`)

type Column struct {
	text     string
	trailing int
}

func MeasureWidth(cols []Column) int {
	total := 0
	for i, col := range cols {
		total += utf8.RuneCountInString(col.text)
		if i < len(cols)-1 {
			total += col.trailing
		}
	}
	return total
}

var Unchanged = errors.New("No changes made")

func FindBlocks(rows [][]Column, numCols int) [][][]Column {
	blockStart := 0
	blockEnd := 0
	blocks := make([][][]Column, 0)

	for row := 0; row <= len(rows); row++ {
		if row < len(rows) && len(rows[row]) >= numCols {
			blockEnd = row + 1
		} else {
			if blockEnd > blockStart {
				blocks = append(blocks, rows[blockStart:blockEnd])
			}

			blockStart = row + 1
			blockEnd = row + 1
		}
	}

	return blocks
}

// This function works by converting the text into a "tabular" format, where
// the text is multiple rows, each containing multiple columns.  A given cell
// stores its text and its trailing length.
func FixTabstops(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	maxCols := 0
	table := make([][]Column, 0)

	changed := false
	for scanner.Scan() {
		line := scanner.Text()
		columns := columnRe.FindAllStringSubmatch(line, -1)

		row := make([]Column, len(columns))
		for i, col := range columns {
			row[i].text = col[1]

			if strings.Contains(col[4], "\t") {
				changed = true
				row[i].trailing = 0
			} else {
				row[i].trailing = len(col[4])
			}
		}

		if len(columns) > maxCols {
			maxCols = len(columns)
		}

		table = append(table, row)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for alignCol := 1; alignCol < maxCols; alignCol++ {
		for _, block := range FindBlocks(table, alignCol+1) {
			widest := 0
			for _, row := range block {
				width := MeasureWidth(row[:alignCol])

				if width > widest {
					widest = width
				}
			}

			for _, row := range block {
				width := MeasureWidth(row[:alignCol]) + row[alignCol-1].trailing
				delta := widest + 2 - width
				changed = changed || delta != 0
				row[alignCol-1].trailing += delta
			}
		}
	}

	if !changed {
		return Unchanged
	}

	bw := bufio.NewWriter(w)
	for _, row := range table {
		for col, cell := range row {
			if _, err := bw.WriteString(cell.text); err != nil {
				return err
			}

			if col < len(row)-1 {
				if _, err := bw.Write(bytes.Repeat([]byte(" "), cell.trailing)); err != nil {
					return err
				}
			}
		}
		if err := bw.WriteByte('\n'); err != nil {
			return err
		}
	}

	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: etabs [FILE]")
		fmt.Fprintln(os.Stderr, "    or etabs - to use stdin/stdout")
		return
	}

	filename := os.Args[1]

	if filename == "-" {
		if err := FixTabstops(os.Stdin, os.Stdout); err != nil {
			panic(err)
		}
	} else {
		infile, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer infile.Close()

		outfile, err := ioutil.TempFile(".", "etab")
		if err != nil {
			panic(err)
		}
		defer outfile.Close()
		defer os.Remove(outfile.Name())

		if err := FixTabstops(infile, outfile); err == Unchanged {
			return
		} else if err != nil {
			panic(err)
		}

		if err := os.Rename(outfile.Name(), filename); err != nil {
			panic(err)
		}
	}
}
