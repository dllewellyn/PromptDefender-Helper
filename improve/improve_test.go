package improve

import (
	"testing"
)

func TestRandomStringLength(t *testing.T) {
	length := 10
	result := RandomString(length)
	if len(result) != length {
		t.Fatalf("expected length %d, got %d", length, len(result))
	}
}

func TestRandomStringCharacters(t *testing.T) {
	length := 20
	result := RandomString(length)
	for _, char := range result {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			t.Fatalf("unexpected character %c in result", char)
		}
	}
}

func TestRandomStringUniqueness(t *testing.T) {
	length := 15
	result1 := RandomString(length)
	result2 := RandomString(length)
	if result1 == result2 {
		t.Fatalf("expected different strings, got %s and %s", result1, result2)
	}
}
