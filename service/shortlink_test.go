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
		{1370657502714138624, "0010MoRZ5C1"},
		{1370657502722527235, "301WMoRZ5C1"},
		{1370660755145232384, "0010Q7Ba5C1"},
		{1370657502718332931, "301GMoRZ5C1"},
		{1370660755149426696, "801GQ7Ba5C1"},
	}

	base64 := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_=")

	for _, test := range tests {
		if output := convertToBase64(test.input, base64); output != test.expected {
			t.Errorf("Test failed. %d input, %s expected, but received %s\n", test.input, test.expected, output)
		}
	}
}

func TestNextIdUnique(t *testing.T) {
	node, _ := snowflake.NewNode(1)

	if getNextID(node) == getNextID(node) {
		t.Error("Test failed. Unique id is not being generated.")
	}
}
