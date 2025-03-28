package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetJobByID(db *gorm.DB, jobID string) (*Job, error) {
	var job Job
	result := db.Where("id = ?", jobID).First(&job)
	if result.Error != nil {
		return nil, result.Error
	}

	return &job, nil
}

func GetLogReportByJobID(db *gorm.DB, jobID string) (*LogReport, error) {
	var logReport LogReport
	result := db.Where("job_id = ?", jobID).First(&logReport)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("record not found")
	}

	type KeywordCount struct {
		Keyword string
		Count   int
	}

	var keywordCounts []KeywordCount
	logReport.TrackedKeywordsCount = make(map[string]int)

	result = db.Table(TrackedKeywordsCount{}.TableName()).Select("keyword", "count").Where("log_report_id = ?", logReport.ID).Scan(&keywordCounts)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, kewWordCount := range keywordCounts {
		logReport.TrackedKeywordsCount[kewWordCount.Keyword] = kewWordCount.Count
	}

	return &logReport, nil
}

type WholeLogReportsAggregate struct {
	UserID           uuid.UUID      `gorm:"column:user_id" json:"userID"`
	TotalLogReports  int            `gorm:"column:total_log_reports" json:"totalLogReports"`
	TotalJobs        int            `gorm:"column:total_jobs" json:"totalJobs"`
	TotalLogs        int            `gorm:"column:total_logs" json:"totalLogs"`
	TrackedKeywords  map[string]int `gorm:"-" json:"totalTrackedKeywords"`
	TotalErrorLogs   int            `gorm:"column:total_error_logs" json:"totalErrorLogs"`
	TotalWarningLogs int            `gorm:"column:total_warning_logs" json:"totalWarningLogs"`
	TotalInfoLogs    int            `gorm:"column:total_info_logs" json:"totalInfoLogs"`
	TotalUniqueIPs   int            `gorm:"column:total_unique_ips" json:"totalUniqueIPs"`
	TotalInvalidLogs int            `gorm:"column:total_invalid_logs" json:"totalInvalidLogs"`
}

func GetWholeLogReportsAggregate(db *gorm.DB, userID uuid.UUID) (*WholeLogReportsAggregate, error) {
	var wholeLogReportsAggregate WholeLogReportsAggregate
	result := db.Raw(`
	SELECT
		jobs.user_id,
		COUNT(DISTINCT log_stats.id) AS total_log_reports,
		COUNT(DISTINCT jobs.id) AS total_jobs,
		COALESCE(SUM(log_stats.total_logs), 0) AS total_logs,
		COALESCE(SUM(log_stats.error_count), 0) AS total_error_logs,
		COALESCE(SUM(log_stats.warn_count), 0) AS total_warning_logs,
		COALESCE(SUM(log_stats.info_count), 0) AS total_info_logs,
		COALESCE(SUM(log_stats.unique_ips), 0) AS total_unique_ips,
		COALESCE(SUM(log_stats.invalid_logs), 0) AS total_invalid_logs
	FROM log_stats
	LEFT JOIN jobs ON log_stats.job_id = jobs.id
	WHERE jobs.user_id = ?
	GROUP BY jobs.user_id
	`, userID).Scan(&wholeLogReportsAggregate)
	if result.Error != nil {
		return nil, result.Error
	}

	type KeywordCount struct {
		Keyword string
		Count   int
	}

	var keywordCounts []KeywordCount

	result = db.Raw(`
	SELECT
		tracked_keywords_counts.keyword,
		COALESCE(SUM(tracked_keywords_counts.count), 0) AS count
	FROM tracked_keywords_counts
	JOIN log_stats ON tracked_keywords_counts.log_report_id = log_stats.id
	JOIN jobs ON log_stats.job_id = jobs.id
	WHERE jobs.user_id = ?
	GROUP BY tracked_keywords_counts.keyword
	`, userID).Scan(&keywordCounts)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Println("keywordCounts", keywordCounts)

	wholeLogReportsAggregate.TrackedKeywords = make(map[string]int)
	for _, kc := range keywordCounts {
		wholeLogReportsAggregate.TrackedKeywords[kc.Keyword] = kc.Count
	}

	return &wholeLogReportsAggregate, nil
}

func AddFailAttemptForJob(db *gorm.DB, jobID string) error {
	result := db.Exec("UPDATE jobs SET attempts = attempts + 1, succeeded = false WHERE id = ?", jobID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}
