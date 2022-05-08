// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"github.com/google/wire"
	"license_kratos/internal/dao"
	"license_kratos/internal/license"
	"license_kratos/internal/server/http"
	"license_kratos/internal/service"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, license.Provider, service.Provider, http.New, NewApp))
}
