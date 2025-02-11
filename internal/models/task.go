package models

import "time"

type Task struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	Desc      string    `json:"description" db:"description"`
}

type Worklog struct {
	TaskID  int    `json:"task_id"`
	Desc    string `json:"description"`
	Hours   int    `json:"hours"`
	Minutes int    `json:"minutes"`
}

type GetUserWorklogsRequest struct {
	UserID    int       `form:"user_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" binding:"required"`
}
