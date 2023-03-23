//
// Copyright (c) 2021-2023 Markku Rossi
//
// All rights reserved.
//

package asciiart

import (
	"strings"
)

type Region struct {
	maxWidth int
	lines    [][]rune
}

func (r *Region) String() string {
	var b strings.Builder

	for row, line := range r.lines {
		if row > 0 {
			b.WriteRune('\n')
		}
		for _, r := range line {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (r *Region) Width() int {
	return r.maxWidth
}

func (r *Region) Height() int {
	return len(r.lines)
}

func (r *Region) Get(row, col int) rune {
	if row < 0 || row >= len(r.lines) {
		return 0
	}
	if col < 0 || col >= len(r.lines[row]) {
		return 0
	}
	return r.lines[row][col]
}

func (r *Region) Set(row, col int, ch rune) {
	if row < 0 || row >= len(r.lines) {
		return
	}
	if col < 0 || col >= len(r.lines[row]) {
		return
	}
	r.lines[row][col] = ch
}

func (r *Region) Clone() *Region {
	lines := make([][]rune, len(r.lines))
	for i, line := range r.lines {
		l := make([]rune, len(line))
		copy(l, line)
		lines[i] = l
	}
	return &Region{
		maxWidth: r.maxWidth,
		lines:    lines,
	}
}

func NewRegion(input []byte) *Region {
	var lines [][]rune
	var width int

	for _, line := range strings.Split(string(input), "\n") {
		runes := []rune(line)
		if len(runes) > width {
			width = len(runes)
		}
		lines = append(lines, runes)
	}
	return &Region{
		maxWidth: width,
		lines:    lines,
	}
}

const (
	FlagUp int = 1 << iota
	FlagDown
	FlagLeft
	FlagRight
)

var properties = map[rune]int{
	'+': FlagUp | FlagDown | FlagLeft | FlagRight,
	'*': FlagUp | FlagDown | FlagLeft | FlagRight,
	'|': FlagUp | FlagDown,
	'-': FlagLeft | FlagRight,

	0x2500: FlagLeft | FlagRight,
	0x2501: FlagLeft | FlagRight,
	0x2502: FlagDown | FlagUp,
	0x2503: FlagDown | FlagUp,
	0x2504: FlagLeft | FlagRight,
	0x2505: FlagLeft | FlagRight,
	0x2506: FlagDown | FlagUp,
	0x2507: FlagDown | FlagUp,
	0x2508: FlagLeft | FlagRight,
	0x2509: FlagLeft | FlagRight,
	0x250A: FlagDown | FlagUp,
	0x250B: FlagDown | FlagUp,
	0x250C: FlagDown | FlagRight,
	0x250D: FlagDown | FlagRight,
	0x250E: FlagDown | FlagRight,
	0x250F: FlagDown | FlagRight,

	0x2510: FlagDown | FlagLeft,
	0x2511: FlagDown | FlagLeft,
	0x2512: FlagDown | FlagLeft,
	0x2513: FlagDown | FlagLeft,
	0x2514: FlagUp | FlagRight,
	0x2515: FlagUp | FlagRight,
	0x2516: FlagUp | FlagRight,
	0x2517: FlagUp | FlagRight,
	0x2518: FlagUp | FlagLeft,
	0x2519: FlagUp | FlagLeft,
	0x251A: FlagUp | FlagLeft,
	0x251B: FlagUp | FlagLeft,
	0x251C: FlagUp | FlagDown | FlagRight,
	0x251D: FlagUp | FlagDown | FlagRight,
	0x251E: FlagUp | FlagDown | FlagRight,
	0x251F: FlagUp | FlagDown | FlagRight,

	0x2520: FlagUp | FlagDown | FlagRight,
	0x2521: FlagUp | FlagDown | FlagRight,
	0x2522: FlagUp | FlagDown | FlagRight,
	0x2523: FlagUp | FlagDown | FlagRight,
	0x2524: FlagUp | FlagDown | FlagLeft,
	0x2525: FlagUp | FlagDown | FlagLeft,
	0x2526: FlagUp | FlagDown | FlagLeft,
	0x2527: FlagUp | FlagDown | FlagLeft,
	0x2528: FlagUp | FlagDown | FlagLeft,
	0x2529: FlagUp | FlagDown | FlagLeft,
	0x252A: FlagUp | FlagDown | FlagLeft,
	0x252B: FlagUp | FlagDown | FlagLeft,
	0x252C: FlagDown | FlagLeft | FlagRight,
	0x252D: FlagDown | FlagLeft | FlagRight,
	0x252E: FlagDown | FlagLeft | FlagRight,
	0x252F: FlagDown | FlagLeft | FlagRight,

	0x2530: FlagDown | FlagLeft | FlagRight,
	0x2531: FlagDown | FlagLeft | FlagRight,
	0x2532: FlagDown | FlagLeft | FlagRight,
	0x2533: FlagDown | FlagLeft | FlagRight,
	0x2534: FlagUp | FlagLeft | FlagRight,
	0x2535: FlagUp | FlagLeft | FlagRight,
	0x2536: FlagUp | FlagLeft | FlagRight,
	0x2537: FlagUp | FlagLeft | FlagRight,
	0x2538: FlagUp | FlagLeft | FlagRight,
	0x2539: FlagUp | FlagLeft | FlagRight,
	0x253A: FlagUp | FlagLeft | FlagRight,
	0x253B: FlagUp | FlagLeft | FlagRight,
	0x253C: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x253D: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x253E: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x253F: FlagUp | FlagDown | FlagLeft | FlagRight,

	0x2540: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2541: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2542: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2543: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2544: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2545: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2546: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2547: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2548: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x2549: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x254A: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x254B: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x254C: FlagLeft | FlagRight,
	0x254D: FlagLeft | FlagRight,
	0x254E: FlagUp | FlagDown,
	0x254F: FlagUp | FlagDown,

	0x2550: FlagLeft | FlagRight,
	0x2551: FlagUp | FlagDown,
	0x2552: FlagDown | FlagRight,
	0x2553: FlagDown | FlagRight,
	0x2554: FlagDown | FlagRight,
	0x2555: FlagDown | FlagLeft,
	0x2556: FlagDown | FlagLeft,
	0x2557: FlagDown | FlagLeft,
	0x2558: FlagUp | FlagRight,
	0x2559: FlagUp | FlagRight,
	0x255A: FlagUp | FlagRight,
	0x255B: FlagUp | FlagLeft,
	0x255C: FlagUp | FlagLeft,
	0x255D: FlagUp | FlagLeft,
	0x255E: FlagUp | FlagDown | FlagRight,
	0x255F: FlagUp | FlagDown | FlagRight,

	0x2560: FlagUp | FlagDown | FlagRight,
	0x2561: FlagUp | FlagDown | FlagLeft,
	0x2562: FlagUp | FlagDown | FlagLeft,
	0x2563: FlagUp | FlagDown | FlagLeft,
	0x2564: FlagDown | FlagLeft | FlagRight,
	0x2565: FlagDown | FlagLeft | FlagRight,
	0x2566: FlagDown | FlagLeft | FlagRight,
	0x2567: FlagUp | FlagLeft | FlagRight,
	0x2568: FlagUp | FlagLeft | FlagRight,
	0x2569: FlagUp | FlagLeft | FlagRight,
	0x256A: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x256B: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x256C: FlagUp | FlagDown | FlagLeft | FlagRight,
	0x256D: FlagDown | FlagRight,
	0x256E: FlagDown | FlagLeft,
	0x256F: FlagUp | FlagLeft,

	0x2570: FlagUp | FlagRight,
	0x2571: 0,
	0x2572: 0,
	0x2573: 0,
	0x2574: FlagLeft,
	0x2575: FlagUp,
	0x2576: FlagRight,
	0x2577: FlagDown,
	0x2578: FlagLeft,
	0x2579: FlagUp,
	0x257A: FlagRight,
	0x257B: FlagDown,
	0x257C: FlagLeft | FlagRight,
	0x257D: FlagUp | FlagDown,
	0x257E: FlagLeft | FlagRight,
	0x257F: FlagUp | FlagDown,
}

var lineDrawing = []rune{
	'+',    //
	0x2575, // 	               Up	╵
	0x2577, // 	          Down		╷
	0x2502, // 	          Down Up	│
	0x2574, // 	     Left			╴
	0x2518, // 	     Left      Up	┘
	0x2510, // 	     Left Down		┐
	0x2524, // 	     Left Down Up	┤
	0x2576, // Right				╶
	0x2514, // Right           Up	└
	0x250C, // Right      Down		┌
	0x251C, // Right      Down Up	├
	0x2500, // Right Left			─
	0x2534, // Right Left      Up	┴
	0x252C, // Right Left Down		┬
	0x253C, // Right Left Down Up	┼
}

var lineDrawingRound = []rune{
	'*',    //
	0x2575, // 	               Up	╵
	0x2577, // 	          Down	 	╷
	0x2502, // 	          Down Up	│
	0x2574, // 	     Left		 	╴
	0x256F, // 	     Left      Up	╯
	0x256E, // 	     Left Down	 	╮
	0x2524, // 	     Left Down Up	┤
	0x2576, // Right			 	╶
	0x2570, // Right           Up	╰
	0x256D, // Right      Down	 	╭
	0x251C, // Right      Down Up	├
	0x2500, // Right Left		 	─
	0x2534, // Right Left      Up	┴
	0x252C, // Right Left Down	 	┬
	0x253C, // Right Left Down Up	┼
}

var backslash = []rune{
	'\\',   //
	0x2570, // 	               Up	╰
	0x256E, // 	          Down		╮
	'\\',   // 	          Down Up
	0x256E, // 	     Left			╮
	0x2518, // 	     Left      Up	┘
	0x256E, // 	     Left Down		╮
	0x2524, // 	     Left Down Up	┤
	0x2570, // Right				╰
	0x2570, // Right           Up	╰
	0x250C, // Right      Down		┌
	0x251C, // Right      Down Up	├
	'\\',   // Right Left
	0x2534, // Right Left      Up	┴
	0x252C, // Right Left Down		┬
	0x253C, // Right Left Down Up	┼
}

var slash = []rune{
	'/',    //
	0x256F, // 	               Up	╯
	0x256D, // 	          Down		╭
	'/',    // 	          Down Up
	0x256F, // 	     Left			╯
	0x256F, // 	     Left      Up	╯
	0x2510, // 	     Left Down		┐
	0x2524, // 	     Left Down Up	┤
	0x256D, // Right				╭
	0x2514, // Right           Up	└
	0x256D, // Right      Down		╭
	0x251C, // Right      Down Up	├
	'/',    // Right Left
	0x2534, // Right Left      Up	┴
	0x252C, // Right Left Down		┬
	0x253C, // Right Left Down Up	┼
}

func isLine(r rune, f int) bool {
	props, ok := properties[r]
	if !ok {
		return false
	}
	return props&f != 0
}

func checkMainLines(input *Region, row, col int) int {
	var index int
	if input.Get(row-1, col) == '|' {
		index |= FlagUp
	}
	if input.Get(row+1, col) == '|' {
		index |= FlagDown
	}
	if input.Get(row, col-1) == '-' {
		index |= FlagLeft
	}
	if input.Get(row, col+1) == '-' {
		index |= FlagRight
	}
	return index
}

func Process(data string) string {
	input := NewRegion([]byte(data))
	output := input.Clone()

	for row := 0; row < input.Height(); row++ {
		for col := 0; col < input.Width(); col++ {
			ch := input.Get(row, col)
			switch ch {
			case '|':
				output.Set(row, col, 0x2502)
			case '-':
				output.Set(row, col, 0x2500)
			case '+', '*', '\'', '\\', '/':
				var index int
				if isLine(input.Get(row-1, col), FlagDown) {
					index |= FlagUp
				}
				if isLine(input.Get(row+1, col), FlagUp) {
					index |= FlagDown
				}
				if isLine(input.Get(row, col-1), FlagRight) {
					index |= FlagLeft
				}
				if isLine(input.Get(row, col+1), FlagLeft) {
					index |= FlagRight
				}
				if index == 0 {
					output.Set(row, col, ch)
				} else {
					switch ch {
					case '+', '\'':
						output.Set(row, col, lineDrawing[index])

					case '*':
						output.Set(row, col, lineDrawingRound[index])

					case '\\':
						mainLines := checkMainLines(input, row, col)
						if mainLines != 0 {
							output.Set(row, col, backslash[mainLines])
						} else {
							output.Set(row, col, backslash[index])
						}

					case '/':
						mainLines := checkMainLines(input, row, col)
						if mainLines != 0 {
							output.Set(row, col, slash[mainLines])
						} else {
							output.Set(row, col, slash[index])
						}
					}
				}
			}
		}
	}

	return output.String()
}
