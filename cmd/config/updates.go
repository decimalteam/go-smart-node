package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type UpdatesInfoStruct struct {
	filename  string
	LastBlock int64            `json:"last_update"` // last height of 'software_upgrade'
	AllBlocks map[string]int64 `json:"all_updates"` // map of executed upgrades. key - plan name, value - height
}

func NewUpdatesInfo(planfile string) *UpdatesInfoStruct {
	return &UpdatesInfoStruct{
		filename:  planfile,
		LastBlock: 0,
		AllBlocks: make(map[string]int64),
	}
}

func (plan *UpdatesInfoStruct) PushNewPlanHeight(planHeight int64) {
	plan.LastBlock = planHeight
}

func (plan *UpdatesInfoStruct) AddExecutedPlan(planName string, planHeight int64) {
	plan.AllBlocks[planName] = planHeight
}

func (plan *UpdatesInfoStruct) Save() error {
	f, err := os.OpenFile(plan.filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return plan.save(f)
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

func (plan *UpdatesInfoStruct) Load() error {
	if !fileExist(plan.filename) {
		err := ioutil.WriteFile(plan.filename, []byte("{}"), 0644)
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

func (plan *UpdatesInfoStruct) load(rdr io.Reader) error {
	bytes, err := ioutil.ReadAll(rdr)
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
