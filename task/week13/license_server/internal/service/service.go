package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"license_kratos/internal/core/util"
	"license_kratos/internal/license"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	pb "license_kratos/api"
	"license_kratos/internal/dao"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.LicenseServer), new(*Service)))

// Service service.
type Service struct {
	ac      *paladin.Map
	dao     dao.Dao
	license *license.License
}

// New new a service and return.
func New(d dao.Dao, license *license.License) (s *Service, cf func(), err error) {
	s = &Service{
		ac:      &paladin.TOML{},
		dao:     d,
		license: license,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s", req.Name)
	return
}

// SayHelloURL bm demo func.
func (s *Service) SayHelloURL(ctx context.Context, req *pb.HelloReq) (reply *pb.HelloResp, err error) {
	reply = &pb.HelloResp{
		Content: "hello " + req.Name,
	}
	fmt.Printf("hello url %s", req.Name)
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}

func (s *Service) GetHardwareCode(ctx context.Context, req *empty.Empty) (resp *pb.HardwareCodeResp, err error) {
	return &pb.HardwareCodeResp{
		HardwareCode: s.license.GetMachineCode(),
	}, nil
}

func (s *Service) GetLicenseAll(ctx context.Context, req *empty.Empty) (resp *pb.LicenseInfo, err error) {
	info, err := s.license.GetLicenseInfo()
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "GetLicenseAll Failed"), log.KV("error", err))
		err = ecode.Errorf(-1, "获取授权信息失败 error(%v)", err)
		return
	}
	resp = &info
	return
}

func (s *Service) GetLicenseProject(ctx context.Context, req *pb.ProjectLicenseReq) (resp *pb.ProjectLicense, err error) {
	licenseInfo, err := s.license.GetProjectLicenseInfo(req.ProjectLicense)
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "GetProjectLicenseInfo Failed"), log.KV("error", err))
		err = ecode.Errorf(-1, "获取项目授权失败 error(%v)", err)
		return
	}
	resp = &licenseInfo
	return
}

func (s *Service) PostLicenseFile(ctx context.Context, req *empty.Empty) (resp *empty.Empty, err error) {
	c := ctx.(*bm.Context)
	buf, err := util.FileHeaderToBuffer(c.Request.MultipartForm.File["file"][0])
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "FileHeaderToBuffer Failed"), log.KV("error", err))
		err = ecode.Errorf(-1, "上传文件失败 error(%v)", err)
		return
	}
	// 文件解密
	err = s.license.UpdateLicFile(string(buf.Bytes()))
	if err != nil {
		log.Errorv(context.Background(), log.KV("log", "UpdateLicFile Failed"), log.KV("error", err))
		err = ecode.Errorf(-1, "上传文件失败 error(%v)", err)
		return
	}
	return
}

func (s *Service) PostLicenseSync(ctx context.Context, req *empty.Empty) (resp *empty.Empty, err error) {
	err = s.license.LicenseSync()
	if err != nil {
		err = ecode.Errorf(-1, "同步授权失败 error(%v)", err)
		return
	}
	return
}

func (s *Service) GetPermissionProject(ctx context.Context, req *pb.ProjectLicenseReq) (resp *pb.PermissionProject, err error) {
	permissionProject, err := s.license.GetPermissionProject(req.ProjectLicense)
	if err != nil {
		err = ecode.Errorf(-1, "获取权限失败 error(%v)", err)
		return
	}
	resp = &pb.PermissionProject{
		Content: permissionProject,
	}
	return
}
