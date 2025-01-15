package models

import (
    "github.com/beego/beego/v2/client/orm"
)

type PropertyDescription struct {
    DestID      string `orm:"pk;column(dest_id)"`
    Images      string `orm:"type(text)"`
    Description string `orm:"type(text)"`
    Rating      float64
    Review      string `orm:"type(text)"`
}

func init() {
    orm.RegisterModel(new(PropertyDescription))
}