package uv

import (
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"webserver/util/ue"
)

// PQ parse query
func PQ(c *gin.Context, out ...interface{}) {
	for _, o := range out {
		if err := c.BindQuery(o); err != nil {
			//if err := c.ShouldBindQuery(o); err != nil {
			EP(err)
		}
		Trim(o)
	}
}

// PPID parse pid with custom idKey
func PPID(c *gin.Context, idName string) int {
	var id = c.Param(idName)
	var i, err = strconv.Atoi(id)
	if err != nil {
		EP(err)
	}
	return i
}

// PID parse pid=1 like
func PID(c *gin.Context) int {
	return PPID(c, "id")
}

// PidStrList 获取字符串形式的id列表 ids=a,b,c,d
func PidStrList(c *gin.Context) []string {
	type ids struct {
		IDS string `form:"ids"` //逗号分割
	}
	var rid ids
	if err := c.BindQuery(&rid); err != nil {
		EP(err)
	}
	return strings.Split(rid.IDS, ",")
}

// PidList parse ids=1,2,3
func PidList(c *gin.Context) []int {
	type ids struct {
		IDS string `form:"ids"` //逗号分割
	}
	var rid ids
	if err := c.BindQuery(&rid); err != nil {
		EP(err)
	}

	var idsUint, err = splitIds(rid.IDS)
	if err != nil {
		EP(err)
	}
	return idsUint

}

func splitIds(ids string) ([]int, error) {
	var idsStr = strings.Split(ids, ",")
	var idsUint []int
	for _, s := range idsStr {
		if i, err := strconv.Atoi(s); err != nil {
			return nil, err
		} else {
			idsUint = append(idsUint, i)
		}
	}
	return idsUint, nil
}

func PEIf(ec string, args ...interface{}) {
	for _, a := range args {
		if err, ok := a.(error); ok && err != nil {
			PE(ec, args...)
		}
	}
}

func PE(ec string, args ...interface{}) {
	panic(ue.NewErr(ec, args...))
}

// EP err parameter
func EP(err error) {

	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)

	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		PE("E_INVALID_PARAM", err)
	}

	errMsg := RemoveTopStruct(errs.Translate(Trans))
	for _, v := range errMsg {
		PE("E_INVALID_PARAM", v)
		break
	}

	//PE("E_INVALID_PARAM", err)
	//ve, ok := err.(validator.ValidationErrors)

	//if !ok {
	//	PE("E_INVALID_PARAM", err)
	//	return
	//}

	//for _, f := range ve {
	//	fmt.Println("xxxxxxxxxxx", f.Field(), f)
	//}

}

// PB parse body
func PB(c *gin.Context, o interface{}) {
	if err := c.Bind(o); err != nil {
		EP(err)
	}
	Trim(o)
}

func RC(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
}
