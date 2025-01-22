package internal

import (
	"encoding/json"
	"errors"
	"net"
	"unicode/utf8"
)

func getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("getLocalIp err")
}

func UnsafeMergeMap(tgt map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	if len(src) == 0 {
		return tgt
	}
	if tgt == nil {
		tgt = map[string]interface{}{}
	}
	for k, v := range src {
		if _, has := tgt[k]; !has {
			tgt[k] = v
		}
	}
	return tgt
}

func SubStrDecodeRuneInString(s string, length int) string {
	var size, n int
	for i := 0; i < length && n < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[n:])
		n += size
	}
	return s[:n]
}

// 将结构体转成map
func Struct2Map(obj interface{}) (map[string]interface{}, error) {
	arr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	dst := map[string]interface{}{}
	err = json.Unmarshal(arr, &dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
