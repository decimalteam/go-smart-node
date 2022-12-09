package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromOldConfig(t *testing.T) {
	oldConfig := []byte(`{"last_update":8123401,
	"all_updates":{"https://repo.decimalchain.com/7010801":7010801,
	"https://repo.decimalchain.com/7098701":7098701,
	"https://repo.decimalchain.com/7348701":7348701,
	"https://repo.decimalchain.com/7519401":7519401,
	"https://repo.decimalchain.com/7944001":7944001,
	"https://repo.decimalchain.com/7980901":7980901,
	"https://repo.decimalchain.com/8037701":8037701,
	"https://repo.decimalchain.com/8123401":8123401}}`)
	updInf := NewUpdatesInfo("")
	r := bytes.NewReader(oldConfig)
	updInf.load(r)
	assert.Equal(t, 8, len(updInf.AllBlocks), "old AllBlock must be 8")
	assert.Equal(t, int64(8123401), updInf.LastBlock, "LastBlock must be in safe")
}

func TestSaveLoad(t *testing.T) {
	updInf := NewUpdatesInfo("")
	updInf.PushNewPlanHeight(1)
	updInf.AddExecutedPlan("1", 1)
	updInf.AddExecutedPlan("2", 1)
	updInf.AddExecutedPlan("3", 1)
	tmp := make([]byte, 0)
	buf := bytes.NewBuffer(tmp)
	updInf.save(buf)
	newinf := NewUpdatesInfo("")
	newinf.load(buf)
	assert.Equal(t, updInf.AllBlocks, newinf.AllBlocks, "AllBlocks must be same")
	assert.Equal(t, updInf.LastBlock, newinf.LastBlock, "LastBlock must be same")
}
