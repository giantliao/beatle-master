package db

import "github.com/kprc/libeth/account"

type LicenseDesc struct {
	Sig string                `json:"sig"`
	CID account.BeatleAddress `json:"-"`
}
