package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsProtectedFile(t *testing.T) {
	tests := []struct {
		protectedFiles []string
		filePath       string
		expected       bool
	}{
		{[]string{"/etc/passwd", "/etc/shadow"}, "/etc/passwd", true},
		{[]string{"/etc/passwd", "/etc/shadow"}, "/etc/hosts", false},
		{[]string{".env", "config.yaml"}, "config.yaml", true},
		{[]string{".env", "config.yaml"}, "README.md", false},
	}

	for _, test := range tests {
		t.Run(test.filePath, func(t *testing.T) {
			result := isProtectedFile(test.protectedFiles, test.filePath)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		csv       string
		separator string
		expected  []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"  a ; b ; c  ", ";", []string{"a", "b", "c"}},
		{"x|y|z", "|", []string{"x", "y", "z"}},
		{" single ", ",", []string{"single"}},
	}

	for _, test := range tests {
		t.Run(test.csv, func(t *testing.T) {
			result := splitAndTrim(test.csv, test.separator)
			assert.Equal(t, test.expected, result)
		})
	}
}
