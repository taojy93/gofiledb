package gofiledb

type Record struct {
	ID   int         `json:"id"`
	Data interface{} `json:"data"`
}
