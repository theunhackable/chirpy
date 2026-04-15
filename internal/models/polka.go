package model

type PolkaData struct {
	UserId string `json:"user_id"`
}

type PolkaReq struct {
	Event string    `json:"event"`
	Data  PolkaData `json:"data"`
}
