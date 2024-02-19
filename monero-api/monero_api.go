package moneroapi

import (
	"fmt"
	"strings"
)

func FiatToXmrMarketprice(fiatAmount, xmrMarketPrice float64) string {
	var depositString = fiatAmount / xmrMarketPrice
	depositXmrAmount := fmt.Sprintf("%.12f", depositString)
	split := strings.Split(depositXmrAmount, ".")
	join := strings.Join(split, " ")
	rawXmrDeposit := strings.ReplaceAll(join, " ", "")
	return rawXmrDeposit
}
