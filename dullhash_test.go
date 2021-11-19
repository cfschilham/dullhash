package dullhash

import (
	"encoding/hex"
	"github.com/dgryski/go-onlinestats"
	"gonum.org/v1/gonum/stat"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

const correlationBatchSize = 1000000
var (
	inputs  []float64
	outputs []float64
)

func init() {
	rand.Seed(time.Now().UnixNano())
	inputs, outputs = make([]float64, correlationBatchSize), make([]float64, correlationBatchSize)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = float64(rand.Int63())
	}
	for i := 0; i < correlationBatchSize; i++ {
		sum := Sum(big.NewInt(int64(i)).Bytes())
		sumbi := big.NewInt(0).SetBytes(sum[:])
		outputs[i] = float64(sumbi.Div(sumbi, big.NewInt(4)).Int64())
	}
}

func TestSumAdjacentCollisions(t *testing.T) {
	colls := 0
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

func BenchmarkSum(b *testing.B) {
	startTime := time.Now()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		 _ = Sum([]byte{0})
	}
	b.Logf("hash rate: %.3f MH/s", (float64(b.N) / time.Since(startTime).Seconds()) / 1000000)
}

func TestSumCorrelationCoefficient(t *testing.T) {
	pearsons := stat.Correlation(inputs, outputs, nil)

	if pearsons > .001 || pearsons < -.001 {
		t.Errorf("pearsons correlation coefficient of %v is too high/low, expected [-0.001, 0.001]\n", pearsons)
	}
	t.Logf("pearsons correlation: %v, batch size: %v\n", pearsons, correlationBatchSize)
}

func TestSumSpearmanRhoCorrelationCoefficient(t *testing.T) {
	spearmanr, p := onlinestats.Spearman(inputs, outputs)

	if spearmanr > .001 || spearmanr < -.001 {
		t.Errorf("spearmanr correlation coefficient of %v is too high/low, expected [-0.001, 0.001]\n", spearmanr)
	}
	t.Logf("spearmanr correlation: %v, associated p-value: %v, batch size: %v\n", spearmanr, p, correlationBatchSize)
}

func TestChunkify(t *testing.T) {
	data := []byte{
		0x66, 0xD9, 0x91, 0x01, 0xF8, 0x3C, 0x19, 0x83, 0xEA, 0x1C, 0x80, 0x8E,
		0x77, 0x7D, 0xA0, 0x5D, 0x5A, 0xFE, 0x35, 0xC4, 0x6A, 0x19, 0x96, 0x5F,
		0xE6, 0x98, 0xDE, 0x66, 0x53, 0xE1, 0xA0, 0x30, 0xF1, 0x10, 0x08, 0x5E,
		0x55, 0x6C, 0xC4, 0x52, 0xE5, 0x33, 0x70, 0x1E, 0x3B, 0x2F, 0x1E, 0xB4,
		0x06, 0x1C, 0x9D, 0x0E, 0x23, 0xBA, 0xF0, 0xAB, 0x51, 0x4D, 0x6E, 0x19,
		0xBB, 0x22, 0xA4, 0x55, 0x39, 0xB8, 0xD2, 0xBB, 0xB5, 0xE0, 0x17, 0x62,
		0x14, 0x46, 0x08, 0x72, 0x6A, 0x1E, 0x88, 0xA1,
	}
	expected := [][16]uint32{
		{
			0x66D99101, 0xF83C1983, 0xEA1C808E, 0x777DA05D, 0x5AFE35C4,
			0x6A19965F, 0xE698DE66, 0x53E1A030, 0xF110085E, 0x556CC452,
			0xE533701E, 0x3B2F1EB4, 0x061C9D0E, 0x23BAF0AB, 0x514D6E19,
			0xBB22A455,
		},
		{
			0x39B8D2BB, 0xB5E01762, 0x14460872, 0x6A1E88A1, 0x80000000,
			0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
			0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
			0x00000288,
		},
	}
	chunks := chunkify(data)
	for i, chunk := range chunks {
		for j, n := range chunk {
			if n != expected[i][j] {
				t.Errorf("unexpected value at chunks[%v][%v]: got %08X expected %08X", i, j, n, expected[i][j])
			}
		}
	}
}
