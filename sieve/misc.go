package sieve

import (
	"math"
	"math/big"
)

var negativeOne = big.NewInt(-1)
var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)
var oneMillion = big.NewInt(1000000)
var firstFewPrimes = generateFirstFewPrimes()

func isEven(n *big.Int) bool {
	rest := big.NewInt(0)
	rest.Mod(n, two)
	return rest.Cmp(zero) == 0
}

func IsPrime(n *big.Int) bool {

	if n.Sign() == -1 {
		panic("negative numbers disallowed")
	}

	cmpTwo := n.Cmp(two)
	cmpMillion := n.Cmp(oneMillion)

	if cmpTwo == -1 {
		return false
	} else if cmpTwo == 0 {
		return true
	} else {
		if cmpMillion <= 0 {
			return isPrimeBruteForce(n)
		} else {
			if isPrimeFirstFew(n) == false {
				return false
			} else if isPrimeMillerRabin(n) == false {
				return false
			} else if isPrimeBruteForce(n) == false {
				return false
			} else {
				return true
			}
		}
	}
}

func isPrimeFirstFew(n *big.Int) bool {
	rest := big.NewInt(0)

	for _, mod := range firstFewPrimes {
		rest.Mod(n, mod)
		if rest.Cmp(zero) == 0 {
			return false
		}
	}

	return true
}

func isPrimeBruteForce(n *big.Int) bool {

	cmp := n.Cmp(two)
	if cmp == 0 {
		return true
	} else if cmp == -1 {
		return false
	} else if isEven(n) == true {
		return false
	}

	sqrtCeil := SquareRootCeil(n)

	rest := big.NewInt(0)

	for mod := big.NewInt(3); mod.Cmp(sqrtCeil) <= 0; mod.Add(mod, two) {
		rest.Mod(n, mod)
		if rest.Cmp(zero) == 0 {
			return false
		}
	}

	return true
}

func IsPrimeBruteForceSmallInt(n int64) bool {
	if n == 2 {
		return true
	} else if n < 2 || n%2 == 0 {
		return false
	}

	max := int64(math.Sqrt(float64(n))) + 1

	for i := int64(3); i <= max; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func isPrimeMillerRabin(n *big.Int) bool {
	return n.ProbablyPrime(n.BitLen())
}

func generateFirstFewPrimes() []*big.Int {
	var primes [200000]*big.Int

	primes[0] = big.NewInt(2)

	count := 1
	for n := int64(3); n <= 1000000; n += 2 {
		if IsPrimeBruteForceSmallInt(n) == true {
			primes[count] = big.NewInt(n)
			count += 1
		}
	}

	bigPrimes := make([]*big.Int, count)

	for i := 0; i < count; i += 1 {
		bigPrimes[i] = primes[i]
	}

	return bigPrimes
}

func SquareRootCeil(n *big.Int) *big.Int {
	if n.Cmp(one) == -1 {
		panic("cannot get square root of a number smaller than one")
	}

	upperLimit := big.NewInt(2)
	lowerLimit := big.NewInt(2)

	upperLimitExp := big.NewInt(int64(math.Ceil(float64(n.BitLen()) / 2)))
	lowerLimitExp := big.NewInt(int64(math.Floor(float64(n.BitLen()-1) / 2)))

	upperLimit.Exp(upperLimit, upperLimitExp, nil)
	lowerLimit.Exp(lowerLimit, lowerLimitExp, nil)

	middle := big.NewInt(0)

	middleSquared := big.NewInt(0)

	// binary search
	for upperLimit.Cmp(lowerLimit) != 0 {

		if upperLimit.Cmp(lowerLimit) == -1 {
			panic("upperlimit < lowerlimit shouldnt happen")
		}

		middle.Add(upperLimit, lowerLimit)
		middle.Div(middle, two)

		middleSquared.Exp(middle, two, nil)

		if middleSquared.Cmp(n) == -1 {

			if lowerLimit.Cmp(middle) == 0 {
				return upperLimit
			}

			lowerLimit.Set(middle)
		} else {
			upperLimit.Set(middle)
		}
	}

	return upperLimit
}
