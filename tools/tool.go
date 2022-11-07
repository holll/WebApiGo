package tools

import "io"

func RepBodyToStr(body io.ReadCloser) string {
	repBody, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(repBody)
}
