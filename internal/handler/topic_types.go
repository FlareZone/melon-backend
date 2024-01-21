package handler

type TopicCreateParamRequest struct {
	Name string `json:"name" binding:"required,max=256"`
}

type TopicResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
