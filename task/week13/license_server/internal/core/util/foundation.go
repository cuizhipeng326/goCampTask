package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

type LineScannerFunc func(line string, scanner *bufio.Scanner)
type Empty struct{}

func GetCurrentPath() string {
	d, err := os.Getwd()
	if err != nil {
		return ""
	}
	return d
}

func GetExePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

func ScanFile(path string, line LineScannerFunc) (*bufio.Scanner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		line(s.Text(), s)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

func SubString(s string, beginIndex int, length int) string {
	b := int(math.Max(0, float64(beginIndex)))
	l := int(math.Min(float64(len(s)), float64(length+beginIndex)))
	return string([]rune(s)[b:l])
}

func Assert(s string) {
	log.Fatal("Ring: ", s)
}

func AssertRingError(clz string, s string) {
	log.Fatal("Ring: [Start] ERROR " + clz + " " + s)
}

func AssertRingError2(clz string, s string, err error) {
	log.Println(err)
	AssertRingError(clz, s)
}

//脱变量
func RegexVarBraces(s string) []string {
	var reg = regexp.MustCompile(`\$\{(.*?)}`)
	p := reg.FindStringSubmatch(s)
	if len(p) == 2 {
		return p
	}
	return nil
}

//脱括号
func RegexBrackets(s string) []string {
	var reg = regexp.MustCompile(`\((.*?)\)`)
	p := reg.FindStringSubmatch(s)
	if len(p) == 2 {
		return p
	}
	return nil
}
func RegexSquareBrackets(s string) []string {
	var reg = regexp.MustCompile(`\[([^\[\]]*)\]`)
	p := reg.FindStringSubmatch(s)
	if len(p) == 2 {
		return p
	}
	return nil
}

//压缩子串
func CompressString(str string) string {
	if str == "" {
		return ""
	}
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

func GoFmtFile(file string) (string, error) {
	cmd := exec.Command("go", "fmt", file)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

func CopyMapInterface(s map[string]interface{}) map[string]interface{} {
	targetMap := make(map[string]interface{})
	for k, v := range s {
		targetMap[k] = v
	}
	return targetMap
}

func AppendMapInterface(source, append map[string]interface{}) {
	for k, v := range append {
		source[k] = v
	}
}

func AppendMapString(source, append map[string]string) {
	for k, v := range append {
		source[k] = v
	}
}

func AppendMapStrings(source, append map[string][]string) {
	for k, v := range append {
		source[k] = v
	}
}
func FmtUrl(url string) string {
	return strings.ReplaceAll(url, "\\", Separator)
}

func ListContainsPrefix(ss []string, s string) bool {
	for _, a := range ss {
		if strings.HasPrefix(s, a) {
			return true
		}
	}
	return false
}

func ListContainsString(ss []string, s string) bool {
	for _, a := range ss {
		if s == a {
			return true
		}
	}
	return false
}

func FormatJsonByte(byte []byte) []byte {
	var str bytes.Buffer
	_ = json.Indent(&str, byte, "", "    ")
	return str.Bytes()
}

func StringArrayInsert(slice *[]string, str string, index int) {
	*slice = append(*slice, "0")
	copy((*slice)[index+1:], (*slice)[index:])
	(*slice)[index] = str
}

func StringSliceRemove(slice *[]string, i int) {
	(*slice)[i] = (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
}

func RemoveSliceLast(slice *[]string) {
	*slice = (*slice)[:len(*slice)-1]
}

func UUID(ex bool) string {
	// 创建
	u := uuid.NewV4()
	uuid := u.String()
	if ex {
		uuid = strings.ReplaceAll(uuid, "-", "")
	}
	return uuid
}

func HostJoin(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func FormatJson(string string) string {
	var str bytes.Buffer
	_ = json.Indent(&str, []byte(string), "", "    ")
	s := str.String()
	if len(s) == 0 {
		s = string
	}
	return s
}

func JoinHost(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func GetBoolString(bol bool) string {
	if bol {
		return "true"
	} else {
		return "false"
	}
}

func StringLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

type PermutationEach func(crr []string) bool

func Permutation(n int, arr []string, placeholder string, each PermutationEach) bool {
	if len(arr) == 0 {
		return false
	}
	for i := n; i < len(arr); i++ {
		var crr []string
		for _, v := range arr {
			crr = append(crr, v)
		}
		crr[i] = placeholder
		if each(crr) {
			return true
		}
		Permutation(i+1, crr, placeholder, each)
	}
	return false
}

//格式化时间，标准
func FormatLocalTimeStandard(formatTime string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00:00", formatTime+":00")
	return t.Local().Format("2006-01-02 15:04:05"), err
}
