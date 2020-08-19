package price

import (
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-protocol/price"
	"github.com/kprc/nbsnetwork/tools"
)

var (
	ethprice        float64
	lastRequestTime int64
)

func GetEthPrice() float64 {

	if tools.GetNowMsTime()-lastRequestTime < 300000 {
		return ethprice
	}

	p, err := price.GetPrice()
	if err != nil {
		if ethprice == 0 {
			ethprice = config.GetCBtlm().CurrentEthPrice
		}
		return ethprice
	}

	lastRequestTime = tools.GetNowMsTime()

	ethprice = p.EthPrice.USD

	return p.EthPrice.USD
}
