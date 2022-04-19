package helpers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"runtime"
)

//InArray 判断切片是否包含某个元素
func InArray(need interface{}, haystack interface{}) bool {
	switch key := need.(type) {
	case int:
		for _, item := range haystack.([]int) {
			if item == key {
				return true
			}
		}
	case string:
		for _, item := range haystack.([]string) {
			if item == key {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if item == key {
				return true
			}
		}
	case float64:
		for _, item := range haystack.([]float64) {
			if item == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}

//ByteFormat 将字节格式化为指定的单位
// refer:https://blog.microdba.com/golang/2016/05/01/golang-byte-conv/
func ByteFormat(size uint64) string {
	sz   := float64(size)
	base := float64(1024)
	unit := []string{"B","KB","MB","GB","TB","EB"}
	i := 0

	for sz >= base {
		sz /= base
		i++
	}
	return fmt.Sprintf("%.2f%s",sz, unit[i])
}


//StrUuid 生成字符串格式的UUID
func StrUuid(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

//ExecCommand 运行系统命令和二进制文件
func ExecCommand(cmd string, params ...string) (string, error) {
	// Print Go Version
	cmdOutput, err := exec.Command(cmd, params...).Output()
	if err != nil {
		return "", err
	}
	return string(cmdOutput), err
}

//JsonEncode 实现Json编码
func JsonEncode(v interface{}) string {
	u,_ := json.Marshal(v)
	return string(u)
}

//JsonDecode 实现json解码
func JsonDecode(data string, v interface{}) error  {
	return json.Unmarshal([]byte(data), v)
}

//Extract 提取json编码数据
func Extract(src interface{}, dst interface{}) error {
	return JsonDecode(JsonEncode(src), dst)
}

// Empty 类似于 PHP 的 empty() 函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

//CurrentFuncName 返回当前函数名
func CurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}