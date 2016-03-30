package models

import (
	"fmt"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	_DB_NAME        = "db/beego.db"
	_SQLITE3_DRIVER = "sqlite3"
)

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}
func AddComment(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	comment := new(Comment)
	comment.Tid = tidNum
	comment.Name = nickname
	comment.Content = content
	comment.Created = time.Now()

	o := orm.NewOrm()

	_, err = o.Insert(comment)
	if err != nil {
		return err
	}

	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("Id", tid).One(topic)
	if err != nil {
		return err
	}
	topic.ReplyCount++
	topic.ReplyTime = time.Now()

	_, err = o.Update(topic)
	if err != nil {
		return err
	}

	return nil
}
func AddTopic(title, category, content string) error {
	o := orm.NewOrm()

	topic := new(Topic)
	topic.Title = title
	topic.Category = category
	topic.Content = content
	topic.Created = time.Now()
	topic.Updated = time.Now()
	topic.ReplyTime = time.Now()

	_, err := o.Insert(topic)
	if err != nil {
		return err
	}

	if category == "" {
		return nil
	}

	cate := new(Category)
	cate.Title = category

	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	cate.TopicCount++

	if err == nil {
		_, err = o.Update(cate)
		if err != nil {
			return err
		}
	} else {
		cate.Created = time.Now()
		cate.TopicTime = time.Now()
		_, err = o.Insert(cate)
		if err != nil {
			return err
		}
	}

	return nil
}
func AddCategory(cname string) error {
	//fmt.Println(cname)
	o := orm.NewOrm()
	cate := new(Category)
	cate.Title = cname
	cate.Created = time.Now()
	cate.TopicTime = time.Now()
	cate.TopicCount = 1

	qs := o.QueryTable("category")
	err := qs.Filter("title", cname).One(cate)
	if err == nil {
		return err
	}

	_, err = o.Insert(cate)
	if err != nil {
		return err
	}

	return nil
}
func DeleteComment(tid, rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	comment := new(Comment)
	comment.Id = ridNum

	o := orm.NewOrm()
	_, err = o.Delete(comment)
	if err != nil {
		return err
	}

	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("Id", tid).One(topic)
	if err != nil {
		return err
	}
	topic.ReplyCount--

	_, err = o.Update(topic)
	if err != nil {
		return err
	}

	return nil
}
func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	_, err = o.Delete(topic)

	return err
}
func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	fmt.Println("DeleteCategory")
	fmt.Println(id)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)

	return err
}
func GetAllTopics(category string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()

	topics := make([]*Topic, 0)

	qs := o.QueryTable("topic")
	var err error
	if len(category) != 0 {
		if isDesc {
			_, err = qs.Filter("Category", category).OrderBy("-created").All(&topics)
		} else {
			_, err = qs.Filter("Category", category).All(&topics)
		}
	} else {
		if isDesc {
			_, err = qs.OrderBy("-created").All(&topics)
		} else {
			_, err = qs.All(&topics)
		}
	}

	return topics, err
}

func GetComments(tid string) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	comments := make([]*Comment, 0)

	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&comments)

	return comments, err
}
func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.All(&cates)

	return cates, err
}
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()

	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)

	return topic, err
}
func ModifyTopic(tid, title, category, content string) error {
	tidnum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()

	topic := &Topic{Id: tidnum}

	if o.Read(topic) == nil {
		topic.Title = title
		topic.Category = category
		topic.Content = content
		topic.Updated = time.Now()
		o.Update(topic)
	}
	return nil
}
