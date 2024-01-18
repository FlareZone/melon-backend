package uuid

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	fmt.Println(s)
	fmt.Println(len(s))
}
