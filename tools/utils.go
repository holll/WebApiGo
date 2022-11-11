package tools

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func RepBodyToStr(body io.ReadCloser) string {
	repBody, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(repBody)
}

func RepBodyToByteSlice(body io.ReadCloser) []byte {
	repBody, err := io.ReadAll(body)
	if err != nil {
		return make([]byte, 0)
	}
	return repBody
}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//IsNotExist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func IsInStrSlice(str string, strSlice []string) bool {
	for _, tmpStr := range strSlice {
		if str == tmpStr {
			return true
		}
	}
	return false
}
