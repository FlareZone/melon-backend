package migrate

import (
	"github.com/FlareZone/melon-backend/internal/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/inconshreveable/log15"
	"xorm.io/xorm"
)

var (
	log = log15.New("m", "migrate")
)

func Schema(dsn string) {
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Error("sync database new xorm fail", "err", err)
		panic(err)
	}
	err = engine.Sync(model.User{}, model.Post{}, model.Topic{}, model.PostTopic{}, model.Group{},
		model.Comment{})
	if err != nil {
		log.Error("sync database migrate fail", "err", err)
		panic(err)
	}
}
