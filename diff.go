package main

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func lineDiff(src1, src2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(src1, src2)
	diffs := dmp.DiffMain(a, b, false)
	result := dmp.DiffCharsToLines(diffs, c)
	return result
}

func diffLines(src1, src2 string) []int {
	diffs := lineDiff(src1, src2)
	lineMap := make(map[int]bool)
	lines := []int{}
	lineCount := 0
	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			lineCount += strings.Count(d.Text, "\n")
		case diffmatchpatch.DiffDelete:
			// TODO: Handle Delete
		case diffmatchpatch.DiffInsert:
			count := strings.Count(d.Text, "\n")
			for i := 0; i < count; i++ {
				if !lineMap[lineCount+i] {
					lineMap[lineCount+i] = true
					lines = append(lines, lineCount+i)
				}
			}
			lineCount += count
		}
	}

	return lines
}
