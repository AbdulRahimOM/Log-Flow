package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ()

type DbTablesWithName interface {
	TableName() string
}

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"column:id;primaryKey"`
	Username  string         `json:"username" gorm:"column:username;unique;not null;index"`
	Email     string         `json:"email" gorm:"column:email"`
	HashedPW  string         `json:"-" gorm:"column:hashed_pw"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;column:deleted_at"`
}

func (u User) TableName() string {
	return "users"
}

type Job struct {
	ID                uuid.UUID `json:"id" gorm:"column:id;primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"column:user_id"`
	FileURL           string    `json:"file_url" gorm:"column:file_url;not null"`
	LogFileUploadedAt time.Time `json:"log_file_uploaded_at" gorm:"column:log_file_uploaded_at"`

	// User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

func (j Job) TableName() string {
	return "jobs"
}

func (job *Job) Create(db *gorm.DB) error {
	job.LogFileUploadedAt = time.Now()
	return db.Create(job).Error
}

type LogReport struct {
	ID                   uuid.UUID      `json:"id" gorm:"column:id;primaryKey"`
	JobID                uuid.UUID      `json:"job_id" gorm:"column:job_id"`
	TotalLogs            int            `json:"total_logs" gorm:"column:total_logs"`
	ErrorCount           int            `json:"error_count" gorm:"column:error_count"`
	WarnCount            int            `json:"warn_count" gorm:"column:warn_count"`
	InfoCount            int            `json:"info_count" gorm:"column:info_count"`
	UniqueIPs            int            `json:"unique_ips" gorm:"column:unique_ips"`
	InvalidLogs          int            `json:"invalid_logs" gorm:"column:invalid_logs"`
	TrackedKeywordsCount map[string]int `json:"tracked_keywords_count" gorm:"-"`
	CreatedAt            time.Time      `json:"created_at" gorm:"column:created_at"`

	Job Job `json:"job" gorm:"foreignKey:JobID;references:ID"`
}

func (lr LogReport) TableName() string {
	return "log_reports"
}

func (lr *LogReport) Create(db *gorm.DB) error {
	lr.CreatedAt = time.Now()
	lr.ID = uuid.New()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(lr).Error; err != nil {
			return fmt.Errorf("Error saving log report: %v", err)
		}

		if len(lr.TrackedKeywordsCount) != 0 {
			trackedKeywordsCounts := make([]TrackedKeywordsCount, 0, len(lr.TrackedKeywordsCount))
			for keyword, count := range lr.TrackedKeywordsCount {
				trackedKeywordsCounts = append(trackedKeywordsCounts, TrackedKeywordsCount{
					LogReportID: lr.ID,
					Keyword:     keyword,
					Count:       count,
				})
			}

			if err := tx.Create(&trackedKeywordsCounts).Error; err != nil {
				return fmt.Errorf("Error saving tracked keywords count: %v", err)
			}
		}

		return nil
	})
	return err

}

type TrackedKeywordsCount struct {
	LogReportID uuid.UUID `json:"log_report_id" gorm:"column:log_report_id;primaryKey"`
	Keyword     string    `json:"keyword" gorm:"column:keyword;primaryKey"`
	Count       int       `json:"count" gorm:"column:count"`

	LogReport LogReport `json:"log_report" gorm:"foreignKey:LogReportID;references:ID"`
}

func (tkc TrackedKeywordsCount) TableName() string {
	return "tracked_keywords_counts"
}
