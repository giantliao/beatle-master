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

type ReceiptsDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

type ReceiptItem struct {
	ReceiptKey   string                `json:"-"` //for eth,trx is transaction hash
	Cid          account.BeatleAddress `json:"cid"`
	FromType     string                `json:"receipt_from"` //eth trx agent
	FromAddress  string                `json:"from_address"`
	CurrentPrice float64               `json:"current_price"`
	Month        int64                 `json:"month"`
	CreateTime   int64                 `json:"create_time"`
	UpdateTime   int64                 `json:"update_time"`
}

func (ri *ReceiptItem) String() string {
	msg := ""
	msg += fmt.Sprintf("ReceiptKey: %s\r\n", ri.ReceiptKey)
	msg += fmt.Sprintf("Cid:%s, FromType:%s, FromAddress:%s\r\n", ri.Cid.String(), ri.FromType, ri.FromAddress)
	msg += fmt.Sprintf("CurrentPrice: %-8.4f, Month: %d, CreateTime: %s, UpdateTime: %s\r\n",
		ri.CurrentPrice, ri.Month, tools.Int64Time2String(ri.CreateTime), tools.Int64Time2String(ri.UpdateTime))

	return msg
}

var (
	receiptsStore     *ReceiptsDb
	receiptsStoreLock sync.Mutex
)

func newReceiptsStore() *ReceiptsDb {
	cfg := config.GetCBtlm()
	db := db.NewFileDb(cfg.GetReceiptsDbFile()).Load()
	return &ReceiptsDb{NbsDbInter: db}
}

func GetReceiptsDb() *ReceiptsDb {
	if receiptsStore == nil {
		receiptsStoreLock.Lock()
		defer receiptsStoreLock.Unlock()
		if receiptsStore == nil {
			receiptsStore = newReceiptsStore()
		}
	}
	return receiptsStore
}

func (rd *ReceiptsDb) Insert(receiptKey string, cid account.BeatleAddress, fromType string, fromAddr string, currentPrice float64, month int64) error {
	rd.dbLock.Lock()
	defer rd.dbLock.Unlock()

	now := tools.GetNowMsTime()

	if _, err := rd.NbsDbInter.Find(receiptKey); err != nil {
		ri := &ReceiptItem{Cid: cid, FromType: fromType, FromAddress: fromAddr, CurrentPrice: currentPrice, Month: month}
		ri.CreateTime = now
		ri.UpdateTime = now

		j, _ := json.Marshal(*ri)
		err = rd.NbsDbInter.Insert(receiptKey, string(j))

		return err
	} else {
		return errors.New("key is existed, row id is: " + receiptKey)
	}

}

func (rd *ReceiptsDb) Iterator() *db.DBCusor {
	rd.dbLock.Lock()
	defer rd.dbLock.Unlock()

	return rd.NbsDbInter.DBIterator()
}

func (rd *ReceiptsDb) Next(cursor *db.DBCusor) (receiptKey string, rdi *ReceiptItem, err error) {
	if cursor == nil {
		return "", nil, errors.New("initialize cursor first")
	}
	rd.dbLock.Lock()
	k, v := cursor.Next()
	if k == "" {
		rd.dbLock.Unlock()
		return "", nil, errors.New("no receipt in list")
	}
	rd.dbLock.Unlock()

	rdi = &ReceiptItem{}

	if err := json.Unmarshal([]byte(v), rdi); err != nil {
		return "", nil, err
	}

	receiptKey = k
	rdi.ReceiptKey = k

	return

}

func (rd *ReceiptsDb) StringAll() string {
	iter := rd.Iterator()

	msg := ""

	for {
		_, v, err := rd.Next(iter)
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
