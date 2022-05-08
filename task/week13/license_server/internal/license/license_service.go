package license

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/google/wire"
	"github.com/wenzhenxi/gorsa"
	"io/ioutil"
	pb "license_kratos/api"
	"license_kratos/internal/core/util"
	"license_kratos/internal/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Provider = wire.NewSet(NewLicense)

const _licensePath = "./license.lic"
const _productIdPath = "./productId"
const _permissionPath = "./permission.json"

type License struct {
	licContext        pb.LicenseInfo                  // 授权信息
	permissionFilter  map[string]model.PermissionList // 过滤后的权限信息
	permissionAllTree map[string][]model.Permission   // 全部权限信息 key:parentId
	permissionAll     map[string]model.Permission     // 全部权限信息 key:resourceId
	sensitive         []model.Permission              // 敏感字段
	roles             []model.Role                    // roles
	permissionGroups  []model.PermissionGroup         // permissionGroups
	product           string                          // 产品ID
	cfg               ClientHttpConfig
}

type ClientHttpConfig struct {
	Ip     string
	Port   int
	Url    string
	Appid  string
	Secret string
}

func NewLicense() (license *License, err error) {
	var (
		cfg ClientHttpConfig
		ct  paladin.Map
	)
	if err = paladin.Get("sync_server.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("ClientHttp").UnmarshalTOML(&cfg); err != nil {
		return
	}

	license = &License{
		cfg:               cfg,
		permissionAll:     make(map[string]model.Permission),
		permissionAllTree: make(map[string][]model.Permission),
		permissionFilter:  make(map[string]model.PermissionList),
	}
	hardwareCode := license.GetMachineCode()
	log.Infov(context.Background(), log.KV("log", "机器码："+hardwareCode))

	err = license.loadProductFile()
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "加载product文件失败"), log.KV("error", err))
		return
	}

	err = license.loadPermission()
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "加载permission文件失败"), log.KV("error", err))
		return
	}

	err2 := license.loadLicFile()
	if err2 != nil {
		log.Errorv(context.Background(), log.KV("log", "加载lic文件失败"), log.KV("error", err))
		return
	}
	return
}

// 更新授权文件
func (l *License) UpdateLicFile(content string) (err error) {
	err = l.decodeContent(content)
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "加载lic文件失败"), log.KV("error", err))
		return
	}
	// 保存到文件
	err = util.WriteFile(_licensePath, []byte(content), true)
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "UpdateLicFile failed"), log.KV("error", err))
		return
	}
	// 更新权限
	l.filterPermission()
	return
}

// 从远程服务同步授权文件
func (l *License) LicenseSync() error {
	log.Infov(context.Background(), log.KV("log", "LicenseSync"))
	httpClient := resty.New()
	host := fmt.Sprintf("http://%s:%d", l.cfg.Ip, l.cfg.Port)
	httpClient.SetHostURL(host)
	httpClient.SetTimeout(3 * time.Second)
	httpReq := httpClient.R()
	httpReq.Method = http.MethodGet
	httpReq.URL = l.cfg.Url
	httpReq.SetQueryParam("hardwareCode", _machineCode).SetQueryParam("productId", l.product)
	httpReq.SetHeader("appid", l.cfg.Appid).SetHeader("secret", l.cfg.Secret)
	response, err := httpReq.Send()
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "HttpLicenseClient Get 失败"), log.KV("error", err))
		return errors.New("同步授权文件失败")
	}
	if response != nil && response.StatusCode() != 200 {
		log.Errorv(context.Background(), log.KV("log", "HttpLicenseClient Get 失败"+string(response.Body())), log.KV("error", err))
		return errors.New("同步授权文件失败")
	}
	//解密
	if err := l.decodeContent(string(response.Body())); err != nil {
		return err
	}
	//保存文件
	err = util.WriteFile(_licensePath, response.Body(), true)
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "create file failed"), log.KV("error", err))
		return err
	}

	// 更新权限
	l.filterPermission()
	return nil
}

// 获取授权信息
func (l *License) GetLicenseInfo() (license pb.LicenseInfo, err error) {
	if len(l.licContext.HardwareCode) == 0 {
		err = errors.New("没有授权信息")
		return
	}
	license = l.checkValidTime(l.licContext)
	return
}

// 获取项目授权信息
func (l *License) GetProjectLicenseInfo(projectLicense string) (projectLicenseInfoRet pb.ProjectLicense, err error) {
	if len(l.licContext.HardwareCode) == 0 {
		err = errors.New("没有授权信息")
		return
	}
	validTime := l.checkValidTime(l.licContext)
	for _, projectLicenseInfo := range validTime.ProjectLicenses {
		if projectLicenseInfo.ProjectLicense == projectLicense {
			projectLicenseInfoRet = *projectLicenseInfo
			return
		}
	}
	err = errors.New("不存在该项目")
	return
}

func (l *License) GetPermissionProject(projectLicense string) (content string, err error) {
	if _, exist := l.permissionFilter[projectLicense]; !exist {
		return "", errors.New("不存在该项目")
	}
	origData, err := json.Marshal(l.permissionFilter[projectLicense])
	if err != nil {
		return
	}
	log.Debugv(context.Background(), log.KV("log", "GetPermissionProject"), log.KV("projectLicense", projectLicense), log.KV("content", string(origData)))
	//加密
	code := "233"
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	hash := md5.Sum([]byte(encoded))
	ss := hex.EncodeToString(hash[:])
	fmt.Println(ss)
	key := hash[:]
	encrypted, _ := AesEncryptCBC(origData, key)
	content = base64.StdEncoding.EncodeToString(encrypted)

	return
}

func testDec(encrypted string) {
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
	fmt.Println("解密结果：", string(decrypted))
	return
}

// 读取lic文件，解密后保存到内存
func (l *License) loadLicFile() error {
	// 读取文件
	file, err := ioutil.ReadFile(_licensePath)
	if err != nil {
		return err
	}
	// 解密
	buf := string(file)
	err = l.decodeContent(buf)
	if err != nil {
		return err
	}
	// 更新权限
	l.filterPermission()
	return nil
}

func (l *License) loadPermission() error {
	// 读取文件
	fileByte, err := ioutil.ReadFile(_permissionPath)
	encoded, _ := base64.StdEncoding.DecodeString(string(fileByte))
	if err != nil {
		return err
	}
	var permissionList model.PermissionList
	err = json.Unmarshal(encoded, &permissionList)
	if err != nil {
		return err
	}
	l.roles = permissionList.Roles
	l.permissionGroups = permissionList.PermissionGroups
	for _, permission := range permissionList.Permissions {
		l.permissionAll[permission.ResourceID] = permission
		l.permissionAllTree[permission.ParentID] = append(l.permissionAllTree[permission.ParentID], permission)
		if permission.ResourceType == model.ResourceTypeSensitive {
			l.sensitive = append(l.sensitive, permission)
		}
	}
	return nil
}

func (l *License) filterPermission() {
	l.permissionFilter = make(map[string]model.PermissionList)
	// 用l.licContext 从全部权限中匹配
	for _, license := range l.licContext.ProjectLicenses {
		permissions := l.sensitive
		for _, module := range license.Modules {
			if module.IsAuth == 1 {
				resourceId := strconv.Itoa(int(module.ModuleId))
				if permission, exist := l.permissionAll[resourceId]; exist {
					// 当前
					permissions = append(permissions, permission)
					// 子
					childPermissions := l.getChildPermissions(permission.ResourceID)
					permissions = append(permissions, childPermissions...)
				}
			}
		}

		l.permissionFilter[license.ProjectLicense] = model.PermissionList{
			Permissions:      permissions,
			Roles:            l.roles,
			PermissionGroups: l.permissionGroups,
		}
	}
}

func (l *License) getChildPermissions(resourceId string) (childPermissions []model.Permission) {
	permissions, exist := l.permissionAllTree[resourceId]
	if !exist {
		return nil
	}
	for _, permission := range permissions {
		childPermissions = append(childPermissions, permission)
		child := l.getChildPermissions(permission.ResourceID)
		childPermissions = append(childPermissions, child...)
	}
	return
}

// 解密后保存到内存
func (l *License) decodeContent(cryptoContent string) error {
	// 解密
	buf := cryptoContent
	// 公钥解密私钥加密
	buf, err := gorsa.PublicDecrypt(buf, Pubkey)
	if err != nil {
		return err
	}
	// 机器码DES
	code := l.GetMachineCode()
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	hash := md5.Sum([]byte(encoded))
	// des 密文数组
	decodeString, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return err
	}
	encryptCBC, _ := AesDecryptCBC(decodeString, hash[:])
	licContext := string(encryptCBC)
	fmt.Println("授权内容", licContext)
	if licContext == "" {
		return errors.New("没有授权信息")
	}
	var licContextTemp pb.LicenseInfo
	err = json.Unmarshal([]byte(licContext), &licContextTemp)
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "文件内容："+licContext), log.KV("error", err))
		return errors.New("非法授权信息")
	}
	if licContextTemp.HardwareCode != l.GetMachineCode() {
		return errors.New("非法授权信息")
	}
	if myProductId, _ := strconv.Atoi(l.product); licContextTemp.ProductId != int32(myProductId) {
		return errors.New("产品不匹配")
	}
	l.licContext = licContextTemp
	return nil
}

func (l *License) loadProductFile() error {
	// 读取文件
	file, err := ioutil.ReadFile(_productIdPath)
	if err != nil {
		return err
	}
	l.product = string(file)
	l.product = strings.TrimFunc(l.product, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	return nil
}

// 时间校验
func (l *License) checkValidTime(licenseInfo pb.LicenseInfo) pb.LicenseInfo {
	var realLicense = licenseInfo
	timeTemplate := "2006-01-02 15:04:05"
	for i, _ := range realLicense.ProjectLicenses {
		if realLicense.ProjectLicenses[i].Permanent != 1 {
			now := time.Now()
			stampAuthBeginTime, _ := time.ParseInLocation(timeTemplate, realLicense.ProjectLicenses[i].AuthBeginTime, time.Local)
			stampAuthExpireTime, _ := time.ParseInLocation(timeTemplate, realLicense.ProjectLicenses[i].AuthExpireTime, time.Local)
			if now.Before(stampAuthBeginTime) || now.After(stampAuthExpireTime) {
				realLicense.ProjectLicenses[i].Status = 2
			}
		}
	}
	return realLicense
}
