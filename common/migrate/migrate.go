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
	log.Info("开始同步数据库！！！")
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Error("sync database new xorm fail", "err", err)
		panic(err)
	}
	err = engine.Sync(
		model.User{}, model.UserFollow{}, model.Point{},
		model.Post{}, model.Topic{}, model.PostTopic{},
		model.PostLike{}, model.PostShare{},
		model.Comment{}, model.CommentLike{},
		model.Group{}, model.UserGroup{},
		model.Asset{},
		model.SigNonce{})
	if err != nil {
		log.Error("sync database migrate fail", "err", err)
		panic(err)
	}
}
