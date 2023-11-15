package features

import "kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common"

// Feature v2ray 组件接口，
// 组件全部实现 Feature 接口，方便进行 Start() Stop() 等生命周期的统一管理，
// 而且后面拓展新组件不需要，修改生命周期代码
type Feature interface {
	common.HasType
	common.Runnable
}
