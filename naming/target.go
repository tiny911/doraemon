package naming

import (
	"github.com/tiny911/doraemon"
)

// DC idc标识
type DC string

const (
	// 线下idc
	OfflineDC DC = "idc-offline"

	// OnlineDC 线上DC idc
	OnlineDC DC = "idc-online"

	// 线下target
	OfflineTg string = "127.0.0.1:8560"

	// OnlineTg 线上 target
	OnlineTg string = "127.0.0.1:8560"
)

// DataCenter 依据env获取datacenter
func DataCenter(env doraemon.Env) DC {
	if env == doraemon.Online || env == doraemon.OnlinePre {
		return OnlineDC
	}

	return OfflineDC
}

// Target 依据env获取target地址
func Target(env doraemon.Env) string {
	if env == doraemon.Online { // notice: 预览环境的地址从vedTarget中获取
		return OnlineTg
	}

	return OfflineTg
}
