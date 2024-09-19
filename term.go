package virtualterm

import "strings"

// VirtualTerm this is created to simulate a terminal,handle the special character such as '\r','\b', "\033[1D".
// For example: if you input "cute\rhat", the result of String() would be "hate"
type VirtualTerm struct {
	// the column of cursor
	curCol int
	// the line of cursor
	curLine int
	// the content of virtual terminal
	content [][]rune
}

type Option func(*VirtualTerm)

// NewOptions create a virtual terminal
func NewOptions(opts ...Option) VirtualTerm {
	vt := VirtualTerm{
		content: [][]rune{},
	}
	for _, o := range opts {
		o(&vt)
	}
	return vt
}

// NewDefault create a default virtual terminal
func NewDefault() VirtualTerm {
	return NewOptions()
}

// Write implements the io.Writer interface
func (vt *VirtualTerm) Write(p []byte) (n int, err error) {
	if lineCnt := len(vt.content); lineCnt <= vt.curLine {
		vt.content = append(vt.content, make([][]rune, vt.curLine+1-lineCnt)...)
	}

	for i := 0; i < len(p); i++ {
		switch p[i] {
		case '\r':
			// Carriage Return
			vt.curCol = 0
		case '\n':
			// NewLine
			// If the cursor is on the last line, add a new line
			vt.curLine++
			vt.curCol = 0
			for vt.curLine >= len(vt.content) {
				vt.content = append(vt.content, nil)
			}
		case '\b':
			if vt.curCol == 0 {
				continue
			} else {
				vt.curCol--
			}
		case '\033':
			// TODO complete this

		default:
			vt.writeRune(rune(p[i]))
		}

	}
	return len(p), nil
}

// writeRune write Rune to content.
func (vt *VirtualTerm) writeRune(r rune) {
	// check the length of line,if it is not enough, make it longer
	if lineLen := len(vt.content[vt.curLine]); lineLen <= vt.curCol {
		vt.content[vt.curLine] = append(vt.content[vt.curLine], []rune(strings.Repeat(" ", vt.curCol+1-lineLen))...)
	}
	vt.content[vt.curLine][vt.curCol] = r
	vt.curCol++
}

// writeRunes write Rune array to content
func (vt *VirtualTerm) writeRunes(rs []rune) {
	for _, r := range rs {
		vt.writeRune(r)
	}
}

// writeString write String to content
func (vt *VirtualTerm) writeString(s string) {
	for _, c := range s {
		vt.writeRune(c)
	}
}

func (vt *VirtualTerm) String() string {
	builder := strings.Builder{}
	for i, line := range vt.content {
		for _, c := range line {
			builder.WriteRune(c)
		}
		if i != len(vt.content)-1 {
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

// Clear all the content in virtual terminal
func (vt *VirtualTerm) Clear() {
	vt.curCol = 0
	vt.curLine = 0
	vt.content = nil
}
