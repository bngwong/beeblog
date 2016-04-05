package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["IsHome"] = true
	c.TplName = "Home.html"

	cate := c.Input().Get("cate")

	topics, err := models.GetAllTopics(cate, true)
	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"] = topics

	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Categories"] = categories
}
