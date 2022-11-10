package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

//MkDir 创建一个目录
func MkDir(path string) error {
	_, err := os.Stat(path)

	// 权限问题
	if os.IsPermission(err) {
		return fmt.Errorf(" Permission denied src: %s", err)
	}

	// 已存在
	if os.IsExist(err) {
		return nil
	}

	// 创建目录
	return os.MkdirAll(path, os.ModePerm)
}

/*// Ext get the file ext
func Ext(fileName string) string {
	return path.Ext(fileName)
}

// BaseName 获取文件的 basename
func BaseName(filename string) string  {
	return path.Base(filename)
}*/

//FileExist 判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// ReadFile 读取文件并将内容写入到一个 string 切片中的
func ReadFile(file string) ([]string, error) {
	inputFile, err := os.Open(file)
	defer inputFile.Close()

	if err != nil {
		return nil, err
	}

	inputReader := bufio.NewReader(inputFile)
	container := []string{}
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF || readerError != nil {
			break
		}
		container = append(container, inputString)
	}

	return container, nil
}

// WriteFile 往 file中写入 content
func WriteFile(file string, content string) (int, error) {
	outputFile, _ := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)

	n, err := outputWriter.WriteString(content)
	if err != nil {
		return n, err
	}

	err = outputWriter.Flush()
	return n, err
}

// WriteFileSimple 以非缓冲的方式写入文件
func WriteFileSimple(file string, content string) (int, error) {
	outputFile, _ := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer outputFile.Close()
	return outputFile.WriteString(content)
}

// FIlePutContent 将数据存入文件, 	如果目录不存在, 则创建
func FIlePutContent(data []byte, to string) error {
	dir := filepath.Dir(to)
	err := MkDir(dir)
	if err != nil {
		return err
	}

	_, err = WriteFile(to, string(data))
	if err != nil {
		return err
	}

	return nil
}

// ReadFileSimple 读取文件
func ReadFileSimple(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

// CopyFile 将文件 srcName 复制到 dstName
func CopyFile(dstName string, srcName string) (int64, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

//RemoteDownload 实现下载文件到本地,获得网络文件的输入流以及本地文件的输出流 ,然后将输入流读取到输出流中
func RemoteDownload(remote string, local string) error {
	res, err := http.Get(remote)
	if err != nil {
		return fmt.Errorf("A error occurred: %v", err)
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	// 获得文件的writer对象
	file, err := os.Create(local)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	return nil
}
