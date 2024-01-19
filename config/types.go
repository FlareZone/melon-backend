package config

type jwtConfig struct {
	Secret string
	Issuer string
}

type melonDBDsn string

func (m melonDBDsn) String() string {
	return string(m)
}

type app struct {
	Name string
	Url  string
}

type eip712 struct {
	ChainID           int64
	Version           string
	Name              string
	VerifyingContract string
}
