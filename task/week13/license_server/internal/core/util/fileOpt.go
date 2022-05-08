package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)


type LineFunc func(inde int, line string, lines *[]string) bool
type LineWillInsertFunc func(index int, line string, lines *[]string) (int, bool)

func GetClassAbsolutePath() string {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return f
}

func ParentDirectory(directory string) string {
	s := strings.ReplaceAll(directory, WindowsSeparator, Separator)
	return SubString(s, 0, strings.LastIndex(s, Separator))
}

func ReadJsonFile(filePath string, object interface{}) error {
	plan, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(plan, object)
	if err != nil {
		return err
	}
	return nil
}

func WriteMapToJsonFile(filePath string, iMap interface{}, format bool) error {
	jsonString, err := json.Marshal(iMap)
	if err != nil {
		return err
	}
	err = WriteFile(filePath, FormatJsonByte(jsonString), true)
	if err != nil {
		return err
	}
	return nil
}

func AppendFile(filePath, content string) (int, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	len, err := file.WriteString(content)
	if err != nil {
		return 0, err
	}
	return len, nil
}

func RemoveAllDirFiles(dirPath string) error {
	if !FileExist(dirPath) {
		return nil
	}
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, d := range dir {
		err = os.RemoveAll(filepath.Join(dirPath, d.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func ArrayInsertLine(content, before string, lines *[]string, line LineFunc, willInsert LineWillInsertFunc) {
	ret := false
	for i, v := range *lines {
		r := line(i, strings.TrimSpace(v), lines)
		if r {
			ret = r
		}
		if ret && strings.TrimSpace(v) == before {
			is := i
			bol := true
			if willInsert != nil {
				is, bol = willInsert(is, v, lines)
			}
			if bol {
				StringArrayInsert(lines, content, is)
				ret = false
			}
		}
	}
}

func FileInsertLine(filePath, content, before string, line LineFunc, willInsert LineWillInsertFunc) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), NewlineSymbol)
	ArrayInsertLine(content, before, &lines, line, willInsert)
	output := strings.Join(lines, NewlineSymbol)
	err = ioutil.WriteFile(filePath, []byte(output), 0666)
	return err
}

func InsertFile(filePath, content string, offset int64, whence int) error {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	} else {
		n, _ := f.Seek(offset, whence)
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return nil
}

func WriteFile(filePath string, byte []byte, overwrite bool) error {
	if !overwrite && FileExist(filePath) {
		return errors.New("-1000, file has exist")
	}
	dir := filepath.Dir(filePath)
	if !FileExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return ioutil.WriteFile(filePath, byte, 0666)
}

func CreateDir(dirPath string, overwrite bool) error {
	if !overwrite && FileExist(dirPath) {
		return errors.New("-1000, dir did exist")
	}
	var err error
	if !FileExist(dirPath) {
		err = os.MkdirAll(dirPath, os.ModePerm)
	}
	return err
}

func FileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 判断所给路径是否为文件夹
func IsDir(path string) (bool, error) {
	s, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

//
//func FileMimeType(filePath string) (string, error)  {
//	f, err := os.Open(filePath)
//	if err != nil {
//		panic(err)
//	}
//	defer f.Close()
//	// Get the content
//	contentType, err := GetFileContentType(f)
//	if err != nil {
//		return "", err
//	}
//	return contentType, nil
//}

func GetFileContentType(out os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func FileHeaderToBuffer(fileHeader *multipart.FileHeader) (*bytes.Buffer, error) {
	formFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, formFile)
	if err != nil {
		return nil, err
	}
	err = formFile.Close()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadFileToJson(path string, object interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, object)
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func CreateDirs(rootPath string, packages map[string][]string, remove, overwrite bool) error {
	if remove {
		err := os.RemoveAll(rootPath)
		if err != nil {
			return err
		}
	}
	for _, v := range packages {
		fp := ""
		if filepath.IsAbs(v[0]) {
			fp = v[0]
		} else {
			fp = filepath.Join(rootPath, v[0])
		}
		err := CreateDir(fp, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyProperties(source, target interface{}) error {
	pr, err := json.Marshal(source)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(pr, target); err != nil {
		return err
	}
	return nil
}
