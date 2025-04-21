package users

import "time"

type FakeCall struct {
    ID        uint       `json:"id"`
    UserID    uint       `json:"user_id"`
    CallTime  time.Time `json:"call_time"`
    CreatedAt time.Time `json:"created_at"`
}

