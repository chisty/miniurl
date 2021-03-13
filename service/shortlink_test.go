package service

import (
	"testing"

	"github.com/bwmarrin/snowflake"
)

func TestConvertToBase64(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{-1, "-1"},
		{64, "10"},
		{-67, "-13"},
		{1370657502722527235, "1C5ZRoMW103"},
		{1370660755145232384, "1C5aB7Q0100"},
		{1370657502718332931, "1C5ZRoMG103"},
		{-1370660755149426696, "-1C5aB7QG108"},
	}

	base64 := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_=")

	for _, test := range tests {
		if output := convertToBase64(test.input, base64); output != test.expected {
			t.Errorf("Test failed. %d input, %s expected, but received %s\n", test.input, test.expected, output)
		}
	}
}

func TestConvertToBase64Fail(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{-1, "1"},
		{64, "11"},
		{-67, "13"},
		{1370657502722527235, "1C5ZRoMG103"},
		{-1370660755149426696, "1C5aB7QG108"},
	}

	base64 := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_=")

	for _, test := range tests {
		if output := convertToBase64(test.input, base64); output == test.expected {
			t.Errorf("Test failed. %d input, %s expected, but received %s\n", test.input, test.expected, output)
		}
	}
}

func TestNextUniqueId(t *testing.T) {
	node, _ := snowflake.NewNode(1)

	if getNextID(node) == getNextID(node) {
		t.Error("Test failed. Unique id is not being generated.")
	}
}
