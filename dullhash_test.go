package dullhash

import (
	"encoding/hex"
	"gonum.org/v1/gonum/stat"
	"math/big"
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
			t.Fatalf("error while reading random bytes: %v\n", err)
		}
		for j := 0; j < len(data1); j++ {
			copy(data2, data1)
			data2[j]++
			sum1, sum2 := Sum(data1), Sum(data2)
			if hex.EncodeToString(sum1[:]) == hex.EncodeToString(sum2[:]) {
				t.Errorf(
					"hash value is the same for data1 and data2 for bytes:\ndata1: %v\ndata2: %v\nat index %v, data1: %v, data2: %v\n",
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
		t.Errorf("total of %v output collisions found\n", colls)
	}
}

func TestSumCorrelationCoefficient(t *testing.T) {
	batchSize := 5000000
	rand.Seed(time.Now().UnixNano())
	inputs, outputs := make([]float64, batchSize), make([]float64, batchSize)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = float64(rand.Int63())
	}
	for i := 0; i < batchSize; i++ {
		sum := Sum(big.NewInt(int64(i)).Bytes())
		sumbi := big.NewInt(0).SetBytes(sum[:])
		outputs[i] = float64(sumbi.Div(sumbi, big.NewInt(4)).Int64())
	}
	coeff := stat.Correlation(inputs, outputs, nil)
	if coeff > .1 || coeff < -.1 {
		t.Errorf("correlation coefficient of %v is too high/low, expected [-0.1, 0.1]\n", coeff)
	}
	t.Logf("correlation: %v, batch size: %v\n", coeff, batchSize)
}
