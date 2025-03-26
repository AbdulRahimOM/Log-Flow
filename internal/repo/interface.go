package repo

type Repository interface {
	SaveLogMetricsToDB(jobID string, totalLogs, errorCount, warnCount, infoCount, uniqueIPs int) error
}
