package dullhash

import (
	"encoding/hex"
	"testing"
)

func TestSum(t *testing.T) {
	sum1, sum2 := Sum([]byte("hello world")), Sum([]byte("hello-world"))
	t.Fatalf("%v\n%v\n", hex.EncodeToString(sum1[:]), hex.EncodeToString(sum2[:]))
}
