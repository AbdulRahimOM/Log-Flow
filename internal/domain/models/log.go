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

type Job struct {
	ID         uuid.UUID `json:"id" gorm:"column:id;primaryKey"`
	UserID     uuid.UUID `json:"userID" gorm:"column:user_id"`
	FileURL    string    `json:"fileURL" gorm:"column:file_url;not null"`
	Attempts   int       `json:"attempts" gorm:"column:attempts;default:0"`
	Succeeded  bool      `json:"succeeded" gorm:"column:succeeded;default:false"`
	UploadedAt time.Time `json:"uploadedAt" gorm:"column:uploaded_at"`
}

func (j Job) TableName() string {
	return "jobs"
}

func (job *Job) Create(db *gorm.DB) error {
	job.UploadedAt = time.Now()
	return db.Create(job).Error
}

type LogReport struct {
	ID                   uuid.UUID      `json:"id" gorm:"column:id;primaryKey"`
	JobID                uuid.UUID      `json:"jobID" gorm:"column:job_id"`
	TotalLogs            int            `json:"totalLogs" gorm:"column:total_logs"`
	ErrorCount           int            `json:"errorCount" gorm:"column:error_count"`
	WarnCount            int            `json:"warnCount" gorm:"column:warn_count"`
	InfoCount            int            `json:"infoCount" gorm:"column:info_count"`
	UniqueIPs            int            `json:"uniqueIPs" gorm:"column:unique_ips"`
	InvalidLogs          int            `json:"invalidLogs" gorm:"column:invalid_logs"`
	TrackedKeywordsCount map[string]int `json:"trackedKeywords_count" gorm:"-"`
	CreatedAt            time.Time      `json:"createdAt" gorm:"column:created_at"`

	Job Job `json:"-" gorm:"foreignKey:JobID;references:ID"`
}

func (lr LogReport) TableName() string {
	return "log_stats"
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

		err := tx.Exec("UPDATE jobs SET succeeded = true, attempts = attempts + 1 WHERE id = ?", lr.JobID).Error
		if err != nil {
			return fmt.Errorf("Error updating job: %v", err)
		}

		return nil
	})
	return err

}

type TrackedKeywordsCount struct {
	LogReportID uuid.UUID `json:"logReportID" gorm:"column:log_report_id;primaryKey"`
	Keyword     string    `json:"keyword" gorm:"column:keyword;primaryKey"`
	Count       int       `json:"count" gorm:"column:count"`

	LogReport LogReport `json:"-" gorm:"foreignKey:LogReportID;references:ID"`
}

func (tkc TrackedKeywordsCount) TableName() string {
	return "tracked_keywords_counts"
}
