package service

import (
	"github.com/FlareZone/melon-backend/common/contract-interaction"
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
	"xorm.io/xorm"
)

type PointService interface {
	AddPoints(userId, invitedBy string, bonusPoints uint64) bool
	UpdatePoints(userId string, bonusPoints uint64)
	DeletePoints(userId string)
	GetUserPoints(userId string) *model.Point
	// 获取用户排行榜，按积分降序排列
	GetUserLeaderboard() []*model.Point
	ExchangePoints(privateKey, userId string, bonusPoints uint64) bool
}

type Point struct {
	xorm *xorm.Engine
}

func NewPoint(xorm *xorm.Engine) PointService {
	return &Point{xorm: xorm}
}

func (p *Point) ExchangePoints(privateKey, userId string, bonusPoints uint64) bool {

	// 开启事务
	session := p.xorm.NewSession()
	defer session.Close()

	//	查询积分余额是否大于0
	point := &model.Point{}
	has, err := session.Table(&model.Point{}).Where("user_id = ?", userId).Get(point)
	if err != nil {
		log.Error("Failed to check if points exist for user", "user_id", userId, "error", err)
		session.Rollback()
	}
	if !has {
		log.Error("User does not have points", "user_id", userId)
		session.Rollback()
	}
	if point.BonusPoints < bonusPoints {
		log.Error("Insufficient points", "user_id", userId, "required_points", bonusPoints, "available_points", point.BonusPoints)
		session.Rollback()

	}
	//调用合约接口兑换
	contract := contractinteraction.NewProposalLogicContract(privateKey)
	isExchangeSuccess := contract.ExchangePoints(bonusPoints)
	if isExchangeSuccess {
		//兑换成功，则修改积分余额
		point.BonusPoints -= bonusPoints
		session.ID(point.ID).Update(point)
		// 提交事务
		if err := session.Commit(); err != nil {
			log.Error("Failed to commit transaction", "error", err)
		}
		return true
	} else {
		session.Rollback()
	}

	return false
}

func (p *Point) AddPoints(userId, invitedBy string, bonusPoints uint64) (isValidInvite bool) {
	// 开启事务
	session := p.xorm.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Error("Failed to start transaction", "error", err)
	}

	// 检查被邀请用户是否已存在积分记录
	userPoint := &model.Point{}
	has, err := session.Table(&model.Point{}).Where("user_id = ?", userId).Get(userPoint)
	if err != nil {
		log.Error("Failed to check if points exist for user", "user_id", userId, "error", err)
		session.Rollback()
	}
	// 如果不存在，则创建新的积分记录
	if !has {
		userPoint = &model.Point{
			UserId:      userId,
			InvitedBy:   invitedBy,
			InviteDate:  time.Now(),
			BonusPoints: bonusPoints,
		}
		_, err := session.Insert(userPoint)
		if err != nil {
			log.Error("Failed to add points for user", "user_id", userId, "error", err)
			session.Rollback()
		}
		isValidInvite = true
		//	如果是已经存在，但是还未被邀请
	} else if has && userPoint.InvitedBy == "" {
		userPoint.InvitedBy = invitedBy
		userPoint.InviteDate = time.Now()
		userPoint.BonusPoints += bonusPoints
		_, err := session.ID(userPoint.ID).Update(userPoint)
		if err != nil {
			log.Error("Invitation failed", "user_id", userId, "error", err)
			session.Rollback()
		}
		isValidInvite = true

	} else {
		log.Error("The user has been invited", "user_id", userId, "error", err)
		session.Rollback()
	}

	if isValidInvite {
		// 更新邀请人的积分
		invitedByPoint := &model.Point{}
		has, err = session.Where("user_id = ?", invitedBy).Get(invitedByPoint)
		if err != nil {
			log.Error("Failed to check if points exist for invitedBy user", "user_id", invitedBy, "error", err)
			session.Rollback()
		}
		if has {
			//邀请人多得10积分
			invitedByPoint.BonusPoints += bonusPoints + 10
			_, err := session.Where("user_id = ?", invitedBy).Update(invitedByPoint)
			if err != nil {
				log.Error("Failed to update points for invitedBy user", "user_id", invitedBy, "error", err)
				session.Rollback()
			}
		} else {
			invitedByPoint = &model.Point{
				UserId:      invitedBy,
				InvitedBy:   "",
				BonusPoints: bonusPoints + 10,
			}
			_, err := session.Insert(invitedByPoint)
			if err != nil {
				log.Error("Failed to add points for invitedBy user", "user_id", invitedBy, "error", err)
				session.Rollback()
			}
		}
	}

	// 提交事务
	if err := session.Commit(); err != nil {
		log.Error("Failed to commit transaction", "error", err)
	}
	return
}

func (p *Point) UpdatePoints(userId string, bonusPoints uint64) {

	// 更新用户的积分
	_, err := p.xorm.Table(&model.Point{}).Where("user_id = ?", userId).
		Update(&model.Point{BonusPoints: bonusPoints})
	if err != nil {
		log.Error("Failed to update points for user", "user_id", userId, "error", err)
	}
}

func (p *Point) DeletePoints(userId string) {
	// 删除用户的积分
	_, err := p.xorm.Table(&model.Point{}).Where("user_id = ?", userId).Delete(&model.Point{})
	if err != nil {
		log.Error("Failed to delete points for user", "user_id", userId, "error", err)
	}
}
func (p *Point) GetUserPoints(userId string) *model.Point {
	point := &model.Point{}
	has, err := p.xorm.Table(&model.Point{}).Where("user_id = ?", userId).Get(point)
	if err != nil {
		log.Error("Failed to get points for user", "user_id", userId, "error", err)
		return nil
	}
	if !has {
		return nil
	}
	return point
}

func (p *Point) GetUserLeaderboard() []*model.Point {
	var points []*model.Point
	err := p.xorm.Desc("bonus_points").Limit(10).Find(&points)
	if err != nil {
		log.Error("Failed to get all user points", "error", err)
		return nil
	}
	return points
}
