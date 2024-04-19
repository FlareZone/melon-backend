package _type

type AliyunOssRequest struct {
	Storage string `json:"storage" binding:"required,ossStorage"`
	Ext     string `json:"ext" binding:"required,ossImageExt"`
}

type AliyunOssResponse struct {
	SignUrl          string `json:"signUrl"`
	ExpiredTimestamp int64  `json:"expired"`
}
