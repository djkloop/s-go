package routers

import (
	"fuck_youku_api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Include(&controllers.UserController{})
	beego.Include(&controllers.VideoController{})
	beego.Include(&controllers.BaseControllers{})
	beego.Include(&controllers.CommentController{})
	beego.Include(&controllers.BarrageController{})
	beego.Include(&controllers.TopController{})
}
