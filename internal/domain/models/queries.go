package models

import "gorm.io/gorm"

func GetJobByID(db *gorm.DB, id string) (*Job, error) {
	var job Job
	result := db.Where("id = ?", id).First(&job)
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

	return &logReport, nil
}
