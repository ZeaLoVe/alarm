package g

import (
	"log"
	"runtime"
)

//change from version 2.0.2
//添加 im,phone 报警通道
//删除 报警合并部分，统一成以级别为区分的报警方式 0：电话、邮件、IM 1-2：邮件、IM 3-6：仅IM
const (
	VERSION = "sdp_2.0.2_001"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
