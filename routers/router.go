// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"image-process/controllers"
	"os"
	"time"

	"github.com/astaxie/beego"
)

// HeartBeatController is used for detecting heart beating
type HeartBeatController struct {
	beego.Controller
}

//Get Reply the timestamp of the heart beat
func (hb *HeartBeatController) Get() {
	hb.Data["json"] = map[string]interface{}{
		"success": true, "timestamp": time.Now()}
	hb.ServeJSON()
}

func init() {
	beego.Router("/image/upload/:type/:userId", &controllers.ImageController{}, "post:UploadFile")
	beego.Router("/image/process/:id", &controllers.ImageController{}, "post:ProcessImage")
	beego.Router("/image", &controllers.ImageController{}, "get:HeartBeat")
	beego.Router("/heart-beat", &HeartBeatController{})
	beego.Debug("file loc", os.Getenv("FILE_LOC"))
	beego.SetStaticPath("/static", os.Getenv("FILE_LOC"))
}
