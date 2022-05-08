package license

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

var _machineCode string

func (l *License) GetMachineCode() string {
	//return "2B26-1D53-1C27-DC4D-8743-0850-A0E9-7889"
	if len(_machineCode) < 1 {
		_machineCode = l.getMachineCode()
		hash := md5.Sum([]byte(_machineCode))
		_machineCode = hex.EncodeToString(hash[:])
		_machineCode = strings.ToUpper(_machineCode)

		var newC []byte
		length := len(_machineCode)
		for i := 0; i < length; i += 4 {
			j := i + 4
			if length < i+4 {
				j = length - i
			}
			newC = append(newC, _machineCode[i:j]...)
			newC = append(newC, '-')
		}
		newC = newC[:len(newC)-1]
		_machineCode = string(newC)
	}

	return _machineCode
}

func (l *License) getMachineCode() string {
	var machineCode string
	id, err := l.GetCpuID()
	if err != nil {
		fmt.Println("license-server", "GetCpuID error", err, nil)
	} else {
		machineCode += id
	}
	boardSN, err := l.GetBoardSerialNumber()
	if err != nil {
		fmt.Println("license-server", "GetBoardSerialNumber error", err, nil)
	} else {
		machineCode += boardSN
	}
	diskSN, err := l.GetDiskSerialNumber()
	if err != nil {
		fmt.Println("license-server", "GetDiskSerialNumber error", err, nil)
	} else {
		machineCode += diskSN
	}
	return machineCode
}
