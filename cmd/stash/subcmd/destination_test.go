package subcmd

import (
	"testing"
)

// This is going to be pretty short because testing user-input functions is tough

// Test that valid directory check works:
func TestIsValidDirectory(t *testing.T) {
	impossibleDir := "/foobar/kraken/who/is/on/first/boom"
	if v, _ := isValidDirectory(impossibleDir); v {
		t.Errorf("There is no way that %s is a valid directory", impossibleDir)
	}

	normalDir := "/tmp"
	if v, _ := isValidDirectory(normalDir); !v {
		t.Errorf("%s should be a valid directory", normalDir)
	}
}
