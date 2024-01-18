package erc4337

import "errors"

const (
	ETHEREUM = "ETH_MAINNET"

	//L2s
	BASE     = "BASE"
	POLYGON  = "POLYGON_POS"
	OPTIMISM = "OPTIMISM"

	// TESTNETS
	MUMBAI          = "POLYGON_MUMBAI"
	GOERLI          = "GOERLI"
	BASE_GOERLI     = "BASE_GOERLI"
	OPTIMISM_GOERLI = "OPTIMISM_GOERLI"
)

func GetNetwork(chainId int64) (string, error) {
	switch chainId {
	case 1:
		return ETHEREUM, nil
	case 5:
		return GOERLI, nil
	case 10:
		return OPTIMISM, nil
	case 137:
		return POLYGON, nil
	case 420:
		return OPTIMISM_GOERLI, nil
	case 84531:
		return BASE_GOERLI, nil
	case 8453:
		return BASE, nil
	case 80001:
		return MUMBAI, nil
	default:
		return "NOT SUPPORTED", errors.New("unsupported chain")
	}
}
