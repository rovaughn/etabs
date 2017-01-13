package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var elasticTabstop = regexp.MustCompile(`( [[:space:]]|[\t\v\f\r])[[:space:]]*`)

type Cell struct {
	text     string
	trailing int
}

type Row []Cell
type Rows []Row

func (c Cell) Width() int {
	return len(c.text) + c.trailing
}

func (r Row) Width() int {
	total := 0
	for _, cell := range r {
		total += cell.Width()
	}
	return total
}

func (rows Rows) FixWidth(cols int, newWidth int) {
	for row := range rows {
		rows[row][cols-1].trailing += newWidth - rows[row][0:cols].Width()
	}
}

func FixTabstops(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	maxCols := 0
	table := make(Rows, 0)

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		trimmedLine := strings.TrimSpace(line)
		indent := line[0:strings.Index(line, trimmedLine)]
		pieces := elasticTabstop.Split(trimmedLine, -1)
		pieces[0] = indent + pieces[0]

		if len(pieces) > maxCols {
			maxCols = len(pieces)
		}

		row := make([]Cell, 0, len(pieces))

		for _, piece := range pieces {
			row = append(row, Cell{
				text:     piece,
				trailing: 0,
			})
		}

		table = append(table, row)
	}

	for cols := 1; cols <= maxCols; cols++ {
		blockStart := 0
		tabstop := 0

		for row := 0; row < len(table); row++ {
			if cols > len(table[row])-1 {
				table[blockStart:row].FixWidth(cols, tabstop)
				blockStart = row + 1
				tabstop = 0
			} else {
				width := table[row][0:cols].Width()

				if width > tabstop {
					tabstop = width
				}
			}
		}
	}

	for _, row := range table {
		for col, cell := range row {
			if _, err := w.Write([]byte(cell.text)); err != nil {
				return err
			}

			if col < len(row)-1 {
				if _, err := w.Write([]byte(strings.Repeat(" ", 2+cell.trailing))); err != nil {
					return err
				}
			}
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
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

		if err := FixTabstops(infile, outfile); err != nil {
			panic(err)
		}

		if err := os.Rename(outfile.Name(), filename); err != nil {
			panic(err)
		}
	}
}
