package bootstrap

import (
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/miners"
)

func CollectBootsTrapList(count int) *miners.BootsTrapMiners {
	mdb := db.GetMinersDb()

	btms := &miners.BootsTrapMiners{}
	cnt := 0

	mdb.Iterator()
	for {
		k, v, err := mdb.Next()
		if k == "" || err != nil {
			break
		}
		bm := &miners.Miner{}
		bm.Port = v.Port
		bm.Ipv4Addr = v.Ipv4Addr
		bm.Location = v.Location
		bm.MinerId = v.ID

		btms.Boots = append(btms.Boots, bm)

		cnt++
		if cnt >= count {
			break
		}
	}

	cfg := config.GetCBtlm()

	btms.EthAccPoint = cfg.EthAccessPoint
	btms.TrxAccPoint = cfg.TrxAccessPoint
	btms.NextDownloadPoint = cfg.BootsTrapDownload

	return btms
}
