package main
import (
  "testing"
)

func TestPlus(t *testing.T) {
  got := Plus(2, 2)
  want := 5

  if got != want {
    t.Errorf("Got %v, wanted %v", got, want)
  }
}
