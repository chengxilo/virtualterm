package test

import (
	"fmt"
	"github.com/chengxilo/virtualterm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVirtualTerm(t *testing.T) {
	str := "hello\rvirtuaa\bl-terminal"
	vt := virtualterm.NewDefault()
	vt.Write([]byte(str))
	assert.Equal(t, vt.String(), "virtual-terminal")
}

func ExampleVirtualTerm() {
	str := "hello\rvirtuaa\bl-terminal"
	vt := virtualterm.NewDefault()
	vt.Write([]byte(str))
	fmt.Print(vt.String())
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
		actual := vt.String()
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
		vt.Write([]byte(te.input))
		actual := vt.String()
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
	}
	for _, te := range tests {
		vt.Write([]byte(te.input))
		actual := vt.String()
		assert.Equal(t, te.output, actual)
		vt.Clear()
	}
}

func TestControlSequenceIntroducer(t *testing.T) {
	vt := virtualterm.NewDefault()
	tests := []struct {
		input  string
		output string
	}{
		{"hello\033[Hworl", "worlo"},
		{"\033[123*", ""},
		{"\033[123helloworld", "elloworld"},
		{"\033[323[21helloworld", "21helloworld"},
		{"\033helloworld", "elloworld"},
		{"\033[2Bhello\033[Hworld", "world\n\nhello"},
		{"hello\033[3Dworld\033[2CTo be or not to be\033[3Bsecond hello\033[2Aworld\033[4Bbalabala",
			"heworld  To be or not to be\n                                       world\n\n                           second hello\n\n                                            balabala"},
	}
	for _, te := range tests {
		_, err := vt.Write([]byte(te.input))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, te.output, vt.String())
		vt.Clear()
	}
}
