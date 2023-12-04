package error

func GetErrorMessage(input error) string {
	if input != nil {
		return input.Error()
	}
	return ""
}
