package components

import (
	"github.com/FlareZone/melon-backend/config"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	DBEngine *xorm.Engine
)

func InitComponents() {
	// init DBEngine
	engine, err := xorm.NewEngine("mysql", config.MelonDBDsn.String())
	if err != nil {
		panic(err)
	}
	DBEngine = engine
}
