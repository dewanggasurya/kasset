package kasset

import (
	"git.kanosolution.net/kano/appkit"
	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type Asset struct {
	orm.DataModelBase `json:"-" bson:"_id"`
	ID                string `json:"_id" bson:"_id"`
	Title             string `json:"title"`
	OriginalFileName  string `json:"originalfilename"`
	NewFileName       string `json:"newfilename"`
	URI               string `json:"uri"`
	ContentType       string `json:"contenttype"`
	Size              int    `json:"size"`
	Tags              string `json:"tags"`
}

func (a *Asset) TableName() string {
	return "assets"
}

func (a *Asset) GetID(c dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{a.ID}
}

func (a *Asset) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		a.ID = keys[0].(string)
	}
}

func (a *Asset) PreSave(c dbflex.IConnection) error {
	if a.ID == "" {
		a.ID = appkit.MakeID("", 32)
	}
	return nil
}
