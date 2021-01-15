package api

import (
	"context"
	"github.com/giantliao/beatles-master/app/cmdcommon"
	"github.com/giantliao/beatles-master/app/cmdpb"
	"github.com/giantliao/beatles-master/bootstrap"
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/db"
	"github.com/giantliao/beatles-master/wallet"
	"github.com/kprc/libeth/account"
	"strconv"

	"time"
)

type CmdStringOPSrv struct {
}

func (cso *CmdStringOPSrv) StringOpDo(cxt context.Context, so *cmdpb.StringOP) (*cmdpb.DefaultResp, error) {
	msg := ""
	switch so.Op {
	case cmdcommon.CMD_WALLET_SHOW:
		msg = cso.showWallet()
	case cmdcommon.CMD_WALLET_LOAD:
		msg = cso.loadWallet(so.Param[0])
	case cmdcommon.CMD_WALLET_CREATE:
		msg = cso.createWallet(so.Param[0])
	case cmdcommon.CMD_BOOTSTRAP_LIST:
		msg = cso.bootstrapList(so.Param[0])
	case cmdcommon.CMD_BOOTSTRAP_ADD:
		msg = cso.bootstrapAdd(so.Param[0], so.Param[1], so.Param[2], so.Param[3], so.Param[4], so.Param[5])
	case cmdcommon.CMD_BOOTSTRAP_DEL:
		msg = cso.bootstrapDel(so.Param[0])
	case cmdcommon.CMD_BOOTSTRAP_PUSH:
		msg = cso.bootstrapPush(so.Param[0])
	case cmdcommon.CMD_MINER_REMOVE:
		msg = cso.removeMiner(so.Param[0])
	case cmdcommon.CMD_DB_SHOW:
		msg = cso.dbShow(so.Param[0])
	case cmdcommon.CMD_DB_SAVE:
		msg = cso.dbSave(so.Param[0])
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}

func int64time2string(t int64) string {
	tm := time.Unix(t/1000, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func (cso *CmdStringOPSrv) showWallet() string {
	if _, err := wallet.GetWallet(); err != nil {
		return err.Error()
	} else {
		var s string
		if s, err = wallet.ShowWallet(); err != nil {
			return err.Error()
		}

		return s
	}
}

func (cso *CmdStringOPSrv) createWallet(auth string) string {
	if err := wallet.LoadWallet(auth); err != nil {
		return err.Error()
	}

	if w, err := wallet.GetWallet(); err != nil {
		return err.Error()
	} else {
		return "create wallet successful, beatles address is : " + w.BtlAddress().String()
	}
}

func (cso *CmdStringOPSrv) loadWallet(auth string) string {
	if err := wallet.LoadWallet(auth); err != nil {
		return err.Error()
	}

	if w, err := wallet.GetWallet(); err != nil {
		return err.Error()
	} else {
		return "load wallet successful, beatles address is : " + w.BtlAddress().String()
	}
}

func (cso *CmdStringOPSrv) bootstrapList(filename string) string {
	if err := bootstrap.Save2File(filename); err != nil {
		return err.Error()
	}

	return "save to file: " + filename + " successful"
}

func (cso *CmdStringOPSrv) bootstrapAdd(owner, repository, filePath, readToken, name, email string) string {
	cfg := config.GetCBtlm()

	err := cfg.AddBootstrap(owner, repository, filePath, readToken, name, email)

	if err != nil {
		return err.Error()
	}

	return "add bootstrap server successful"
}

func (cso *CmdStringOPSrv) bootstrapDel(idxstr string) string {

	idx, err := strconv.Atoi(idxstr)
	if err != nil {
		return err.Error()
	}

	cfg := config.GetCBtlm()
	err = cfg.DelBootstrap(idx)
	if err != nil {
		return err.Error()
	}

	return "delete bootstrap server index: " + idxstr + " successful"
}

func (cso *CmdStringOPSrv) bootstrapPush(idxstr string) string {
	idx, err := strconv.Atoi(idxstr)
	if err != nil {
		return err.Error()
	}
	var msg string

	msg, err = bootstrap.Push2Github(idx)
	if err != nil {
		return err.Error()
	}

	return msg
}
func (cso *CmdStringOPSrv) removeMiner(id string) string {
	mdb := db.GetMinersDb()
	if _, err := mdb.Find(account.BeatleAddress(id)); err != nil {
		return "id: " + id + " not found"
	}

	mdb.Delete(account.BeatleAddress(id))

	return "id: " + id + " delete successfully"

}

func (cso *CmdStringOPSrv) dbShow(dbname string) string {
	if dbname == config.MinersDBName {
		return db.GetMinersDb().StringAll()
	} else if dbname == config.LicenseDBName {
		return db.GetLicenseDb().StringAll()
	} else if dbname == config.ReceiptsName {
		return db.GetReceiptsDb().StringAll()
	} else {
		return "no " + dbname + " !"
	}

}
func (cso *CmdStringOPSrv) dbSave(dbname string) string {
	if dbname == config.MinersDBName {
		db.GetMinersDb().Save()
	} else if dbname == config.LicenseDBName {
		db.GetLicenseDb().Save()
	} else if dbname == config.ReceiptsName {
		db.GetReceiptsDb().Save()
	} else {
		return "no " + dbname + " !"
	}

	return "save " + dbname + " success"
}
