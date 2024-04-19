package _type

import (
	_ "embed"
	"encoding/json"
)

//go:embed embed/banners.json
var bannerBytes []byte

type Banner struct {
	Image string `json:"image"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

var (
	Banners []*Banner
)

func init() {
	_ = json.Unmarshal(bannerBytes, &Banners)
}

type BannerListResponse struct {
	Banners []*Banner `json:"banners"`
}

func (b *BannerListResponse) WithBanners(banners []*Banner) *BannerListResponse {
	result := &BannerListResponse{Banners: make([]*Banner, 0)}
	for idx := range banners {
		result.Banners = append(result.Banners, &Banner{Link: banners[idx].Link, Title: banners[idx].Title})
	}
	return result
}
