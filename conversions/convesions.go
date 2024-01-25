package conversions

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

func ParseNanoSecondsToDay(ns int64) int64 {
	interval := time.Duration(ns)
	hours := interval.Hours()

	days := hours / 24
	return int64(days)
}

func ParseDaysToNanoSeconds(days int64) time.Duration {
	nanoSsecondsInt := days * 24 * 60 * 60 * 1e9
	return time.Duration(nanoSsecondsInt)
}

func ParseTransferAmountFloat(token string, amount int64) float64 {
	var divisor int
	if token == "USDC" || token == "USDT" {
		divisor = 6
	} else {
		divisor = 18
	}
	minorFactor := math.Pow10(divisor)
	parsedAmount := float64(amount) / minorFactor

	return parsedAmount
}

// ParseFloatAmountToIntDenomination converts the float64 value to the token's denominational  integer value
// USDC, USDT are denominated in MWei
func ParseFloatAmountToIntDenomination(token string, amount float64) int64 {
	var divisor int
	if token == "USDC" || token == "USDT" {
		divisor = 6
	} else {
		divisor = 18
	}
	minorFactor := math.Pow10(divisor)
	parsedAmount := int64(amount * minorFactor)

	return parsedAmount
}

func ParseTransferAmount(token string, amount float64) *big.Int {
	var divisor int
	if token == "USDC" || token == "USDT" {
		divisor = 6
	} else {
		divisor = 18
	}
	minorFactor := math.Pow10(divisor)
	parsedAmount := int64(amount * minorFactor)

	return big.NewInt(parsedAmount)
}

func ParseAmountToMwei(amount int64) *big.Int {
	etherInMWei := new(big.Int)
	return etherInMWei.SetInt64(amount)
}

func ParseAmountToWei(amount any) (*big.Int, error) {
	etherInWei := new(big.Int)
	etherInWei.SetString("1000000000000000000", 10)

	switch v := amount.(type) {
	case *big.Int:
		weiAmount := new(big.Int).Mul(v, etherInWei)
		return weiAmount, nil
	case *big.Float:
		weiAmount := new(big.Int)
		weiAmountFloat := new(big.Float).Mul(v, big.NewFloat(1e18))
		weiAmountFloat.Int(weiAmount)
		return weiAmount, nil
	default:
		return nil, fmt.Errorf("unsupported input type: %T", amount)
	}
}
