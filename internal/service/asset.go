package service

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
	"xorm.io/xorm"
)

type AssetService interface {
	Create(uuid, cosPath string) (asset *model.Asset)
	QueryByUuid(uuid string) (asset *model.Asset)
}

type Asset struct {
	xorm *xorm.Engine
}

func NewAsset(xorm *xorm.Engine) AssetService {
	return &Asset{xorm: xorm}
}

func (a *Asset) Create(uuid, cosPath string) (asset *model.Asset) {
	asset = &model.Asset{
		UUID:      uuid,
		CosPath:   cosPath,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err := a.xorm.Table(&model.Asset{}).Insert(asset)
	if err != nil {
		log.Error("insert to asset table fail", "uuid", uuid, "cos_path", cosPath, "err", err)
	}
	return
}

func (a *Asset) QueryByUuid(uuid string) (asset *model.Asset) {
	asset = new(model.Asset)
	_, err := a.xorm.Table(&model.Asset{}).Where("uuid = ?", uuid).Get(asset)
	if err != nil {
		log.Error("query asset fail", "uuid", uuid, "err", err)
	}
	return
}
