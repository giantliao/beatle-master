package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/giantliao/beatles-master/config"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type LicenseDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

type LicenseDesc struct {
	Sig        string                `json:"sig"`
	LastTx     string                `json:"last_tx"`
	ServerId   account.BeatleAddress `json:"server_id"`
	CID        account.BeatleAddress `json:"-"`
	Name       string                `json:"name"`
	Email      string                `json:"email"`
	Cell       string                `json:"cell"`
	ExpireTime int64                 `json:"expire_time"`
	CreateTime int64                 `json:"create_time"`
	UpdateTime int64                 `json:"update_time"`
}

func (ld *LicenseDesc)String() string  {
	msg:=""
	msg += fmt.Sprintf("ServerId: %s  , CID: %s\r\n",ld.ServerId.String(),ld.CID.String())
	msg += fmt.Sprintf("Sig : %s\r\n",ld.Sig)
	msg += fmt.Sprintf("LastTx : %s\r\n",ld.LastTx)
	msg += fmt.Sprintf("Name:%s, Email:%s, Cell:%s\r\n",ld.Name,ld.Email,ld.Cell)
	msg += fmt.Sprintf("ExpireTime: %s, CreateTime:%s,  UpdateTime:%s\r\n",
		tools.Int64Time2String(ld.ExpireTime),
		tools.Int64Time2String(ld.CreateTime),
		tools.Int64Time2String(ld.UpdateTime))

	return msg
}


var (
	licenseStore     *LicenseDb
	licenseStoreLock sync.Mutex
)

func newLicenseStore() *LicenseDb {
	cfg := config.GetCBtlm()

	db := db.NewFileDb(cfg.GetLicensesDbFile()).Load()

	return &LicenseDb{NbsDbInter: db}
}

func GetLicenseDb() *LicenseDb {
	if licenseStore == nil {
		licenseStoreLock.Lock()
		defer licenseStoreLock.Unlock()

		if licenseStore == nil {
			licenseStore = newLicenseStore()
		}
	}
	return licenseStore
}

func (ld *LicenseDb) Insert(cid, serverId account.BeatleAddress, sig, lastTx, name, email, cell string, expireTime int64) error {
	ld.dbLock.Lock()
	defer ld.dbLock.Unlock()

	now := tools.GetNowMsTime()

	if _, err := ld.NbsDbInter.Find(cid.String()); err != nil {
		lDesc := &LicenseDesc{Sig: sig, LastTx: lastTx, Name: name, Email: email, Cell: cell, ServerId: serverId, ExpireTime: expireTime}
		//lDesc.ExpireTime = tools.Moth2Expire(0, month)
		lDesc.CreateTime = now
		lDesc.UpdateTime = now

		j, _ := json.Marshal(*lDesc)
		ld.NbsDbInter.Insert(cid.String(), string(j))

		return nil
	} else {
		return errors.New("key is existed, row id is " + cid.String())
	}
}

func (ld *LicenseDb) Update(cid, serverId account.BeatleAddress, sig, lastTx, name, email, cell string, expireTime int64) error {
	ld.dbLock.Lock()
	defer ld.dbLock.Unlock()

	now := tools.GetNowMsTime()

	if lDescStr, err := ld.NbsDbInter.Find(cid.String()); err != nil {
		lDesc := &LicenseDesc{Sig: sig, LastTx: lastTx, Name: name, Email: email, Cell: cell, ServerId: serverId, ExpireTime: expireTime}
		//lDesc.ExpireTime = tools.Moth2Expire(0, month)
		lDesc.CreateTime = now
		lDesc.UpdateTime = now

		j, _ := json.Marshal(*lDesc)
		ld.NbsDbInter.Insert(cid.String(), string(j))

		return nil
	} else {
		lDesc := &LicenseDesc{}
		json.Unmarshal([]byte(lDescStr), lDesc)

		if sig == lDesc.Sig {
			return errors.New("nothing to update")
		}

		lDesc.ServerId = serverId
		lDesc.Sig = sig
		lDesc.Name = name
		lDesc.Email = email
		lDesc.Cell = cell
		//lDesc.ExpireTime = tools.Moth2Expire(lDesc.ExpireTime, month)
		lDesc.ExpireTime = expireTime
		lDesc.UpdateTime = now

		j, _ := json.Marshal(*lDesc)
		ld.NbsDbInter.Update(cid.String(), string(j))

		return nil
	}
}

func (ld *LicenseDb) Find(cid account.BeatleAddress) *LicenseDesc {
	ld.dbLock.Lock()
	defer ld.dbLock.Unlock()
	if lDescStr, err := ld.NbsDbInter.Find(cid.String()); err != nil {
		return nil
	} else {
		lDesc := &LicenseDesc{}
		json.Unmarshal([]byte(lDescStr), lDesc)
		return lDesc
	}
}

func (ld *LicenseDb) Delete(cid account.BeatleAddress) {
	ld.dbLock.Lock()
	defer ld.dbLock.Unlock()

	licenseStore.NbsDbInter.Delete(cid.String())
}

func (ld *LicenseDb) Save() {
	ld.dbLock.Lock()
	defer ld.dbLock.Unlock()

	ld.NbsDbInter.Save()
}

func (ld *LicenseDb) Iterator() *db.DBCusor{
	ld.dbLock.Lock()
	defer ld.dbLock.Unlock()

	return ld.NbsDbInter.DBIterator()
}

func (ld *LicenseDb) Next(cursor *db.DBCusor) (cid account.BeatleAddress, lDesc *LicenseDesc, err error) {
	if cursor == nil {
		return "", nil, errors.New("initialize cursor first")
	}
	ld.dbLock.Lock()
	k, v := cursor.Next()
	if k == "" {
		ld.dbLock.Unlock()
		return "", nil, errors.New("no license in list")
	}
	ld.dbLock.Unlock()

	lDesc = &LicenseDesc{}
	if err := json.Unmarshal([]byte(v), lDesc); err != nil {
		return "", nil, err
	}

	cid = account.BeatleAddress(k)
	lDesc.CID = cid

	return
}



func (ld *LicenseDb)StringAll() string  {
	iter := ld.Iterator()

	msg := ""

	for {
		_, v, err := ld.Next(iter)
		if err != nil {
			if msg == "" {
				msg = "no miners in db"
			}
			return msg
		}
		if msg != "" {
			msg += "\r\n"
		}
		msg += v.String()
		msg += "============================================"

	}

	return msg
}
