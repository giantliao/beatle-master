package coin

import "C"
import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/giantliao/beatles-master/abires"
	"github.com/giantliao/beatles-master/config"
	"github.com/giantliao/beatles-master/contract"
	"github.com/kprc/libeth/wallet"
	"strings"
	"sync"

	"math/big"
)


type BTLCoinToken struct {
	ethAccessPoint string
	coinAddr string
}

var gBTLCoinTokenInst *BTLCoinToken
var gBTLCoinTokenLock sync.Mutex

func GetBTLCoinToken() *BTLCoinToken {
	if gBTLCoinTokenInst != nil{
		return gBTLCoinTokenInst
	}

	gBTLCoinTokenLock.Lock()
	defer gBTLCoinTokenLock.Unlock()

	if gBTLCoinTokenInst != nil{
		return gBTLCoinTokenInst
	}

	cfg:=config.GetCBtlm()

	gBTLCoinTokenInst = &BTLCoinToken{
		ethAccessPoint: cfg.BTLCAccessPoint,
		coinAddr: cfg.BTLCoinAddr,
	}

	return gBTLCoinTokenInst

}

func (bcw *BTLCoinToken)BtlCoinBalance(addr common.Address) (*big.Int,error)  {
	ec, err := ethclient.Dial(bcw.ethAccessPoint)
	if err != nil {
		return nil,err
	}
	defer ec.Close()
	var btlc *contract.BtlCoin
	btlc,err=contract.NewBtlCoin(common.HexToAddress(bcw.coinAddr),ec)
	if err!=nil{
		return nil, err
	}
	return btlc.BalanceOf(nil,addr)
}

func (bcw *BTLCoinToken)BtlCoinTransfer(toAddr common.Address, tokenNum float64, key *ecdsa.PrivateKey) (hashStr string,err error) {
	ec, err := ethclient.Dial(bcw.ethAccessPoint)
	if err != nil {
		return "",err
	}
	defer ec.Close()
	var btlc *contract.BtlCoin
	btlc,err = contract.NewBtlCoin(common.HexToAddress(bcw.coinAddr),ec)
	if err!=nil {
		return "", err
	}

	opts:=bind.NewKeyedTransactor(key)
	val:=wallet.BalanceEth(tokenNum)

	var tx *types.Transaction

	tx,err = btlc.Transfer(opts,toAddr,val)
	if err!=nil{
		fmt.Println("BTLCoin Transer error",err.Error())
		return "",err
	}

	return tx.Hash().String(),nil
}

func (bcw *BTLCoinToken)CheckHashAndGet(hash common.Hash,nonce uint64,cnt int) (coin float64, fromAddr, toAddr common.Address,err error) {
	var ec *ethclient.Client
	ec, err = ethclient.Dial(bcw.ethAccessPoint)
	if err != nil {
		return
	}

	if cnt <1{
		cnt = 1
	}

	var roundCheck int

	defer ec.Close()
	var tx *types.Transaction
	var isPending bool
	for {
		tx, isPending, err = ec.TransactionByHash(context.TODO(),hash)
		if err!=nil{
			return
		}
		if isPending{
			if roundCheck < cnt {
				log.Println("wait for confirm: ", hash.String())
				time.Sleep(time.Second)
				roundCheck++
				continue
			}else{
				return coin,fromAddr,toAddr,errors.New("pending, waiting")
			}
		}

		if nonce != tx.Nonce(){
			return coin,fromAddr,toAddr,errors.New("not a correct nonce")
		}

		coin,toAddr,err = decodeMethod(tx.Data())
		if err!=nil{
			return
		}
		//tx.AsMessage()
		var chainId *big.Int
		if chainId, err = ec.NetworkID(context.Background()); err != nil {
			return 0, common.Address{}, common.Address{}, err
		}
		var msg types.Message
		if msg, err = tx.AsMessage(types.NewEIP155Signer(chainId)); err != nil {
			return 0, common.Address{}, common.Address{}, err
		}
		fromAddr = msg.From()
		toAddr = *tx.To()

		log.Printf("GetSuccess: coin:%-10.4f fromaddr:%s, toaddr:%s ",coin,fromAddr,toAddr)

		return

	}
}

var gERC20ABI *abi.ABI

func init()  {
	data,err:=abires.Asset("contract/ERC20.abi")
	if err!=nil{
		panic("load ERC20.abi failed: " + err.Error())
	}

	var aj abi.ABI
	aj,err = abi.JSON(strings.NewReader(string(data)))
	if err!=nil{
		panic("abi json failed: "+err.Error())
	}

	gERC20ABI = &aj
}

func decodeMethod(payload []byte) (float64,common.Address,error)  {
	if  bytes.Compare(gERC20ABI.Methods["transfer"].ID, payload[:4]) != 0{
		return 0,common.Address{},errors.New("not a transfer function")
	}
	method,err:=gERC20ABI.MethodById(payload)
	if err!=nil{
		return 0,common.Address{},err
	}

	params:=make(map[string]interface{})
	method.Inputs.UnpackIntoMap(params,payload[4:])

	for k :=range params{
		fmt.Println("key is :",k)
	}
	toAddr := common.Address{}
	if to,ok:=params["to"];ok{
		toAddr = to.(common.Address)
	}

	paramv := 0.0
	if value,ok:=params["value"];ok{
		bigvalue := value.(*big.Int)
		paramv = wallet.BalanceHuman(bigvalue)
	}

	return paramv,toAddr,nil
}

//func TransferERCToken(target string, tokenNo float64, key *ecdsa.PrivateKey) (string, error) {
//
//	t, err := config.SysEthConfig.NewTokenClient()
//	if err != nil {
//		fmt.Println("[TransferERCToken]: tokenConn err:", err.Error())
//		return "", err
//	}
//	defer t.Close()
//
//	opts := bind.NewKeyedTransactor(key)
//	val := util.BalanceEth(tokenNo)
//
//	fmt.Printf("\n----->%.2f", util.BalanceHuman(val))
//
//	tx, err := t.Transfer(opts, common.HexToAddress(target), val)
//	if err != nil {
//		fmt.Println("[TransferERCToken]: Transfer err:", err.Error())
//		return "", err
//	}
//
//	return tx.Hash().Hex(), nil
//}