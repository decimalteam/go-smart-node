package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type UpdatesInfoStruct struct {
	LastBlock int64            `json:"last_update"` // last height of 'software_upgrade'
	AllBlocks map[string]int64 `json:"all_updates"` // map of executed upgrades. key - plan name, value - height
	filename  string
}

func NewUpdatesInfo(planfile string) *UpdatesInfoStruct {
	return &UpdatesInfoStruct{
		LastBlock: 0,
		AllBlocks: make(map[string]int64),
		filename:  planfile,
	}
}

func (plan *UpdatesInfoStruct) PushNewPlanHeight(planHeight int64) {
	plan.LastBlock = planHeight
}

func (plan *UpdatesInfoStruct) AddExecutedPlan(planName string, planHeight int64) {
	plan.AllBlocks[planName] = planHeight
}

func (plan *UpdatesInfoStruct) Load() error {
	if !fileExist(plan.filename) {
		err := os.WriteFile(plan.filename, []byte("{}"), 0644)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(plan.filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return plan.load(f)
}

func (plan *UpdatesInfoStruct) Save() error {
	// This file is diagnostic only. Never expose a partially overwritten JSON
	// document to readers, and never use it as proof of a committed migration.
	bz, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(plan.filename), ".updates-*")
	if err != nil {
		return err
	}
	temporary := f.Name()
	defer os.Remove(temporary)
	if _, err = f.Write(bz); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = os.Rename(temporary, plan.filename); err != nil {
		return err
	}
	// Directory sync makes the rename durable on supported Unix filesystems.
	// Windows does not support syncing directories through os.File.Sync.
	if runtime.GOOS != "windows" {
		dir, openErr := os.Open(filepath.Dir(plan.filename))
		if openErr != nil {
			return openErr
		}
		defer dir.Close()
		return dir.Sync()
	}
	return nil
}

func (plan *UpdatesInfoStruct) save(wrt io.Writer) error {
	bytes, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	_, err = wrt.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (plan *UpdatesInfoStruct) load(rdr io.Reader) error {
	bytes, err := io.ReadAll(rdr)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, plan)
	if err != nil {
		return err
	}

	return nil
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
