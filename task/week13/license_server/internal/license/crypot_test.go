package license

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/wenzhenxi/gorsa"
	"io/ioutil"
	"log"
	"testing"
)

func TestPermissionDec(t *testing.T) {
	encrypted := "xxxx" // 待加密的数据
	code := "233"
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	hash := md5.Sum([]byte(encoded))
	ss := hex.EncodeToString(hash[:])
	fmt.Println(ss)
	key := hash[:]
	bytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		panic(err)
	}
	decrypted, _ := AesDecryptCBC(bytes, key)
	log.Println("解密结果：", string(decrypted))
	return
}

func TestAesDecryptCBC(t *testing.T) {
	origData := []byte("nihao") // 待加密的数据
	code := "233"
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	hash := md5.Sum([]byte(encoded))
	ss := hex.EncodeToString(hash[:])
	fmt.Println(ss)
	key := hash[:]
	//key = []byte("ABCDEFGHIJKLMNOP")              // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ CBC模式 --------------------")
	encrypted, _ := AesEncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted, _ := AesDecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))
	return
	log.Println("------------------ ECB模式 --------------------")
	encrypted = AesEncryptECB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AesDecryptECB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	log.Println("------------------ CFB模式 --------------------")
	encrypted = AesEncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = AesDecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}

func TestDesCBCEncrypt(t *testing.T) {
	origData := []byte("nihao") // 待加密的数据
	//key := []byte("12345678")              // 加密的密钥
	// 机器码DES
	code := ""
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	hash := md5.Sum([]byte(encoded))
	key := hash[:]
	log.Println("原文：", string(origData))
	log.Println("------------------ CBC模式 --------------------")
	encrypted, err := AesEncryptCBC(origData, key)
	if err != nil {
		panic(err)
	}
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted, err := AesDecryptCBC(encrypted, key)
	if err != nil {
		panic(err)
	}
	log.Println("解密结果：", string(decrypted))
}

func TestRsa1(t *testing.T) {
	// 公钥加密私钥解密
	if err := applyPubEPriD(); err != nil {
		log.Println(err)
	}
	// 公钥解密私钥加密
	if err := applyPriEPubD(); err != nil {
		log.Println(err)
	}
}

// 公钥加密私钥解密
func applyPubEPriD() error {
	pubenctypt, err := gorsa.PublicEncrypt(`hello world`, Pubkey)
	if err != nil {
		return err
	}

	pridecrypt, err := gorsa.PriKeyDecrypt(pubenctypt, Pirvatekey)
	if err != nil {
		return err
	}
	if string(pridecrypt) != `hello world` {
		return errors.New(`解密失败`)
	}
	return nil
}

// 公钥解密私钥加密
func applyPriEPubD() error {
	prienctypt, err := gorsa.PriKeyEncrypt(`nihao`, Pirvatekey)
	if err != nil {
		return err
	}

	pubdecrypt, err := gorsa.PublicDecrypt(prienctypt, Pubkey)
	if err != nil {
		return err
	}
	if string(pubdecrypt) != `nihao` {
		return errors.New(`解密失败`)
	}
	return nil
}

func TestDes(t *testing.T) {
	// 读取文件
	file, err := ioutil.ReadFile("license.lic")
	if err != nil {
		panic(err)
	}
	// 解密
	buf := string(file)
	// 公钥解密私钥加密
	buf, err = gorsa.PublicDecrypt(buf, Pubkey)
	if err != nil {
		panic(err)
	}
	// 机器码DES
	code := GetMachineCode()
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	hash := md5.Sum([]byte(encoded))
	sss := hex.EncodeToString(hash[:])
	fmt.Println(sss)

	decodeString, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		panic(err)
	}
	encryptCBC, _ := AesDecryptCBC(decodeString, hash[:])
	fmt.Println(string(encryptCBC))
}
