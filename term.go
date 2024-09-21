package virtualterm

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

var INF = math.MaxInt - 1
var ErrCannotHandle = errors.New("this csi is not supported or syntax error")
var ErrNonDetermistics = errors.New("non-deterministic")

// VirtualTerm this is created to simulate a terminal,handle the special character such as '\r','\b', "\033[1D".
// For example: if you input "cute\rhat", the result of String() would be "hate"
type VirtualTerm struct {
	// the column of cursor
	cx int
	// the line of cursor
	cy int
	// the content of virtual terminal
	content [][]rune

	// xOffset is the offset on X.
	xOffset int

	// silence will shut down the log.By default, it is true
	silence bool
}

type Option func(*VirtualTerm)

// NewOptions create a virtual terminal
func NewOptions(opts ...Option) VirtualTerm {
	vt := VirtualTerm{
		content: emptyContent(),
		cx:      0,
		cy:      0,
		silence: true,
	}

	for _, o := range opts {
		o(&vt)
	}
	return vt
}

// log is used to output something
func (vt *VirtualTerm) log(str string) {
	if vt.silence {
		return
	}
	log.Println(str)
}

// NewDefault create a default virtual terminal
func NewDefault() VirtualTerm {
	return NewOptions()
}

// OptionSilence set whether the vt is silence or not
func OptionSilence(val bool) Option {
	return func(vt *VirtualTerm) {
		vt.silence = val
	}
}

// handleCSI handle control sequence introducer
// implement according to https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ECMA-48-5
func (vt *VirtualTerm) handleCSI(csi string) error {
	// an example \033[A
	// csi should not be shorter than 3
	if len(csi) < 3 {
		return fmt.Errorf("too short, %w", ErrCannotHandle)
	}
	param := csi[2 : len(csi)-1]
	// final byte is a single character
	final := csi[len(csi)-1]
	var err error
	switch final {
	case 'A':
		// move the cursor up
		var cnt int
		if len(param) == 0 {
			cnt = 1
		} else {
			cnt, err = strconv.Atoi(param)
		}
		if err != nil {
			return fmt.Errorf("%w:ESC[#A param may be wrong", err)
		}
		vt.cursorMove(0, -cnt)
	case 'B':
		// move the cursor down
		var cnt int
		if len(param) == 0 {
			cnt = 1
		} else {
			cnt, err = strconv.Atoi(param)
		}
		if err != nil {
			return fmt.Errorf("%w:ESC[#B param may be wrong", err)
		}
		vt.cursorMove(0, cnt)
	case 'C':
		// move the cursor right
		var cnt int
		if len(param) == 0 {
			cnt = 1
		} else {
			cnt, err = strconv.Atoi(param)
		}
		if err != nil {
			return fmt.Errorf("%w:ESC[#C param may be wrong", err)
		}
		vt.cursorMove(cnt, 0)
	case 'D':
		// move the cursor left
		var cnt int
		if len(param) == 0 {
			cnt = 1
		} else {
			cnt, err = strconv.Atoi(param)
		}
		if err != nil {
			return fmt.Errorf("%w:ESC[#D param may be wrong", err)
		}
		vt.cursorMove(-cnt, 0)
	case 'H':
		// move the cursor to the home position
		vt.cursorHome()
	default:
		return fmt.Errorf("%q%w", csi, ErrCannotHandle)
	}
	return nil
}

// isIntermediateByte check whether it is CSI param byte
func isIntermediateRune(c rune) bool {
	return c >= 0x20 && c <= 0x2F
}

// isParameterByte check whether it is CSI parameter byte
func isParameterRune(c rune) bool {
	return c >= 0x30 && c <= 0x3F
}

// isCSIFinalByte check whether it is CSI final byte
func isCSIFinalRune(c rune) bool {
	return c >= 0x40 && c <= 0x7E
}

// Write implements the io.Writer interface
func (vt *VirtualTerm) Write(b []byte) (n int, err error) {
	return vt.WriteRunes([]rune(string(b)))
}

// WriteString write string to virtual terminal
func (vt *VirtualTerm) WriteString(s string) (n int, err error) {
	return vt.Write([]byte(s))
}

// runeWidth
func (*VirtualTerm) runeWidth(r rune) int {
	if len(string(r)) >= 3 {
		return 2
	} else {
		return 1
	}
}

// cursorMove can control the cursor
func (vt *VirtualTerm) cursorMove(x int, y int) {
	// handle the offset first
	if vt.xOffset != 0 {
		x += vt.xOffset
		vt.xOffset = 0
	}
	if x < 0 {
		// move the cursor left
		// this may cause offset
		x = -x
		var far int
		for vt.cx > 0 && x > 0 {
			far = vt.runeWidth(vt.content[vt.cy][vt.cx-1])
			if x >= far {
				x -= far
				vt.cx--
			} else {
				break
			}
		}
		if x == 1 && vt.cx != 0 {
			vt.xOffset = -1
		}
	} else if x > 0 {
		// move the cursor right
		// this situation will not cause offset
		for x > 0 {
			// if cx is out of bound, add empty element
			if vt.cx >= len(vt.content[vt.cy])-1 {
				vt.content[vt.cy] = append(vt.content[vt.cy], ' ')
			}
			far := vt.runeWidth(vt.content[vt.cy][vt.cx])
			if far <= x {
				x -= far
				vt.cx++
			} else {
				break
			}
		}
		if x == 1 {
			vt.xOffset = -1
			vt.cx++
		}
	}

	// avoid index out of bound
	vt.cy = max(vt.cy+y, 0)
	for vt.cy >= len(vt.content) {
		vt.content = append(vt.content, []rune{' '})
	}
	for vt.cx >= len(vt.content[vt.cy]) {
		vt.content[vt.cy] = append(vt.content[vt.cy], ' ')
	}
}

// cursorHome can make the cursor back to 0,0
func (vt *VirtualTerm) cursorHome() {
	vt.cx = 0
	vt.cy = 0
	vt.xOffset = 0
}

// writeRune write Rune to content.
func (vt *VirtualTerm) writeRune(r rune) error {
	// if the offset of cursor is not zero, means that there will be non-deterministic for the output
	// For example, if your output is "你好\bCOOL", than it might be "你好OOL"(git bash in Windows) or 你 COOL("Windows powershell")
	// So it should be treated as an error.
	if vt.xOffset != 0 {
		return ErrNonDetermistics
	}
	// get the width of rune
	far := vt.runeWidth(r)
	vt.content[vt.cy][vt.cx] = r
	vt.cursorMove(far, 0)

	return nil
}

// WriteRunes write Rune array to content
func (vt *VirtualTerm) WriteRunes(p []rune) (n int, err error) {
	for i := 0; i < len(p); i++ {
		switch p[i] {
		case '\r':
			// Carriage Return
			vt.cursorMove(-INF, 0)
		case '\n':
			// NewLine
			// If the cursor is on the last line, add a new line
			vt.cursorMove(-INF, 1)
		case '\b':
			vt.cursorMove(-1, 0)
		case '\033':
			idx := i + 1
			for idx < len(p) && (!isCSIFinalRune(p[idx]) || (idx == i+1 && p[idx] == '[')) {
				idx++
			}
			// stopped because it is the final byte
			if idx != len(p) {
				csi := p[i : idx+1]
				err = vt.handleCSI(string(csi))
				if err != nil {
					if errors.Is(err, ErrCannotHandle) {
						vt.log("warning, some sci is not supported")
					} else {
						return 0, err
					}
				}
			}
			i = idx
		default:
			err = vt.writeRune(p[i])
			if err != nil {
				return 0, err
			}
		}
	}
	return len(p), nil
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
		for j, c := range line {
			if j == len(line)-1 {
				break
			}
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
	vt.cursorHome()
	vt.content = emptyContent()
}

// emptyContent get the empty content
func emptyContent() [][]rune {
	runes := make([][]rune, 1, 20)
	runes[0] = append(runes[0], ' ')
	return runes
}

// Process get the output directly without explicitly create a virtual terminal
func Process(input string) (string, error) {
	vt := NewDefault()
	_, err := vt.WriteString(input)
	if err != nil {
		return "", err
	}
	return vt.String(), err
}
