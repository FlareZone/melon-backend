package config

type jwtConfig struct {
	Secret string
	Issuer string
}

type melonDBDsn string

func (m melonDBDsn) String() string {
	return string(m)
}
