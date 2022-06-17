package bloomfile

import (
	"testing"
)

func TestIsBfTxtChunk(t *testing.T) {
	a := "_bf0.example.com."
	b := "example.com."
	c := "exxample.com."
	d := "__bf0.example.com."

	if !isBfTxtChunk(b, a) {
		t.Errorf("should have been true")
	}
	if isBfTxtChunk(c, a) {
		t.Errorf("should have been false")
	}
	if isBfTxtChunk(b, d) {
		t.Errorf("should have been false")
	}
}

func TestExtractGlobalIndex(t *testing.T) {
	a := "_bf1234.example.com."
	b := "example.com."

	if k, _ := extractGlobalIndex(b, a); k != 1234 {
		t.Errorf("should have been true")
	}

	if k, _ := extractGlobalIndex(b, a); k == 234 {
		t.Errorf("should have been false")
	}

	if k, _ := extractGlobalIndex(b, a); k == 123 {
		t.Errorf("should have been false")
	}
}
