package repo

import (
	"context"
	v1 "file-transfer/pkg/api/v1"
	"file-transfer/pkg/config"
	"file-transfer/pkg/db/dbmongo"
	"path/filepath"
	"testing"
)

func Test_query_tag(t *testing.T) {
	configFile, err := filepath.Abs("../../../_output/file-transfer.yaml")
	if err != nil {
		t.Logf("file not found: %s", configFile)
	}
	ctx := context.Background()
	config.ReadConfig(configFile)
	client := dbmongo.GetClient(ctx)
	defer dbmongo.CloseClient(ctx)

	msgRepo := newMessageRepo(client)
	result, err := msgRepo.Query(ctx, &v1.MessageQuery{
		UserId:   "A",
		PageNum:  1,
		PageSize: 2,
	})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	t.Logf("result: %v", result)
}
