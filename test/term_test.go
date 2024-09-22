package test

import (
	"errors"
	"fmt"
	"github.com/chengxilo/virtualterm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestVirtualTerm(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"hello\rvirtuaa\bl-terminal", "virtual-terminal"},
		{"100% |██████████| (100/100 B, 100 B/s, 100 it/s)", "100% |██████████| (100/100 B, 100 B/s, 100 it/s)"},
		{"🦐🦑🐙\r🦞🦀\n🐚\b\b🦆🐓", "🦞🦀🐙\n🦆🐓"},
	}
	for _, te := range tests {
		s, err := virtualterm.Process(te.input)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, te.output, s)
	}
}

func ExampleVirtualTerm() {
	str := "hello\rvirtuaa\bl-terminal"
	ns, _ := virtualterm.Process(str)
	fmt.Print(ns)
	// Output:
	// virtual-terminal
}

func TestCarriageReturn(t *testing.T) {
	vt := virtualterm.NewDefault()
	tests := []struct {
		input  string
		output string
	}{
		{"cute\rhat", "hate"},
		{"hello\rworld", "world"},
		{"strawberry\rpie", "pieawberry"},
		{"pie\rapple", "apple"},
	}
	for _, te := range tests {
		_, err := vt.Write([]byte(te.input))
		if err != nil {
			t.Fatal(err)
		}
		actual, _ := vt.String()
		assert.Equal(t, te.output, actual)
		vt.Clear()
	}
}

func TestNewLine(t *testing.T) {
	vt := virtualterm.NewDefault()
	tests := []struct {
		input  string
		output string
	}{
		{"hello\nworld", "hello\nworld"},
		{"strawberry\npie", "strawberry\npie"},
		{"pie\napple", "pie\napple"},
		{"hello\nI am\na hobo.\n", "hello\nI am\na hobo.\n"},
		{"I want to be the\rworld444\nWhat is this\n", "world444o be the\nWhat is this\n"},
	}
	for _, te := range tests {
		_, err := vt.Write([]byte(te.input))
		if err != nil {
			t.Fatal(err)
		}
		actual, _ := vt.String()
		assert.Equal(t, te.output, actual)
		vt.Clear()
	}
}

func TestBackspace(t *testing.T) {
	vt := virtualterm.NewDefault()
	tests := []struct {
		input  string
		output string
	}{
		{"I want to be the\rworld444\nWhat is this\bh\n", "world444o be the\nWhat is thih\n"},
		{"I want to be the\rworld444\nWhat is this\n\bh\n", "world444o be the\nWhat is this\nh\n"},
		{"锄禾\b\b日\b\b\b\b当午", "当午"},
	}
	for _, te := range tests {
		_, err := vt.Write([]byte(te.input))
		if err != nil {
			t.Fatal(err)
		}
		actual, _ := vt.String()
		assert.Equal(t, te.output, actual)
		vt.Clear()
	}
}

func TestCSI(t *testing.T) {
	vt := virtualterm.NewDefault()
	tests := []struct {
		input    string
		expected string
	}{
		{"123\r嗨", "嗨3"},
		{"\033[123*", ""},
		{"你好\r\033[4C啊", "你好啊"},
		{"你好\r\033[C", "你好"},
		{"Hello\033[1D\bWorld", "HelWorld"},
		{"Hello\033[2D\b\bWorld", "HWorld"},
		{"Hello\033[5CWorld", "Hello     World"},
		{"Hello\033[1AWorld", "HelloWorld"},
		{"Hello\033[1BWorld", "Hello\n     World"},
		{"Hello\033[1D\033[1CWorld", "HelloWorld"},
		{"Hello\033[1A\033[1BWorld", "Hello\n     World"},
		{"Hello\033[2D\b\b\033[2CWorld", "HelWorld"},
		{"Hello\033[1A\033[1B\033[1D\bWorld", "Hello\n   World"},
		{"\033helloworld", "elloworld"},
		{"\033[2Bhello\033[Hworld", "world\n\nhello"},
		{"hello\033[3Dworld\033[2CTo be or not to be\033[3Bsecond hello\033[2Aworld\033[4Bbalabala",
			"heworld  To be or not to be\n                                       world\n\n                           second hello\n\n                                            balabala"},
		{"锄禾日\b\033[1C当午", "锄禾日当午"},
		{"你好\b\033[2C呀", "你好 呀"},
		{"锄禾日当午\n汗滴禾下土\033[H谁知盘中餐\r粒粒皆辛苦", "粒粒皆辛苦\n汗滴禾下土"},
		{"🏀🐤\033[2D\b\b🐓你太美", "🐓你太美"},
		{"\033[323[21helloworld", "21helloworld"},
		{"\033[123helloworld", "elloworld"},
		{"\bレ\033[2Dモン", "モン"},
	}
	for i, te := range tests {
		_, err := vt.WriteString(te.input)
		if err != nil {
			t.Fatal(err, i)
		}
		actual, _ := vt.String()
		if actual != te.expected {
			log.Print("actual: "+actual+"expected: ", te.expected, "test index: ", i)
			t.Fail()
		}
		vt.Clear()
	}
}

func TestInvalidInput(t *testing.T) {
	vt := virtualterm.NewDefault()
	tests := []struct {
		input string
	}{
		{"你好\ba"},
		{"我是\b猫"},
		{"我是\033[1D猫"},
		{"我是\bhero"},
		{"锄禾日\b\033当[1C午"},
		{"\bレ\033[Dモン"},
	}
	for _, te := range tests {
		vt.WriteString(te.input)
		if _, err := vt.String(); !errors.Is(err, virtualterm.ErrNonDeterministic) {
			t.Fatal(err)
		}
		vt.Clear()
	}
}
