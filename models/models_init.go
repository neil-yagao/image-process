package models

import (
	"os"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //using mysql driver
	uuid "github.com/satori/go.uuid"
)

//Image default struct
type Image struct {
	Id          string `orm:"pk"`
	Type        string
	Relatedid   string
	Loc         string
	Name        string
	Extrastring string
	Isactive    bool      `orm:"default(1)"`
	Createdat   time.Time `orm:"auto_now_add;type(date)"`
	Updatedat   time.Time `orm:"auto_now_add;type(date)"`
}

type ImageProcessConfig struct {
	Name  string `json:name`
	Width uint   `json:width`
}

// TableName redefine to match previous
func (i *Image) TableName() string {
	return "Images"
}

var o orm.Ormer

func init() {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)

	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	databseURL := os.Getenv("MYSQL_URL")
	database := os.Getenv("DB_NAME")
	connectString := user + ":" + password + "@tcp(" + databseURL + ")/" + database + "?charset=utf8"
	orm.RegisterDataBase("default", "mysql", connectString)
	orm.RegisterModel(new(Image))
	o = orm.NewOrm()
	o.Using("default")
}

//GenerateNewImage generate new image database record
func GenerateNewImage(i *Image) (string, error) {
	duplicate := new(Image)
	duplicate.Type = i.Type
	duplicate.Relatedid = i.Relatedid
	o.Read(duplicate)
	if len(duplicate.Id) == 0 {
		i.Id = uuid.Must(uuid.NewV4()).String()
		_, err := o.Insert(i)
		if err != nil {
			beego.Error("err", err)
			return "", err
		}
	} else {
		_, err := o.Update(i)
		if err != nil {
			beego.Error("err", err)
		}
	}

	beego.Debug("id", i.Id)
	return i.Id, nil
}

func FindImageById(i string) *Image {
	image := new(Image)
	image.Id = i
	o.Read(image)
	return image
}
