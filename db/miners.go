package db

import (
	"encoding/json"
	"github.com/giantliao/beatles-master/config"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type MinersDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	minersStore     *MinersDb
	minersStoreLock sync.Mutex
)

func newMinersStore() *MinersDb {
	cfg := config.GetCBtlm()

	db := db.NewFileDb(cfg.GetMinersDbFile()).Load()

	return &MinersDb{NbsDbInter: db}
}

func GetMinersDb() *MinersDb {
	if minersStore == nil {
		minersStoreLock.Lock()
		defer minersStoreLock.Unlock()

		if minersStore == nil {
			minersStore = newMinersStore()
		}

	}

	return minersStore
}

type MinerDesc struct {
	Ipv4Addr   string                `json:"ipv_4_addr"`
	Port       int                   `json:"port"`
	Location   string                `json:"location"`
	ID         account.BeatleAddress `json:"id"`
	CreateTime int64                 `json:"create_time"`
	UpdateTime int64                 `json:"update_time"`
}

func (mdb *MinersDb) Insert(ipv4 string, port int, location string, id account.BeatleAddress) error {
	minersStoreLock.Lock()
	defer minersStoreLock.Unlock()

	md := MinerDesc{Ipv4Addr: ipv4, Port: port, Location: location, ID: id}

	if _, err := minersStore.NbsDbInter.Find(id.String()); err != nil {
		now := tools.GetNowMsTime()
		md.CreateTime = now
		md.UpdateTime = now
		j, _ := json.Marshal(md)
		minersStore.NbsDbInter.Insert(id.String(), string(j))
		return nil
	} else {
		return err
	}
}

func (mdb *MinersDb) Update(ipv4 string, port int, location string, id account.BeatleAddress) error {
	minersStoreLock.Lock()
	defer minersStoreLock.Unlock()

	return nil
}
