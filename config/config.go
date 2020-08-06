package config

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"os"
	"path"
	"sync"
)

const (
	BTLM_HomeDir      = ".btlmaster"
	BTLM_CFG_FileName = "btlmaster.json"
	BTLM_DB_PATH      = "db"
)

type BtlMasterConf struct {
	EthAccessPoint string `json:"eth_access_point"`
	TrxAccessPoint string `json:"trx_access_point"`

	CmdListenPort  string `json:"cmdlistenport"`
	HttpServerPort int    `json:"http_server_port"`
	WalletSavePath string `json:"wallet_save_path"`

	ApiPath           string `json:"api_path"`
	PurchasePath      string `json:"purchase_path"`
	ListMinerPath     string `json:"list_miner_path"`
	RegisterMinerPath string `json:"register_miner_path"`

	MinersDbPath   string `json:"miners_db_path"`
	LicenseDbPath  string `json:"license_db_path"`
	ReceiptsDbPath string `json:"receipts_db_path"`

	CurrentPrice float64 `json:"-"`
	LastPrice    float64 `json:"-"`

	BootsTrapDownload []string `json:"boots_trap_download"`
}

var (
	btlmcfgInst     *BtlMasterConf
	btmlcfgInstLock sync.Mutex
)

func (bc *BtlMasterConf) InitCfg() *BtlMasterConf {
	bc.HttpServerPort = 50101
	bc.CmdListenPort = "127.0.0.1:50500"
	bc.WalletSavePath = "wallet.json"

	bc.ApiPath = "api"
	bc.PurchasePath = "purchase"
	bc.ListMinerPath = "list"
	bc.RegisterMinerPath = "reg"

	bc.MinersDbPath = "miners.db"
	bc.LicenseDbPath = "miners.db"
	bc.ReceiptsDbPath = "receipts.db"

	bc.CurrentPrice = 0.01
	bc.LastPrice = 0.01

	return bc
}

func (bc *BtlMasterConf) Load() *BtlMasterConf {
	if !tools.FileExists(GetBtlmCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBtlmCFGFile())
	if err != nil {
		log.Println("load file failed", err)
		return nil
	}

	err = json.Unmarshal(jbytes, bc)
	if err != nil {
		log.Println("load configuration unmarshal failed", err)
		return nil
	}

	return bc

}

func newBtlmCfg() *BtlMasterConf {

	bc := &BtlMasterConf{}

	bc.InitCfg()

	return bc
}

func GetCBtlm() *BtlMasterConf {
	if btlmcfgInst == nil {
		btmlcfgInstLock.Lock()
		defer btmlcfgInstLock.Unlock()
		if btlmcfgInst == nil {
			btlmcfgInst = newBtlmCfg()
		}
	}

	return btlmcfgInst
}

func PreLoad() *BtlMasterConf {
	bc := &BtlMasterConf{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BtlMasterConf {
	bc := &BtlMasterConf{}

	bc.InitCfg()

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	btmlcfgInstLock.Lock()
	defer btmlcfgInstLock.Unlock()
	btlmcfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BtlMasterConf) *BtlMasterConf) *BtlMasterConf {
	btmlcfgInstLock.Lock()
	defer btmlcfgInstLock.Unlock()

	lbc := newBtlmCfg().Load()

	if lbc != nil {
		btlmcfgInst = lbc
	} else {
		lbc = newBtlmCfg()
	}

	btlmcfgInst = initfromcmd(lbc)

	return btlmcfgInst
}

func GetBtlmCHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BTLM_HomeDir)
}

func GetBtlmCFGFile() string {
	return path.Join(GetBtlmCHomeDir(), BTLM_CFG_FileName)
}

func (bc *BtlMasterConf) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBtlmCHomeDir()) {
		os.MkdirAll(GetBtlmCHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBtlmCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func (bc *BtlMasterConf) GetPurchasePath() string {
	return "http://" + bc.ApiPath + "/" + bc.PurchasePath
}

func (bc *BtlMasterConf) GetLittMinerPath() string {
	return "http://" + bc.ApiPath + "/" + bc.ListMinerPath
}

func IsInitialized() bool {
	if tools.FileExists(GetBtlmCFGFile()) {
		return true
	}

	return false
}

func (bc *BtlMasterConf) mkdirDbPath() string {
	dbPath := path.Join(GetBtlmCHomeDir(), BTLM_DB_PATH)

	if !tools.FileExists(dbPath) {
		os.MkdirAll(dbPath, 0755)
	}
	return dbPath
}

func (bc *BtlMasterConf) GetMinersDbFile() string {
	return path.Join(bc.mkdirDbPath(), bc.MinersDbPath)
}

func (bc *BtlMasterConf) GetLicensesDbFile() string {
	return path.Join(bc.mkdirDbPath(), bc.LicenseDbPath)
}

func (bc *BtlMasterConf) GetReceiptsDbFile() string {
	return path.Join(bc.mkdirDbPath(), bc.ReceiptsDbPath)
}

func (bc *BtlMasterConf) GetWalletSavePath() string {
	return path.Join(GetBtlmCHomeDir(), bc.WalletSavePath)
}

func (bc *BtlMasterConf) GetpurchaseWebPath() string {
	return "/" + bc.ApiPath + "/" + bc.PurchasePath
}

func (bc *BtlMasterConf) GetListMinersWebPath() string {
	return "/" + bc.ApiPath + "/" + bc.ListMinerPath
}

func (bc *BtlMasterConf) GetRegisterMinerWebPath() string {
	return "/" + bc.ApiPath + "/" + bc.RegisterMinerPath
}
