package dullhash

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func TestSumAdjacentCollisions(t *testing.T) {
	colls := 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 16; i++ {
		data1, data2 := make([]byte, 256 + i), make([]byte, 256 + i)
		if _, err := rand.Read(data1); err != nil {
			t.Fatalf("error while reading random bytes: %v", err)
		}
		for j := 0; j < len(data1); j++ {
			copy(data2, data1)
			data2[j]++
			sum1, sum2 := Sum(data1), Sum(data2)
			if hex.EncodeToString(sum1[:]) == hex.EncodeToString(sum2[:]) {
				t.Errorf(
					"hash value is the same for data1 and data2 for bytes:\ndata1: %v\ndata2: %v\nat index %v, data1: %v, data2: %v",
					data1,
					data2,
					j,
					data1[j],
					data2[j],
				)
				colls++
			}
		}
	}
	if colls > 0 {
		t.Errorf("total of %v output collisions found", colls)
	}
}
