package bootstrap

import (
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/giantliao/beatles-protocol/miners"
	"github.com/kprc/nbsnetwork/tools"
	"os"
	"path"
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

func GenBootstrapFileContent() (string, error) {
	btm := CollectBootsTrapList(8)

	cipherTxt, err := btm.Marshal(miners.SecKey())
	if err != nil {
		return "", err
	}

	m := meta.Meta{}

	m.Marshal("meta sender", cipherTxt)

	return m.ContentS, nil

}

func Save2File(fileName string) error {

	f := fileName

	if !path.IsAbs(fileName) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		f = path.Join(cwd, fileName)
	}

	data, err := GenBootstrapFileContent()
	if err != nil {
		return err
	}

	err = tools.Save2File([]byte(data), f)
	if err != nil {
		return err
	}

	return nil
}
