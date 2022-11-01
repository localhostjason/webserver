package uv

import (
	"webserver/server/util/ue"
)

const (
	E_DB_ERROR      = "E_DB_ERROR"
	E_INVLIAD_PARAM = "E_INVALID_PARAM"
)

var eMap = map[string]ue.Error{
	E_DB_ERROR:      {Code: E_DB_ERROR, Desc: "DB操作错误", Msg: "%v"},
	E_INVLIAD_PARAM: {Code: E_INVLIAD_PARAM, Desc: "无效的参数", Msg: "%v"},
}

func init() {
	ue.RegErrors(eMap)
}
