package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 100, 5, 0, 100, 5)

	fmt.Println(aoiMgr)
}
