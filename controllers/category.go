package controllers

import (
	"beeblog/models"
	"fmt"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	fmt.Println("category get")
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["IsCategory"] = true
	op := c.Input().Get("op")
	fmt.Println(op)
	switch op {
	case "add":
		{
			fmt.Println("add")
			name := c.Input().Get("cname")
			if len(name) == 0 {
				break
			}

			err := models.AddCategory(name)
			if err != nil {
				beego.Error(err)
			}
			c.Redirect("/category", 301)
		}
	case "del":
		{
			fmt.Println("del")
			id := c.Input().Get("id")
			if len(id) == 0 {
				break
			}
			err := models.DeleteCategory(id)
			if err != nil {
				beego.Error(err)
			}
			c.Redirect("/category", 301)
		}
	}
	fmt.Println("none")
	c.TplName = "Category.html"

	var err error
	c.Data["Categories"], err = models.GetAllCategories()

	if err != nil {
		beego.Error(err)
	}
}
