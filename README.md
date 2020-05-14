A branch to show Golang RESTful with [Gin-Gonic](https://github.com/gin-gonic/gin) & [XORM](https://github.com/go-xorm/xorm) with MySQL.

### Model

本例以Pet模型为例，其模型定义如下：

models/pet.go

```go
package models

import (
	"time"
)

type (
	Pet struct {
		Id int `json:"id" binding:"required" form:"id"`
		Name string `json:"name" xorm:"varchar(20)" binding:"required" form:"name"`
		Age int `json:"age" binding:"required" form:"age"`
		Photo string `json:"photo" xorm:"varchar(30)"`
		Ctime time.Time `json:"created_at" xorm:"ctime"`
		Utime time.Time `json:"updated_at" xorm:"utime"`
	}
)
```

对应的SQL如下

sql/schema.sql

```sql
USE `test`;

DROP TABLE IF EXISTS `pets`;
CREATE TABLE `pets`
(
    `id`    INT(10) AUTO_INCREMENT NOT NULL COMMENT '宠物编号',
    `name`  VARCHAR(20)            NOT NULL COMMENT '宠物名称',
    `age`   TINYINT(3)             NOT NULL COMMENT '宠物年龄',
    `photo` VARCHAR(30) DEFAULT NULL COMMENT '宠物图片',
    `ctime` DATETIME    DEFAULT CURRENT_TIMESTAMP,
    `utime` DATETIME    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB,
  DEFAULT CHARSET = utf8mb4 COMMENT ='宠物表';

INSERT INTO `pets` (ID, NAME, AGE)
VALUES (1, 'cat', '1');
INSERT INTO `pets` (ID, NAME, AGE)
VALUES (2, 'dog', '2');
INSERT INTO `pets` (ID, NAME, AGE)
VALUES (3, 'mouse', '3');
```

### Routing

在routing/web.go的routing方法中定义了路由:

```go
...
type WebService struct{}

func (w *WebService) routing(db *xorm.Engine) {
	petController := controller.PetController{DB: db}

	r := gin.Default()
	v1 := r.Group("/pets")
	v1.GET("/", petController.Index)
	v1.POST("/", petController.Create)
	v1.GET("/:id", petController.Show)
	v1.PUT("/:id", petController.Update)
	v1.DELETE("/:id", petController.DeleteById)

	err := r.Run()
	if err != nil {
		fmt.Printf("fail to start")
	}
}
```

### API

和routing对应的，在controller/pet_controller.go中定义了具体API的实现:

```go
package controller

import (
	"fmt"
	"go-restful-xorm/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type PetController struct {
	DB *xorm.Engine
}

func (p *PetController) Index(c *gin.Context) {
	var pets []models.Pet
	err := p.DB.Table("pets").Find(&pets)
	if err != nil {
		fmt.Printf("%v", err)
	}

	if len(pets) <= 0 {
		c.JSON(404, gin.H{"status": 404, "message": "not found."})
		return
	}
	c.JSON(200, gin.H{"status": 200, "data": pets})
}

func (p *PetController) Create(c *gin.Context) {
	var pet models.Pet

	fmt.Printf("Create id: %s\n", c.Param("id"))

	if err := c.ShouldBind(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Insert data: %v\n", pet)

	insert, err := p.DB.Table("pets").Insert(&pet)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	if insert != 1 {
		c.JSON(404, gin.H{"error": "Fail to insert, maybe exist?"})
	}

	c.JSON(200, gin.H{"status": 200, "message": "pet item created."})
}

func (p *PetController) Show(c *gin.Context) {
	var pet models.Pet
	has, err := p.DB.Table("pets").Where("id = ?", c.Param("id")).Get(&pet)
	if err != nil {
		c.JSON(404, gin.H{"status": 404, "message": "pet select error"})
		return
	}
	if !has {
		c.JSON(404, gin.H{"status": 404, "message": "pet not found"})
		return
	}

	if pet.Id == 0 {
		c.JSON(404, gin.H{"status": 404, "message": "pet not found"})
		return
	}

	c.JSON(200, gin.H{"status": 200, "data": pet})
}

func (p *PetController) Update(c *gin.Context) {
	petId, _ := strconv.Atoi(c.Param("id"))
	petName := c.PostForm("name")
	petAge, _ := strconv.Atoi(c.PostForm("age"))

	fmt.Printf("Update id: %s\n", c.Param("id"))

	res, err := p.DB.Exec("UPDATE `pets` SET `name` = ?, `age` = ? WHERE `id` = ?", petName, petAge, petId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		c.JSON(404, gin.H{"status": 404, "message": "fail to update, maybe record not exist?"})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "pet data updated."})
}

func (p *PetController) DeleteById(c *gin.Context) {
	res, err := p.DB.Exec("DELETE FROM `pets` WHERE `id` = ?", c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		c.JSON(404, gin.H{"status": 404, "message": "fail to delete, maybe record not exist?"})
		return
	}
	c.JSON(200, gin.H{"status": 200, "message": "pet data deleted."})
}
```

ORM采用的是[XORM](https://github.com/go-xorm/xorm)

### Config

在项目的根目录下使用config.json简单定义了数据库的相关配置;

项目启动时通过config/get_config.go加载相关的配置信息，并使用GetConfig方法获取相关配置；

```go
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	appConfig map[string]string
)

func init() {
	err := loadConfig()
	if err != nil {
		fmt.Printf("failed to load config")
	}
}

func loadConfig() error {
	appConfig = map[string]string{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to load config")
	}

	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config")
	}

	return nil
}

func GetConfig(key string) string {
	if value, exist := appConfig[key]; exist {
		return value
	} else {
		return "error"
	}
}
```

### Deployment

在根目录下的server.go作为整个项目的入口;

```go
package main

import (
	"fmt"
	"go-restful-xorm/routing"
)

func main() {
	server := routing.WebService{}
	server.Run()
	fmt.Println("Server started!")
}
```

server.go中调用了routing/web.go中的Run()方法，完成了创建SQL连接Engine，并对8080（默认）端口对应路由进行监听：

```go
......
func (w *WebService) Run() {
	engine, err := xorm.NewEngine(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
			c.GetConfig("dbusername"),
			c.GetConfig("dbpassword"),
			c.GetConfig("dbhost"),
			c.GetConfig("dbport"),
			c.GetConfig("dbname")))
	if err != nil {
		fmt.Printf("fail to connect database")
	}

	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(15)

	//名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engine.SetTableMapper(core.SnakeMapper{})

	w.routing(engine)
}
```

### Test

可以使用Postman或其他测试工具测试；

