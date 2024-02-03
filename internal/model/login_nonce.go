package model

type SigNonce struct {
	ID         uint64 `xorm:"pk autoincr 'id'" json:"id"`
	EthAddress string `xorm:"char(42) unique notnull 'eth_address'" json:"eth_address"`
	NonceToken string `xorm:"char(45) unique notnull 'nonce_token'" json:"nonce_token"`
	UsedNonce  int64  `xorm:"index 'used_nonce'" json:"-"`
}

func (s SigNonce) TableName() string {
	return "sig_nonce"
}
