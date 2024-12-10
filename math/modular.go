package math

func GCDExtended(a, b int) (int, int, int) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, xPrime, yPrime := GCDExtended(b%a, a)
	return gcd, yPrime - (b/a)*xPrime, xPrime
}

func ModularInverse(x, base int) int {
	gcd, x, _ := GCDExtended(x, base)
	if gcd != 1 || base == 0 {
		panic("Failed to get modular inverse")
	}

	return (base + (x % base)) % base
}

func Exponentiation(base, exponent, mod int64) int64 {
	if mod == 1 {
		return 0
	}

	if exponent < 0 {
		return -1
	}

	var result int64 = 1

	base = base % mod

	for exponent > 0 {
		if exponent%2 == 1 {
			result = (result * base) % mod
		}
		exponent = exponent >> 1
		base = (base * base) % mod
	}
	return result
}
