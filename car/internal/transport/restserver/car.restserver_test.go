package restserver

import "testing"

func TestIntentionalCIFailure(t *testing.T) {
	t.Fatal("Intentional CI failure")
}
