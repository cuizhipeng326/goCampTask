// 作业题：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
package main

import (
	"errors"
	"fmt"
)

var ErrNoRows = errors.New("no rows matched")

// DaoMocker mock process db request
func DaoMocker(req interface{}) error {
	// do some other process

	// dao层作为中间件会被很多其他业务处理逻辑调用，不应该封装太多东西，直接把原始的错误信息抛到上层
	return daoProcess(req)
}

// daoProcess mock process a specific db operation
func daoProcess(req interface{}) error {
	return errors.New("no rows matched")
}

// BizMocker mock business who need get some resource from db
func BizMocker() error {
	var req interface{}
	if err := bizProcess1(req); errors.Is(err, ErrNoRows) {
		// 业务最上层的接口中把打印相关日志和err中的其他信息或堆栈跟踪信息，方便根据定位err出现的根因
		fmt.Printf("process error:%+v\n", err)
		return err
	}

	return nil
}

func bizProcess1(req interface{}) error {
	// do some other process

	return bizProcess2(req) // 业务中间层不处理错误就只把错误往上层传递
}

func bizProcess2(req interface{}) error {
	// do some other process

	if err := DaoMocker(req); err != nil {
		if errors.Is(err, ErrNoRows) {
			// 业务中直接调用dao层的接口负责把相关定位问题用的信息或堆栈跟踪信息封装到err中并抛给上层
			return fmt.Errorf("no record found in db: %v", err)
		}

		return fmt.Errorf("db process failed, %v", err)
	}

	return nil
}

func main() {
	BizMocker()
}
