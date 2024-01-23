package handler

type PageRequest struct {
	Size   int    `json:"size"`
	NextID string `json:"next_id"`
}

type PageResponse struct {
	Data   interface{} `json:"data"`
	NextID string      `json:"next_id"`
}
