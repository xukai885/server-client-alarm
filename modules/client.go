package modules

import "time"

type Client struct {
	Id    string    `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Ip    string    `json:"ip" db:"ip"`
	Count int64     `json:"count" db:"count"`
	Time  time.Time `json:"time"`
}
