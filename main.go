package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	flagWordSize = flag.Int("size", 2, "size")
)

const (
	contextType = iota
	oldLineType
	newLineType
)

type Worker struct {
	size int
	r    io.Reader
	w    io.Writer
}

func main() {
	flag.Parse()
	if *flagWordSize < 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	worker := Worker{
		size: *flagWordSize,
		r:    os.Stdin,
		w:    os.Stdout,
	}

	worker.Do()
}

type Block struct {
	currentType int
	switches    int
	oldLines    []string
	newLines    []string
}

func NewBlock() Block {
	return Block{
		currentType: contextType,
	}
}

func (obj Worker) Do() {
	sc := bufio.NewScanner(obj.r)
	block := NewBlock()
	for sc.Scan() {
		line := sc.Text()
		lineType := contextType
		if strings.HasPrefix(line, "+") {
			lineType = newLineType
			block.newLines = append(block.newLines, strings.TrimPrefix(line, "+"))
		} else if strings.HasPrefix(line, "-") {
			lineType = oldLineType
			block.oldLines = append(block.oldLines, strings.TrimPrefix(line, "-"))
		}

		if block.currentType == lineType {
			continue
		}

		// fmt.Printf("line: %q, Block: %+v\n", line, block)
		if lineType == contextType || block.switches > 0 {
			obj.DoBlock(block)
			block = NewBlock()
			continue
		}

		if block.currentType != contextType {
			block.switches += 1
		}
		block.currentType = lineType
	}
}

func (obj Worker) DoBlock(block Block) {
	if len(block.newLines) == 0 || len(block.oldLines) == 0 {
		obj.Print(block.oldLines, block.newLines)
	}

	oldL := []string{}
	newL := []string{}

	oldi := len(block.oldLines)
	newi := len(block.newLines)
	maxi := oldi
	if maxi < newi {
		maxi = newi
	}

	for i := 0; i < maxi; i++ {
		if i > oldi-1 {
			newL = append(newL, block.newLines[i])
			continue
		}
		if i > newi-1 {
			oldL = append(oldL, block.oldLines[i])
			continue
		}

		// equal lines. skipping
		if obj.FormatLine(block.newLines[i]) == obj.FormatLine(block.oldLines[i]) {
			continue
		}

		oldL = append(oldL, block.oldLines[i])
		newL = append(newL, block.newLines[i])
	}
	obj.Print(oldL, newL)
}

func (obj Worker) FormatLine(line string) string {
	res := ""
	cur := ""
	curIsNumber := true
	curIsWord := false
	for _, r := range line {
		if isAlphaNumeric(r) {
			if curIsNumber && isLetter(r) {
				curIsNumber = false
			}
			cur += string(r)
			curIsWord = true
			continue
		}

		if curIsWord {
			if curIsNumber && len(cur) <= obj.size*2 {
				cur = "-"
			} else if len(cur) <= obj.size {
				cur = "_"
			}
		} else {
			cur += string(r)
		}

		res += cur
		cur = ""
		curIsNumber = true
		curIsWord = false
	}

	return res
}

func isAlphaNumeric(ch rune) bool {
	return isDecimal(ch) || isLetter(ch)
}

func isLetter(ch rune) bool {
	return isLowerLetter(ch) || isUpperLetter(ch)
}

func isLowerLetter(ch rune) bool {
	return ch >= 'a' && ch <= 'z'
}

func isUpperLetter(ch rune) bool {
	return ch >= 'A' && ch <= 'Z'
}

func isDecimal(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (obj Worker) Print(oldL, newL []string) {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	for _, v := range oldL {
		red.Fprintf(obj.w, "-%s\n", v)
	}
	for _, v := range newL {
		green.Fprintf(obj.w, "+%s\n", v)
	}
}
