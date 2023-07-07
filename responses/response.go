package responses

type ResponseWithData struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type ResponseNoData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
