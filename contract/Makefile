SHELL=PATH='$(PATH)' /bin/sh

.PHONY: all
all:
	abigen --abi ERC20.abi --pkg contract --type BtlCoin --out BtlCoin.go

abi:
	solc --abi --overwrite ERC20.sol -o ./

clean:
#	rm -f *.abi
	rm -f *.go
