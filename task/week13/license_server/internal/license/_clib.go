
package license

// #cgo !windows LDFLAGS: -L ${SRCDIR}/../../../../../../cpp_store/cbb/license_client/1.0.0/lib/gcc_centos64/ -llicense_client64.1.0.0
// #cgo !windows LDFLAGS: -L ${SRCDIR}/../../../../../../cpp_store/third/openssl/1.1.0f/lib/gcc_centos64/ -lssl.1.1.0f -lcrypto.1.1.0f
// #cgo !windows CFLAGS: -D_LINUX_
// #cgo LDFLAGS: -lstdc++ -ldl -lc
// #cgo CFLAGS: -I ${SRCDIR}/../../../../../../cpp_store/cbb/license_client/1.0.0/include/ -I ${SRCDIR}/../../../../../../cpp_store/third/openssl/1.1.0f/include/
// #include <stdlib.h>
// #include "crypt.h"
import "C"
import "unsafe"

func GetMachine() string {
	var bytes [C.MACHINE_CODE_LENGTH]byte
	//defer C.free(unsafe.Pointer(&bytes[0]))
	_ = C.license_client_machine_code((*C.uchar)(unsafe.Pointer(&bytes[0])))
	return string(bytes[:])
}