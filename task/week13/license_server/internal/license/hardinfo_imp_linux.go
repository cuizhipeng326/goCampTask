package license

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

func (l *License) GetCpuID() (string, error) {
	cmd := exec.Command("dmidecode", "-t", "processor")
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}

	compile := regexp.MustCompile(`ID:(.*)`)
	if compile == nil {
		return "", errors.New("regexp err")
	}
	submatch := compile.FindStringSubmatch(string(buf))
	if len(submatch) != 2 {
		return "", errors.New("regexp match err")
	}
	s := submatch[1]
	space := strings.Replace(s, " ", "", -1)
	fmt.Println(space)
	return space, nil
}

func (l *License) GetBoardSerialNumber() (string, error) {
	cmd := exec.Command("dmidecode", "-t", "baseboard")
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}

	compile := regexp.MustCompile(`Serial Number:(.*)`)
	if compile == nil {
		return "", errors.New("regexp err")
	}
	submatch := compile.FindStringSubmatch(string(buf))
	if len(submatch) != 2 {
		return "", errors.New("regexp match err")
	}
	s := submatch[1]
	space := strings.Replace(s, " ", "", -1)
	fmt.Println(space)
	return space, nil
}

func (l *License) GetDiskSerialNumber() (string, error) {
	diskName, err2 := l.getDiskName()
	if err2 != nil {
		return "", err2
	}
	argName := fmt.Sprintf("--name=%s", diskName)
	cmd := exec.Command("udevadm", "info", " --query=all", argName)
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}

	compile := regexp.MustCompile(`ID_SERIAL=(.*)`)
	if compile == nil {
		return "", errors.New("regexp err")
	}
	submatch := compile.FindStringSubmatch(string(buf))
	if len(submatch) != 2 {
		return "", errors.New("regexp match err")
	}
	s := submatch[1]
	space := strings.Replace(s, " ", "", -1)
	fmt.Println(space)
	return space, nil
}

func (l *License) getDiskName() (string, error) {
	output, err := exec.Command(
		"lsblk",
		"-n",
		"-l",
		"-s",
	).Output()
	if err != nil {
		return "", err
	}
	var diskName string
	reader := strings.NewReader(string(output))
	newReader := bufio.NewReader(reader)
	for {
		line, err := newReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}

		if strings.Index(line, "/") == len(line)-1 {
			break
		}
	}
	for {
		line, err := newReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		if strings.Contains(line, "disk") {
			indexByte := strings.IndexByte(line, ' ')
			diskName = line[0:indexByte]
		}
	}
	if len(diskName) < 1 {
		return "", errors.New("not found disk name")
	}
	return diskName, nil
}
