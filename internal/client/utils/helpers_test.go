package utils

import (
	"os"
	"testing"
)

func TestPrinter_Scan(t *testing.T) {
	origStdin := os.Stdin
	//nolint:reassign //it's for testing
	defer func() { os.Stdin = origStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Error creating pipe: %v", err)
	}
	defer func() {
		_ = r.Close()
	}()

	//nolint:reassign //it's for testing
	os.Stdin = r

	input := "123\n"
	_, err = w.WriteString(input)
	if err != nil {
		t.Fatalf("Error writing to pipe: %v", err)
	}
	_ = w.Close()

	p := Printer{}
	var i int
	n, scanErr := p.Scan(&i)

	if scanErr != nil {
		t.Errorf("Expected no error from Scan, got %v", scanErr)
	}

	if n != 1 {
		t.Errorf("Expected to scan 1 item, scanned %d", n)
	}

	if i != 123 {
		t.Errorf("Expected to scan '123', scanned %d", i)
	}
}
