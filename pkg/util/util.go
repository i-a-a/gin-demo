package util

import (
	"io"
	"math/big"
	"net"
	"reflect"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Ip2Int(ip string) int {
	i := big.NewInt(0).SetBytes(net.ParseIP(ip).To4()).Int64()
	return int(i)
}

func Int2Ip(ip int) string {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// 日志自动分割
func NewWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    50,
		MaxAge:     1,
		MaxBackups: 2,
		LocalTime:  true,
		Compress:   false, // 不压缩 （自己清理）
	}
}

// 结构体指针转map，尽量不使用该方法
// @param2 是否忽略空指针字段
func StructToMap(structPtr interface{}, ignoreNil bool) map[string]interface{} {
	obj := reflect.ValueOf(structPtr)
	// 不是结构体指针直接返回
	if obj.Kind() != reflect.Ptr || obj.Elem().Kind() != reflect.Struct {
		return nil
	}
	v := obj.Elem()
	t := v.Type()
	n := v.NumField()

	myMap := make(map[string]interface{})
	for i := 0; i < n; i++ {
		itemVal := v.Field(i)
		isPtr := itemVal.Kind() == reflect.Ptr
		// 忽略空指针
		if ignoreNil && isPtr && itemVal.IsNil() {
			continue
		}
		// 需要有json标签
		key := t.Field(i).Tag.Get("json")
		if key == "" || key == "-" {
			continue
		}

		myMap[key] = itemVal.Interface()
	}
	return myMap
}
