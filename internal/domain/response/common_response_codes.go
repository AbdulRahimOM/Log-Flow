package response


const (
	Success = "SUCCESS"
	Created = "CREATED"

	InternalServerError = "INTERNAL_SERVER_ERROR"
	Unauthorized        = "UNAUTHORIZED"
	Forbidden           = "FORBIDDEN"

	ParsingFormDataErr   = "PARSING_FORM_DATA_ERROR"
	BadRequest           = "BAD REQUEST DATA"
	WrongInput           = "WRONG_INPUT"
	InvalidURLParam      = "INVALID_URL_PARAM"
	UrlParamNotFound     = "URL_PARAM_NOT_FOUND"
	InvalidFileExtension = "INVALID_FILE_EXTENSION"

	DbError              = "DB_ERROR"
	AlreadyExist         = "DATA ALREADY_EXIST"
	Duplication          = "DUPLICATION"
	NotFound             = "NOT_FOUND"
	UniqueFieldViolation = "UNIQUE_FIELD_VIOLATION"

	Bug = "BUG_IN_BACK_END. REPORT_BUG"
)
