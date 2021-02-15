## Go-MySQL-Server-Demo

### **go-mysql-server简介**

**仓库地址**如下：

-   https://github.com/dolthub/go-mysql-server

根据这个库宣称，**go-mysql-server**会100%兼容MySQL语法（go-mysql-server is to support 100% of the statements that MySQL does.）；

目前该库支持的语法见：

-   [SUPPORTED](https://github.com/dolthub/go-mysql-server/blob/master/SUPPORTED.md)

下面是go-mysql-server的简介：

>   **go-mysql-server** is a SQL engine which parses standard SQL (based on MySQL syntax) and executes queries on data sources of your choice. A simple in-memory database and table implementation are provided, and you can query any data source you want by implementing a few interfaces.
>
>   **go-mysql-server** also provides a server implementation compatible with the MySQL wire protocol. That means it is compatible with MySQL ODBC, JDBC, or the default MySQL client shell interface.
>
>   [Dolt](https://www.doltdb.com/), a SQL database with Git-style versioning, is the main database implementation of this package. Check out that project for reference implmenentations.

翻译如下：

go-mysql-server是一个SQL引擎，它解析标准SQL（基于MySQL语法）并对选择的数据源执行查询。

本库提供了一个简单的基于内存的数据库和表实现，但是你可以通过实现几个接口来查询所需的任何数据源（例如，基于文件、网络等）。

本库还提供了一个与MySQL wire协议兼容的服务器实现。所以本库可以和MySQL ODBC、JDBC或默认的MySQL-client-shell接口兼容！

[Dolt](https://www.doltdb.com/), 是一个具有Git风格版本控制的SQL数据库，是这个包的主要数据库实现。可以查看该项目以获取参考实现。

这个库的**详细文档**：

-   [go-mysql-server godoc](https://godoc.org/github.com/dolthub/go-mysql-server)

<br/>

### **一个简单的例子**

go-mysql-server库已经包含了一个SQL引擎和服务器实现。所以，如果要启动服务器，首先需要实例化引擎并传递给你自己的`sql.Database`实现，该引擎将负责处理检索数据的所有逻辑；

下面是一个初始化并插入数据的例子：

```go
package main

import (
	"fmt"
	"time"

	sqle "github.com/dolthub/go-mysql-server"

	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
)

const (
	user = "user"
	passwd = "pass"
	address = "localhost"
	port = "13306"
	dbName    = "test"
	tableName = "my_table"
)

func main() {
	engine := sqle.NewDefault()
	engine.AddDatabase(createTestDatabase())

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%s", address, port),
		Auth:     auth.NewNativeSingle(user, passwd, auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	go func() {
		s.Start()
	}()

	fmt.Println("mysql-server started!")

	<- make(chan interface{})
}

func createTestDatabase() *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()

	rows := []sql.Row{
		sql.NewRow("John Doe", "jasonkay@doe.com", []string{"555-555-555"}, time.Now()),
		sql.NewRow("John Doe", "johnalt@doe.com", []string{}, time.Now()),
		sql.NewRow("Jane Doe", "jane@doe.com", []string{}, time.Now()),
		sql.NewRow("Evil Bob", "jasonkay@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()),
	}

	for _, row := range rows {
		_ = table.Insert(ctx, row)
	}

	return db
}

```

在文件中，我们定义了一些数据库会使用到的常量；

在main中首先，初始化了一个默认的SQL引擎，并调用`AddDatabase`添加了一个数据库；

createTestDatabase首先创建了一个Database，并创建了一个Table（table具有name、email、phone_numbers等信息），最后向database中加入了这个表，并插入了一些数据，最后返回了`*memory.Database`；

代码最后创建了数据库配置，并使用引擎engine和数据库配置config创建了真正的server；

在新的gorontine中启动了server，主协程阻塞等待；

下面我们来测试这个数据库；

既然这个库已经支持了MySQL-Cli，我们就直接使用客户端连接；

首先启动服务器，随后在终端连接：

```bash
$ mysql --host=127.0.0.1 --port=13306 -uuser -ppass test -e "SELECT * FROM my_table"
+----------+--------------------+-------------------------------+----------------------------+
| name     | email              | phone_numbers                 | created_at                 |
+----------+--------------------+-------------------------------+----------------------------+
| John Doe | jasonkay@doe.com   | ["555-555-555"]               | 2021-02-14 11:57:34.785025 |
| John Doe | johnalt@doe.com    | []                            | 2021-02-14 11:57:34.785026 |
| Jane Doe | jane@doe.com       | []                            | 2021-02-14 11:57:34.785026 |
| Evil Bob | jasonkay@gmail.com | ["555-666-555","666-666-666"] | 2021-02-14 11:57:34.785027 |
+----------+--------------------+-------------------------------+----------------------------+
```

mysql命令指定了host、port、username、passwd、数据库名、执行SQL，最后获取到了结果；

数据库端输出日志：

```diff 
$ go run app.go 
mysql-server started!
+ INFO: NewConnection: client 1
+ INFO: ConnectionClosed: client 1
```

除了可以通过`-e`指令直接执行之外，我们甚至可以直接连接数据库：

```bash
$ mysql --host=127.0.0.1 --port=13306 -uuser -ppass test
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Welcome to the MariaDB monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.9-Vitess 

Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

MySQL [test]> SELECT count(name) from my_table;
+----------------------+
| COUNT(my_table.name) |
+----------------------+
|                    4 |
+----------------------+
1 row in set (0.00 sec)

MySQL [test]> SELECT email FROM my_table WHERE name = 'Evil Bob';
+--------------------+
| email              |
+--------------------+
| jasonkay@gmail.com |
+--------------------+
1 row in set (0.00 sec)

MySQL [test]> SELECT name,year(created_at) FROM my_table;
+----------+---------------------------+
| name     | YEAR(my_table.created_at) |
+----------+---------------------------+
| John Doe |                      2021 |
| John Doe |                      2021 |
| Jane Doe |                      2021 |
| Evil Bob |                      2021 |
+----------+---------------------------+
4 rows in set (0.01 sec)
```

可以看到，使用终端我们同样完成了查询！

<br/>

### **自定义数据源&索引实现**

上面使用到的database、table等，其实都是memory包提供的一个数据源的简单实现；

我们也可以手动实现一些接口，来定义自己的数据源；

这里官方文档已经描述的很清楚了，这里不再赘述：

-   https://github.com/dolthub/go-mysql-server#custom-data-source-implementation

除了可以自定义数据源之外，还可以自定义索引：

-   https://github.com/dolthub/go-mysql-server#indexes

<br/>

### **使用go-mysql-server做项目测试**

go-mysql-server除了可以作为一个简单的MySQL数据库的实现之外，也可以用在项目的测试场景！

熟悉Java的同学可能听说过H2数据库，这里go-mysql-server的使用和H2数据库类似；

下面给出一个例子，在这个例子中，我们创建了一个宠物表，并使用[xorm](https://github.com/go-xorm/xorm)完成对宠物的CRUD测试；

pets表结构：

schema.sql

```mysql
USE `test`;

DROP TABLE IF EXISTS `pets`;
CREATE TABLE `pets`
(
    `id`    INT(10) AUTO_INCREMENT NOT NULL COMMENT '宠物编号',
    `name`  VARCHAR(20)            NOT NULL COMMENT '宠物名称',
    `age`   TINYINT(3)             NOT NULL COMMENT '宠物年龄',
    `photo` VARCHAR(30) DEFAULT NULL COMMENT '宠物图片',
    `ctime` DATETIME    DEFAULT CURRENT_TIMESTAMP,
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

首先，我们创建pet基本类型：

models/pet.go

```go
package models

import (
	"time"
)

type (
	Pet struct {
		Id    int       `json:"id" binding:"required" form:"id"`
		Name  string    `json:"name" xorm:"varchar(20)" binding:"required" form:"name"`
		Age   int       `json:"age" binding:"required" form:"age"`
		Photo string    `json:"photo" xorm:"varchar(30)" form:"photo"`
		Ctime time.Time `json:"created_at" xorm:"ctime"`
	}
)
```

以及数据库连接：

models/models.go

```go
package models

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

const (
	Username = ""
	Passwd   = ""
	Host     = ""
	Port     = ""
	Dbname   = ""
)

func InitDb(username, passwd, host, port, dbname string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
			username,
			passwd,
			host,
			port,
			dbname))
	if err != nil {
		return nil, err
	}

	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(15)

	//名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engine.SetTableMapper(core.SnakeMapper{})

	return engine, nil
}
```

随后创建DAO层：

dao/pet.go

```go
package dao

import (
	"fmt"
	"go-mysql-server-demo/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type PetDAO struct {
	DB *xorm.Engine
}

func (p *PetDAO) CreatePet(pet *models.Pet) error {
	insert, err := p.DB.Table("pets").Insert(pet)
	if err != nil {
		return err
	}
	if insert != 1 {
		return fmt.Errorf("error, fail to insert, maybe exist")
	}

	return nil
}

func (p *PetDAO) FindPetById(id int) (*models.Pet, error) {
	var pet models.Pet
	has, err := p.DB.Table("pets").Where("id = ?", id).Get(&pet)
	if err != nil {
		return nil, err
	}
	if !has || pet.Id == 0 {
		return nil, fmt.Errorf("pet not found")
	}

	return &pet, nil
}

func (p *PetDAO) Update(petId, petAge int, petName string) error {
	res, err := p.DB.Exec("UPDATE `pets` SET `name` = ?, `age` = ? WHERE `id` = ?", petName, petAge, petId)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		return fmt.Errorf("fail to update, maybe record not exist")
	}

	return nil
}

func (p *PetDAO) DeleteById(petId int) error {
	res, err := p.DB.Exec("DELETE FROM `pets` WHERE `id` = ?", petId)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		return fmt.Errorf("fail to delete, maybe record not exist")
	}

	return nil
}
```

最后就是DAO层的测试代码了：

dao/pet_test.go

```go
package dao

import (
	"fmt"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"go-mysql-server-demo/models"
)

const (
	user      = "user"
	passwd    = "pass"
	address   = "localhost"
	port      = "13306"
	dbName    = "test"
	tableName = "pets"
)

var petDAO *PetDAO

func TestMain(m *testing.M) {
	db, err := models.InitDb(user, passwd, address, port, dbName)
	if err != nil {
		panic(err)
	}

	go initMySQL()

	petDAO = &PetDAO{DB: db}

	m.Run()
}

func initMySQL() {
	engine := sqle.NewDefault()
	engine.AddDatabase(createTestDatabase())

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%s", address, port),
		Auth:     auth.NewNativeSingle(user, passwd, auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	go s.Start()

	fmt.Println("mysql-server started!")
}

func createTestDatabase() *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "id", Type: sql.Int64, Nullable: false, Source: tableName},
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "age", Type: sql.Int64, Nullable: false, Source: tableName},
		{Name: "photo", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "ctime", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()

	rows := []sql.Row{
		sql.NewRow(1, "cat", 11, "", time.Now()),
		sql.NewRow(2, "dog", 21, "", time.Now()),
		sql.NewRow(3, "mouse", 31, "", time.Now()),
	}

	for _, row := range rows {
		_ = table.Insert(ctx, row)
	}

	return db
}

func TestPetDAO_CreatePet(t *testing.T) {
	err := petDAO.CreatePet(&models.Pet{
		Name:  "tiger",
		Age:   2,
		Photo: "haha.jpg",
		Ctime: time.Now(),
	})
	if err != nil {
		panic(err)
	}
}

func TestPetDAO_FindPetById(t *testing.T) {
	pet, err := petDAO.FindPetById(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(pet)
}

func TestPetDAO_Update(t *testing.T) {
	err := petDAO.Update(1, 99, "mouse")
	if err != nil {
		panic(err)
	}
}

func TestPetDAO_DeleteById(t *testing.T) {
	err := petDAO.DeleteById(1)
	if err != nil {
		panic(err)
	}
}
```

在测试代码中，我们首先在`TestMain`中创建了一个在内存中的数据库，并且创建了数据库的连接；

最后调用`m.Run()`启动了测试；

执行测试，最终输出：

```bash
=== RUN   TestPetDAO_CreatePet
mysql-server started!
[xorm] [info]  2021/02/15 14:47:58.173553 [SQL] INSERT INTO `pets` (`id`,`name`,`age`,`photo`,`ctime`) VALUES (?, ?, ?, ?, ?) []interface {}{0, "tiger", 2, "haha.jpg", "2021-02-15 14:47:58"}
INFO: NewConnection: client 1
--- PASS: TestPetDAO_CreatePet (0.32s)
=== RUN   TestPetDAO_FindPetById
[xorm] [info]  2021/02/15 14:47:58.477560 [SQL] SELECT `id`, `name`, `age`, `photo`, `ctime` FROM `pets` WHERE (id = ?) LIMIT 1 []interface {}{1}
&{1 cat 11  2021-02-15 06:47:58 +0800 CST}
--- PASS: TestPetDAO_FindPetById (0.00s)
=== RUN   TestPetDAO_Update
[xorm] [info]  2021/02/15 14:47:58.478561 [SQL] UPDATE `pets` SET `name` = ?, `age` = ? WHERE `id` = ? []interface {}{"mouse", 99, 1}
--- PASS: TestPetDAO_Update (0.00s)
=== RUN   TestPetDAO_DeleteById
[xorm] [info]  2021/02/15 14:47:58.479562 [SQL] DELETE FROM `pets` WHERE `id` = ? []interface {}{1}
--- PASS: TestPetDAO_DeleteById (0.00s)
PASS
```

可以看到所有的测试都通过了！

>   为了简单起见，这里的测试用例写的都比较简单；
>
>   你也可以编写更加复杂的测试用例；

使用go-mysql-server做测试的好处就是，在测试时，我们就已经自带了一个MySQL环境，所以只需要修改不同环境（如prod、dev、test）下的配置，即可完成DAO（数据存取）层的测试；

同时，这样也可以避免编写大量的打桩代码，而只需在每次测试时初始化数据库即可！

<br/>

## **附录**

**原文地址：**

-   Github Pages：[使用纯Go实现的MySQL数据库](https://jasonkayzk.github.io/2021/02/14/%E4%BD%BF%E7%94%A8%E7%BA%AFGo%E5%AE%9E%E7%8E%B0%E7%9A%84MySQL%E6%95%B0%E6%8D%AE%E5%BA%93/)
-   国内Gitee镜像：[使用纯Go实现的MySQL数据库](https://jasonkay.gitee.io/2021/02/14/%E4%BD%BF%E7%94%A8%E7%BA%AFGo%E5%AE%9E%E7%8E%B0%E7%9A%84MySQL%E6%95%B0%E6%8D%AE%E5%BA%93/)

**go-mysql-server仓库地址：**

-   https://github.com/dolthub/go-mysql-server

<br/>