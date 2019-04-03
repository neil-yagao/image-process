package controllers

import (
	"encoding/json"
	"image-process/models"
	"image-process/service"
	"log"

	"github.com/astaxie/beego"
)

// ImageController about Image
type ImageController struct {
	beego.Controller
}

// UploadFile handler
// @Description upload image
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /upload/:type/:userId [post]
func (u *ImageController) UploadFile() {
	imageService := new(service.DefaultImageService)
	imageInfo := new(models.Image)
	imageInfo.Relatedid = u.Ctx.Input.Param(":userId")
	imageInfo.Type = u.Ctx.Input.Param(":type")
	f, h, err := u.GetFile("file")
	if err != nil {
		log.Fatal("getfile err ", err)
	}
	imageInfo.Name = h.Filename
	defer f.Close()
	beego.Debug("h", h.Filename)
	// json.Unmarshal(u.Ctx.Input.RequestBody, user)
	u.Data["json"], _ = imageService.UploadFile(f, imageInfo)
	//imageService.ProcessImage(imageInfo)
	u.ServeJSON()
}

// ProcessImage handle
// @Description
// @router /process/:id [post]
func (u *ImageController) ProcessImage() {
	imageService := new(service.DefaultImageService)
	configs := new([]models.ImageProcessConfig)
	imageID := u.Ctx.Input.Param(":id")
	json.Unmarshal(u.Ctx.Input.RequestBody, configs)
	beego.Debug("config", configs)
	imageService.ProcessImageWithConfig(imageID, *configs)
	u.Data["json"] = map[string]bool{"success": true}
	u.ServeJSON()
}

// HeartBeat testing controler is working
// @router / [get]
func (u *ImageController) HeartBeat() {
	u.Data["json"] = map[string]bool{"success": true}
	u.ServeJSON()
}
