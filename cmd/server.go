package cmd

import (
	"github.com/Bnei-Baruch/wfdb/api"
	"github.com/Bnei-Baruch/wfdb/common"
)

func Init() {
	a := api.App{}
	a.Initialize(common.ACC_URL, common.SKIP_AUTH)
	a.InitDB()
	a.Run(common.LISTEN_ADDRESS)
}
