package utils

import (
	"reflect"
	"testing"
)

func TestProcessFileAbs(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected File
	}{
		{
			name:  "Folder Input ProcessFile",
			input: "/Users/Phil/Desktop/MediaTest",
			expected: File{
				Path:       "/Users/Phil/Desktop/MediaTest",
				AbsPath:    "/Users/Phil/Desktop/MediaTest",
				PathTokens: []string{"Users", "Phil", "Desktop", "MediaTest"},
			},
		},
		{
			name:  "File Input ProcessFile",
			input: "/Users/Phil/Desktop/MediaTest/oliver.wooHoo_yeehboi.mp4",
			expected: File{
				Path:       "/Users/Phil/Desktop/MediaTest/",
				AbsPath:    "/Users/Phil/Desktop/MediaTest/oliver.wooHoo_yeehboi.mp4",
				PathTokens: []string{"Users", "Phil", "Desktop", "MediaTest", "oliver.wooHoo_yeehboi.mp4"},

				FileName:  "oliver.wooHoo_yeehboi",
				Ext:       ".mp4",
				PrintName: "oliver wooHoo yeehboi",
			},
		},
	}

	// Test the absolute path
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ProcessFile(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("\n\tExpected: %+v \n\tGot: %+v", tt.expected, actual)
			}
		})
	}
}

func TestProcessFileFile(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Mac File Path",
			input:    "/Users/Phil/Desktop/MediaTest",
			expected: "/Users/Phil/Desktop/MediaTest",
		},
		{
			name:     "Windows File Path",
			input:    "C:\\Users\\Phil\\Desktop\\MediaTest",
			expected: "C:\\Users\\Phil\\Desktop\\MediaTest",
		},
	}

	// Test the absolute path
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ProcessFile(tt.input)
			if actual.AbsPath != tt.expected {
				t.Fatalf("\n\tExpected: %s \n\tGot: %s", tt.expected, actual.AbsPath)
			}
		})
	}
}
