package cerror

func DefaultErrorMsg(httpCode int) string {
	switch httpCode {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Access Denied"
	case 404:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}
