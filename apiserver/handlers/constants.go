package handlers

const (
	headerContentType = "Content-Type"
)

const (
	apiRoot         = "/v1/"
	apiSummary      = apiRoot + "summary"
	apiUsers        = apiRoot + "users"
	apiSessions     = apiRoot + "sessions"
	apiSessionsMine = apiSessions + "/mine"
	apiUsersMe      = apiUsers + "/me"
)

const (
	charsetUTF8         = "charset=utf-8"
	contentTypeJSON     = "application/json"
	contentTypeJSONUTF8 = contentTypeJSON + "; " + charsetUTF8
)
