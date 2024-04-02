package domain

import "time"

type User struct {
	Id        int64
	StudentId string
	Password  string
	Avatar    string
	Nickname  string
	New       bool
	Utime     time.Time
	Ctime     time.Time
}
