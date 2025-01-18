package utils

import (
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Tang-RoseChild/mahonia"
)

const (
	Error   = "error"
	Running = "running"
	Stop    = "stop"
)

func RandAuthToken() string {
	buf := make([]byte, 32)
	_, err := crand.Read(buf)
	if err != nil {
		return RandString(64)
	}

	return fmt.Sprintf("%x", buf)
}

// RandString 生成长度为length的随机字符串
func RandString(length int64) string {
	sources := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sourceLength := len(sources)
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(sourceLength)])
	}

	return string(result)
}

// Md5 生成32位MD5摘要
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))

	return hex.EncodeToString(m.Sum(nil))
}

// RandNumber 生成0-max之间随机数
func RandNumber(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(max)
}

// GBK2UTF8 GBK编码转换为UTF8
func GBK2UTF8(s string) (string, bool) {
	dec := mahonia.NewDecoder("gbk")

	return dec.ConvertStringOK(s)
}

// ReplaceStrings 批量替换字符串
func ReplaceStrings(s string, old []string, replace []string) string {
	if s == "" {
		return s
	}
	if len(old) != len(replace) {
		return s
	}

	for i, v := range old {
		s = strings.Replace(s, v, replace[i], 1000)
	}

	return s
}

func InStringSlice(slice []string, element string) bool {
	element = strings.TrimSpace(element)
	for _, v := range slice {
		if strings.TrimSpace(v) == element {
			return true
		}
	}

	return false
}

// EscapeJson 转义json特殊字符
func EscapeJson(s string) string {
	specialChars := []string{"\\", "\b", "\f", "\n", "\r", "\t", "\""}
	replaceChars := []string{"\\\\", "\\b", "\\f", "\\n", "\\r", "\\t", "\\\""}

	return ReplaceStrings(s, specialChars, replaceChars)
}

// FileExist 判断文件是否存在及是否有权限访问
func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}

	return true
}

// GetMondayTimes 获取近14周的开始时间
func GetMondayTimes() []time.Time {
	weekCount := 13
	var d int
	if time.Now().Weekday() == 0 {
		d = 6
	} else {
		d = int(time.Now().Weekday()) - 1
	}

	//本周周一的开始时间
	monday, _ := time.Parse("2006-01-02", time.Unix(time.Now().Unix()-int64(d*3600*24), 0).Format("2006-01-02"))

	times := make([]time.Time, 14)
	//周一之前13周的开始时间
	start := time.Unix(monday.Unix()-int64(3600*24*7*weekCount), 0)
	for i := range times {
		times[i] = time.Unix(start.Unix()+int64(3600*24*7*i), 0)
	}
	return times
}
