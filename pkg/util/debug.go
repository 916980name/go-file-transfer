package util

import (
	"context"
	"encoding/json"
	"file-transfer/pkg/common"
	"file-transfer/pkg/log"
)

func DebugPrintObj(ctx context.Context, obj interface{}) {
	if common.FLAG_DEBUG {
		s, e := json.Marshal(obj)
		if e != nil {
			return
		}
		log.C(ctx).Debugw(string(s))
	}
	return
}
