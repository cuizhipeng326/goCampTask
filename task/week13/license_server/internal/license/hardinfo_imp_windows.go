package license

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func (l *License) GetCpuID() (string, error) {
	command := "wmic"
	args := "path Win32_Processor get ProcessorId"

	var argArray []string
	argArray = strings.Split(args, " ")
	cmd := exec.Command(command, argArray...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
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
		if line == "ProcessorId" {
			continue
		}
		return line, nil
	}
	return "", errors.New("not found")
}

func (l *License) GetBoardSerialNumber() (string, error) {
	command := "wmic"
	args := "path Win32_BaseBoard get SerialNumber"

	var argArray []string
	argArray = strings.Split(args, " ")
	cmd := exec.Command(command, argArray...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
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
		if line == "SerialNumber" {
			continue
		}
		return line, nil
	}
	return "", errors.New("not found")
}

func (l *License) GetDiskSerialNumber() (string, error) {
	id, err := getDiskId()
	if err != nil {
		return "", err
	}
	command := "wmic"
	args := fmt.Sprintf("path Win32_DiskDrive get SerialNumber,DeviceID")

	var argArray []string
	argArray = strings.Split(args, " ")
	cmd := exec.Command(command, argArray...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	reader := strings.NewReader(string(output))
	newReader := bufio.NewReader(reader)

	phys := fmt.Sprintf("PHYSICALDRIVE%d", id)
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
		if strings.Contains(line, phys) {
			split := strings.Split(line, " ")
			if len(split) > 1 {
				return split[len(split)-1], nil
			}
		}
	}
	return "", errors.New("not found")
}

func getDiskId() (int, error) {
	command := "wmic"
	args := "path Win32_OperatingSystem get Name"

	var argArray []string
	argArray = strings.Split(args, " ")
	cmd := exec.Command(command, argArray...)
	output, err := cmd.Output()
	if err != nil {
		return -1, err
	}
	reader := strings.NewReader(string(output))
	newReader := bufio.NewReader(reader)
	compile := regexp.MustCompile("Harddisk([0-9]+)")
	for {
		line, err := newReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return -1, err
			}
		}
		if line == "Name" {
			continue
		}
		submatch := compile.FindStringSubmatch(line)
		if len(submatch) != 2 {
			return -1, errors.New("regexp match failed")
		}

		return strconv.Atoi(submatch[1])
	}
	return -1, errors.New("not found")
}
