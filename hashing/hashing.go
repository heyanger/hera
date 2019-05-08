package hashing

import (
	// "container/list"
	"fmt"
	"github.com/funkytennisball/hera/common"
	//"gotest.tools/assert"
	"math/rand"
	"sort"
	"strings"
)

const K int = 5
const R int = 3
const RMAX int = 32768

type HashRing struct {
	ring  []RingLoc //list.List // RingLoc list
	nodes map[string]int
}

type RingLoc struct {
	Location uint64
	Nodes    []string
}

type byLoc []RingLoc

func (ring byLoc) Len() int {
	return len(ring)
}
func (ring byLoc) Swap(i, j int) {
	ring[i], ring[j] = ring[j], ring[i]
}
func (ring byLoc) Less(i, j int) bool {
	return ring[i].Location < ring[j].Location
}

/*func (m map[RangeKey]Range) updateNodes(k RangeKey, Nodes []string) {
	m[k].Nodes = Nodes
}*/

func (r *HashRing) rehash(m1 common.RangeMap, m2 common.NodeMap) {
	// INVARIANT: ring sorted

	// Remove old nodes
	//m2["bb"] = common.Node{}

	// Add new nodes
	flag := false
	for k := range m2 {
		if _, present := r.nodes[k]; !present {

			flag = true
			for j := 0; j < K; j++ {
				r.ring = append(r.ring, RingLoc{uint64(rand.Intn(RMAX)), []string{k}})
			}

			/*i := 0
			for ; i < len(r.ring) && r.ring[i].Location <= 100; i++ {
			}
			r.ring = append(r.ring[:i], RingLoc{}, r.ring[i:])*/

			/*fmt.Println(r.ring)
			e := r.ring.Front()
			for ; e != nil && RingLoc(e.Value).Location < 100; e = e.Next() {
			}
			fmt.Println(e)*/
		}
	}
	if flag {
		//fmt.Println(r.ring)
		sort.Sort(byLoc(r.ring))
		fmt.Println(r.ring)
		for i := 0; i < len(r.ring); i++ {
			if len(r.ring[i].Nodes) < R {
				// Update ring tuple
				for j := (i + 1) % len(r.ring); j != i; j = (j + 1) % len(r.ring) {
					// add to r.ring[i].Nodes
					flag_add := true
					for idx := 0; idx < len(r.ring[i].Nodes); idx++ {
						if strings.Compare(r.ring[i].Nodes[idx], r.ring[j].Nodes[0]) == 0 {
							flag_add = false
							break
						}
					}
					if flag_add {
						r.ring[i].Nodes = append(r.ring[i].Nodes, r.ring[j].Nodes[0])
					}
				}

				// Update m1 mapping
				for k := range m1 {
					// ASSERT: len(r.ring) >= 2
					if k.RingLocation <= r.ring[i].Location &&
						k.RingLocation > r.ring[(i-1+len(r.ring))%len(r.ring)].Location {
						// Remap
						tmp := m1[k]
						tmp.Nodes = r.ring[i].Nodes
						m1[k] = tmp
						//m1[k] = common.Range{r.ring[i].Nodes, 42} //r.ring[i].Nodes
						//fmt.Println(m1[k])
					}
				}
			}
		}
		fmt.Println(r.ring)
	}
}
