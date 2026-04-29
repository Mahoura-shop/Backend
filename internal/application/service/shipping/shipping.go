package shipping

// Flat shipping rates per province ID (in Tomans).
// Province IDs match the seeded provinces (Tehran = 1).
var provinceRates = map[uint]uint{
	1: 30000,  // Tehran
	2: 50000,  // Isfahan
	3: 50000,  // Fars
	4: 50000,  // Khorasan Razavi
	5: 60000,  // East Azerbaijan
	6: 60000,  // West Azerbaijan
	7: 60000,  // Gilan
	8: 60000,  // Mazandaran
	9: 60000,  // Khuzestan
	10: 60000, // Kerman
}

const defaultRate uint = 70000

func CalculateShipping(provinceID uint) uint {
	if rate, ok := provinceRates[provinceID]; ok {
		return rate
	}
	return defaultRate
}
