package bootstrap

import (
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/miners"
)

func GetBootsTrapList(count int) []*miners.BootsTrapMiners {
	mdb := db.GetMinersDb()

	var btms []*miners.BootsTrapMiners
	cnt := 0

	mdb.Iterator()
	for {
		k, v, err := mdb.Next()
		if k == "" || err != nil {
			break
		}
		bm := &miners.BootsTrapMiners{}
		bm.Port = v.Port
		bm.Ipv4Addr = v.Ipv4Addr
		bm.Location = v.Location
		bm.MinerId = v.ID

		btms = append(btms, bm)

		cnt++
		if cnt >= count {
			break
		}
	}
	return btms
}
