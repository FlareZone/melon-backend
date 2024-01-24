package pages

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"strings"
	"time"
	"xorm.io/builder"
)

type PageRequest struct {
	Size   int    `json:"size"`
	NextID string `json:"next_id"`
}

type PageResponse struct {
	List   interface{} `json:"list"`
	NextID string      `json:"next_id"`
}

func BuildPostListOrders(current *model.Post, orderParam string) (cond builder.Cond, orderBy string) {
	orderArr := strings.Split(orderParam, ",")
	if len(orderArr) == 0 {
		orderArr = append(orderArr, "-created_at")
	}
	var orders []string
	for _, order := range orderArr {
		if strings.EqualFold(order[:1], "-") {
			orders = append(orders, "posts."+order[1:]+" desc")
		} else if strings.EqualFold(order[:1], "+") {
			orders = append(orders, "posts."+order[1:]+" asc")
		}
	}
	orderBy = strings.Join(orders, ",")
	if current.ID <= 0 {
		cond = nil
		return
	}
	switch orderArr[0] {
	case "-created_at":
		cond = builder.Lte{"posts.created_at": current.CreatedAt.Format(time.DateTime)}
	case "+created_at":
		cond = builder.Gte{"posts.created_at": current.CreatedAt.Format(time.DateTime)}
	default:
		cond = nil
	}
	return
}
