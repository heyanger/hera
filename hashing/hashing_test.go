package hashing

import (
	"fmt"
	"github.com/funkytennisball/hera/common"
	"math/rand"
	"testing"
)

func TestTest(t *testing.T) {
	m1 := map[common.RangeKey]common.Range{}
	m2 := map[string]common.Node{}

	m1[common.RangeKey{"A", "C", uint64(rand.Intn(32768))}] = common.Range{}
	m1[common.RangeKey{"D", "F", uint64(rand.Intn(32768))}] = common.Range{}

	m2["n1"] = common.Node{}
	m2["n2"] = common.Node{}

	fmt.Println("before:")
	fmt.Println(m1)
	fmt.Println(m2)

	r := HashRing{}
	r.rehash(m1, m2)

	fmt.Println()
	fmt.Println("after:")
	fmt.Println(m1)
	fmt.Println(m2)
}

func TestNoError(t *testing.T) {
	//t.Error("Expected XX, got ", "Y")
}
