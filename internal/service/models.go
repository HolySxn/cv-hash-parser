package service

import "time"

type Report struct {
	CvURL     string    `json:"cv_url"`
	Hash      string    `json:"hash"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}
