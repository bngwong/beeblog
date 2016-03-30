package controllers

import (
	"beeblog/models"
	//"fmt"
	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["IsTopic"] = true
	c.TplName = "Topic.html"

	topics, err := models.GetAllTopics("", false)
	if err != nil {
		beego.Error(err)
	}

	c.Data["Topics"] = topics
}

func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	tid := c.Input().Get("tid")
	category := c.Input().Get("category")

	var err error
	if len(tid) == 0 {
		err = models.AddTopic(title, category, content)
	} else {
		err = models.ModifyTopic(tid, title, category, content)
	}
	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic", 302)
}

func (c *TopicController) Add() {
	c.TplName = "topic_add.html"
}

func (c *TopicController) View() {
	c.TplName = "topic_view.html"

	/*
	 *  NOTE:2016.03.28
	 *  fatal error:Typed map does not support indexing
	 *  c.Ctx.Input.Params is a map pointer
	 *  right access type like this:
	 *  para := c.Ctx.Input.Params()
	 *  topic, err := models.GetTopic(para["0"])
	 *  wrong access type like this:
	 *  topic, err := models.GetTopic(c.Ctx.Input.Params("0"))
	 */
	para := c.Ctx.Input.Params()
	tid := para["0"]
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	c.Data["Topic"] = topic
	c.Data["Tid"] = tid

	comments, err := models.GetComments(tid)
	if err != nil {
		beego.Error(err)
		return
	}

	c.Data["Comments"] = comments
	c.Data["IsLogin"] = checkAccount(c.Ctx)
}

func (c *TopicController) Modify() {
	c.TplName = "topic_modify.html"
	tid := c.Input().Get("tid")

	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
}

func (c *TopicController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	tid := c.Input().Get("tid")

	err := models.DeleteTopic(tid)
	if err != nil {
		beego.Error(err)
	}

	c.Redirect("/topic", 302)
}
