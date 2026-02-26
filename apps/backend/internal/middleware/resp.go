package middleware

type CommonResp struct {
	Code    uint8  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
