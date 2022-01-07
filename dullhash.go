package dullhash

import (
	"encoding/binary"
	"github.com/cfschilham/dullhash/sieve"
	"math"
	"math/big"
)

const (
	h0d uint32 = 0x64756C6C
	h1d uint32 = 0x68617368
	h2d uint32 = 0x20697320
	h3d uint32 = 0x6120706F
	h4d uint32 = 0x6F722068
	h5d uint32 = 0x61736820
	h6d uint32 = 0x66756E63
	h7d uint32 = 0x74696F6E
)

var MaxSum = [32]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255,
}

func chunkify(data []byte) [][16]uint32 {
	data = append(data, 128)

	// Data length in bits to store in the final 64 bits of the last chunk.
	dataLen := uint64(len(data) * 8)

	// If less than 64 bits remain in what would have been the final chunk, not
	// enough space is left to append the data length at the end. Instead 0s are
	// appended until the chunk is full and the length will be stored in the
	// next chunk.
	if len(data)%64 > 56 {
		for i := 0; i < len(data)%64; i++ {
			data = append(data, 0)
		}
	}

	// Make sure the final chunk can be divided into 4-byte pieces.
	for i := 0; i < (len(data)%64)%4; i++ {
		data = append(data, 0)
	}

	chunks := make([][16]uint32, (len(data)/64)+1)

	// Make all chunks of uint32s besides the final chunk.
	for i := 0; i < len(chunks)-1; i++ {
		for j := 0; j < 16; j++ {
			chunks[i][j] = binary.BigEndian.Uint32(data[(i*64)+(j*4) : (i*64)+(j*4)+4])
		}
	}

	// Make the final chunk.
	for i := 0; i < (len(data)%64)/4; i++ {
		chunks[len(chunks)-1][i] = binary.BigEndian.Uint32(data[((len(chunks)-1)*64)+(i*4) : ((len(chunks)-1)*64)+(i*4)+4])
	}

	// Set the final two uint32s of the final chunk equal to the length of the
	// initial data.
	chunks[len(chunks)-1][14] = uint32(dataLen >> 32)
	chunks[len(chunks)-1][15] = uint32(dataLen & 0x0000FFFF)

	return chunks
}

func addOverflow(x, y uint32) uint32 {
	if y > x {
		x, y = y, x
	}
	if y > math.MaxUint32-x {
		return y - (math.MaxUint32 - x)
	}
	return x + y
}

func leftRotate(x, n uint32) uint32 {
	n %= 32
	return x<<n | x>>(32-n)
}

func rightRotate(x, n uint32) uint32 {
	n %= 32
	return x>>n | x<<(32-n)
}

func Sum(data []byte) [32]byte {
	chunks := chunkify(data)
	h0, h1, h2, h3, h4, h5, h6, h7 := h0d, h1d, h2d, h3d, h4d, h5d, h6d, h7d
	for _, chunk := range chunks {
		a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7

		for i := 0; i < 8; i++ {
			for i := 0; i < len(chunk); i++ {
				x := a ^ d ^ rightRotate(chunk[i], 11)
				y := f ^ h ^ (x & e)
				z := (b << 10) | (y >> 22) ^ c

				a, b, c, d, e, f, g, h = addOverflow(b, x), h, f, addOverflow(e, z), addOverflow(d, y), a, c, g
			}
		}

		h0, h1, h2, h3 = addOverflow(h0, a), addOverflow(h1, b), addOverflow(h2, c), addOverflow(h3, d)
		h4, h5, h6, h7 = addOverflow(h4, e), addOverflow(h5, f), addOverflow(h6, g), addOverflow(h7, h)
	}

	sum := [32]byte{}
	for i, n := range []uint32{h0, h1, h2, h3, h4, h5, h6, h7} {
		nbe := make([]byte, 4)
		binary.BigEndian.PutUint32(nbe, n)
		sum[i*4], sum[(i*4)+1], sum[(i*4)+2], sum[(i*4)+3] = nbe[0], nbe[1], nbe[2], nbe[3]
	}
	return sum
}

func UsefulSum(data []byte) ([32]byte, []*big.Int) {
	sum := Sum(data)
	// TODO: determine the size of the integer to be extracted from the hash.
	n := big.NewInt(0).SetBytes(sum[:8]) // First 64 bits of the hash (for now).
	f1, f2 := sieve.Factorize(n)
	if f1 != nil && f2 != nil {
		data = append(data, f1.Bytes()...)
		data = append(data, f2.Bytes()...)
		return Sum(data), []*big.Int{f1, f2}
	}
	return sum, nil
}
