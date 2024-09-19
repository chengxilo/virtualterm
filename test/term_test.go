package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"virtualterm"
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
		vt.Write([]byte(te.input))
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
