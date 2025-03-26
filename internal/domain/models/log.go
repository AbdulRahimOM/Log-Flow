package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"column:id;primaryKey"`
	Username  string         `json:"username" gorm:"column:username;unique;not null;index"`
	Email     string         `json:"email" gorm:"column:email"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;column:deleted_at"`
}

type UserCredentials struct {
	ID       string `json:"id" gorm:"column:id;primaryKey"`
	HashedPW string `json:"hashed_pw" gorm:"column:hashed_pw"`

	User User `json:"user" gorm:"foreignKey:ID;references:ID"`
}

type Job struct {
	ID                int       `json:"id" gorm:"column:id;primaryKey"`
	UserID            string    `json:"user_id" gorm:"column:user_id"`
	FileURL           string    `json:"file_url" gorm:"column:file_url;unique;not null"`
	LogFileUploadedAt time.Time `json:"log_file_uploaded_at" gorm:"column:log_file_uploaded_at"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

type LogReport struct {
	ID                       int       `json:"id" gorm:"column:id;primaryKey"`
	JobID                    string    `json:"job_id" gorm:"column:job_id"`
	TotalLogs                int       `json:"total_logs" gorm:"column:total_logs"`
	ErrorCount               int       `json:"error_count" gorm:"column:error_count"`
	WarnCount                int       `json:"warn_count" gorm:"column:warn_count"`
	InfoCount                int       `json:"info_count" gorm:"column:info_count"`
	UniqueIPs                int       `json:"unique_ips" gorm:"column:unique_ips"`
	ErrorWhileParsing        bool      `json:"error_while_parsing" gorm:"column:error_while_parsing"`
	ErrorWhileParsingMessage string    `json:"error_while_parsing_message" gorm:"column:error_while_parsing_message"`
	CreatedAt                time.Time `json:"created_at" gorm:"column:created_at"`

	Job Job `json:"job" gorm:"foreignKey:JobID;references:ID"`
}

type TrackedKeywordsCount struct {
	LogReportID int    `json:"log_report_id" gorm:"column:log_report_id;primaryKey"`
	Keyword     string `json:"keyword" gorm:"column:keyword;primaryKey"`
	Count       int    `json:"count" gorm:"column:count"`

	LogReport LogReport `json:"log_report" gorm:"foreignKey:LogReportID;references:ID"`
}
