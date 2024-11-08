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

func TestNewLlmImprover(t *testing.T) {
	mockImproveFunc := func(input string, randomSequence string) (string, error) {
		return "Improved: " + input, nil
	}

	improver := NewLlmImprover(mockImproveFunc)
	if improver == nil {
		t.Fatalf("expected non-nil LlmImprover")
	}
}

func TestLlmImprover_Improve(t *testing.T) {
	mockImproveFunc := func(input string, randomSequence string) (string, error) {
		return "Improved: " + input, nil
	}

	improver := NewLlmImprover(mockImproveFunc)
	result, err := improver.Improve("test input")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Improved: test input"
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestLlmPromptImproverInput(t *testing.T) {
	input := LlmPromptImproverInput{
		StartingPrompt: "test prompt",
		RandomSequence: "random sequence",
	}

	if input.StartingPrompt != "test prompt" {
		t.Fatalf("expected %v, got %v", "test prompt", input.StartingPrompt)
	}

	if input.RandomSequence != "random sequence" {
		t.Fatalf("expected %v, got %v", "random sequence", input.RandomSequence)
	}
}

func TestNewLlmImprover_NilFunc(t *testing.T) {
	err := NewLlmImprover(nil)
	if err != nil {
		t.Fatalf("expected response to be nil")
	}
}
