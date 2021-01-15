package bootstrap

import (
	"errors"
	"fmt"
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-master/wallet"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/giantliao/beatles-protocol/miners"
	"github.com/giantliao/beatles-protocol/token"
	"github.com/kprc/libgithub"
	"github.com/kprc/nbsnetwork/tools"
	"strings"

	"os"
	"path"
	"strconv"
)

func CollectBootsTrapList(count int) *miners.BootsTrapMiners {
	mdb := db.GetMinersDb()

	btms := &miners.BootsTrapMiners{}
	cnt := 0

	iter := mdb.Iterator()
	for {
		k, v, err := mdb.Next(iter)
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

	w, _ := wallet.GetWallet()

	cfg := config.GetCBtlm()

	btms.EthAccPoint = cfg.EthAccessPoint
	btms.TrxAccPoint = cfg.TrxAccessPoint
	btms.BTLCoinAddr = cfg.BTLCoinAddr
	btms.BTLCPrice = cfg.BTLCoinPrice
	btms.BtlcAccPoint = cfg.BTLCAccessPoint

	for i := 0; i < len(cfg.BootsTrapDownload); i++ {
		btms.NextDownloadPoint = append(btms.NextDownloadPoint, cfg.BootsTrapDownload[i].DownloadPoint)
	}

	btms.BeatlesMasterAddr = w.BtlAddress()
	btms.BeatlesEthAddr = w.AccountString()

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

func Push2Github(idx int) (msg string, err error) {
	cfg := config.GetCBtlm()
	if len(cfg.BootsTrapDownload) <= idx || idx < 0 {
		return "", errors.New("github access point index error")
	}
	content, err := GenBootstrapFileContent()
	if err != nil {
		return "", err
	}

	err = push2github(cfg.BootsTrapDownload[idx], content)
	if err != nil {
		return strconv.Itoa(idx) + " : " + cfg.BootsTrapDownload[idx].String() + " failed", nil
	}

	return strconv.Itoa(idx) + " : " + cfg.BootsTrapDownload[idx].String() + " success", nil
}

func Push2Githubs() (msg string, err error) {
	cfg := config.GetCBtlm()
	if len(cfg.BootsTrapDownload) == 0 {
		return "", errors.New("please input github access point")
	}

	content, err := GenBootstrapFileContent()
	if err != nil {
		return "", err
	}

	succ := ""
	fail := ""

	for i := 0; i < len(cfg.BootsTrapDownload); i++ {
		err = push2github(cfg.BootsTrapDownload[i], content)
		if err != nil {
			if fail != "" {
				fail += "\r\n"
			}
			fail += strconv.Itoa(i) + " : " + cfg.BootsTrapDownload[i].String()
		} else {
			if succ != "" {
				succ += "\r\n"
			}
			succ += strconv.Itoa(i) + " : " + cfg.BootsTrapDownload[i].String()
		}
	}

	if succ != "" {
		msg = "success: \r\n"
	}
	msg += succ
	if succ != "" || fail != "" {
		msg += "\r\n"
	}
	if fail != "" {
		msg += "failed: \r\n"
		msg += fail
	}

	return msg, nil
}

func push2github(ap *config.GithubAccessPoint, content string) error {
	gc := libgithub.NewGithubClient(token.TokenRevert(ap.DownloadPoint.ReadToken),
		ap.DownloadPoint.Owner,
		ap.DownloadPoint.Repository,
		ap.DownloadPoint.Path,
		ap.Name,
		ap.Email)

	_, hash, err := gc.GetContent()
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "404 Not Found") {
			if err != gc.CreateFile("license master create", content) {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	err = gc.UpdateFile2("license master update", content, hash)
	if err != nil {
		fmt.Println("2", err.Error())
		return err
	}
	return nil
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
