package routing

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"

	c "go-restful-xorm/config"
	"go-restful-xorm/controller"
)

type WebService struct{}

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
