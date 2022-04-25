# gin-api
gin-api 是一款基于 `gin` 开发的 golang api 框架，其已涵盖日常接口开发时所需要的基础功能。


## 启动项目
第一步：下载
```go
$ git clone git@github.com:Node1650665999/go-api.git 
```

第二步：生成 .env 文件，并设置相关参数，例如数据库账号密码，redis账号密码等：
```go
$ cp .env.example .env

$ cat .env
APP_NAME=go-api
APP_ENV=local
APP_KEY=base64:QDvuvqT2HD2s6CEXvXe/gDbv3iGkjluwQJIUdXkf8Dg=
APP_DEBUG=true
APP_URL=http://localhost
APP_PORT=8089

#日志相关
LOG_SIZE=5
EXCEPT_URI=/asset/
LOG_TYPE=daily

DB_DRIVER=mysql
DB_USER=root
DB_PASSWORD="YjERlZxXhyB$#Cjf"
DB_HOST=8.129.104.24
DB_PORT=3506
DB_DATABASE=hzg_union

REDIS_HOST=192.168.1.1
REDIS_PORT=1234
REDIS_AUTH=123456

EMAIL_NOTICE=false
EMAIL_TO=to
EMAIL_PORT=456
EMAIL_USER=user
EMAIL_PASSWORD=password
EMAIL_SSl=true
EMAIL_FROM=from
```
第三步：启动项目
```go
$ go run index.go
```

查看帮助信息：
```go
$ go run index.go -h
Default will run "serve" command, you can use "-h" flag to see all subcommands

Usage:
[command]

Available Commands:
completion  Generate the autocompletion script for the specified shell
help        Help about any command
serve       Start web server

Flags:
-e, --env string   load .env file, example: --env=testing will use .env.testing file
-h, --help         help for this command

Use " [command] --help" for more information about a command.

```

# 使用说明
go-api 包含如下功能：
- 配置文件
- 路由
- 中间件
- 限流
- 认证
- 接口请求
- 接口响应
- 静态资源打包
- 日志处理
- 数据库
- 容器化部署

下面我们分别介绍相关功能的使用。

## 配置文件

项目中的所有配置文件位于根目录下的 `config` 目录中，目前内置了系统相关的一些基本配置文件，用户如有需要，可自行扩充。
``` css
config
├── app.go       //系统初始化相关的配置
├── cache.go     //缓存相关的配置
├── database.go  //数据库相关的配置
├── email.go     //邮件服务器相关的配置
├── log.go       //日志相关的配置
└── redis.go     //redis相关的配置 
```

## 路由
路由文件位于根目录下的 `route` 目录中，目前按照业务划分已内置了四类路由文件：
``` css
route
├── admin.go    //管理后台相关的
├── api.go      //api接口相关
├── static.go   //静态资源相关
└── web.go      //web相关
```

路由非常简单的以 `r.Group(prefix)` 设置前缀来区分的各类接口。

例如 api 接口定义如下：
```go
func RegisterApiRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	
	api.Any("/foo", func(ctx *gin.Context) {
		response.Json(errcode.Success, "", nil)
		return
	})
	...

	return r
}
```

## 中间件
中间件类型分为`全局中间件`和`局部中间件`。 中间件所在的位置位于 `application/middleware` 中。

我们以编写一个支持跨域的全局中间件为例：
```go
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//允许客户端传递校验信息比如 cookie
			//c.Header("Access-Control-Allow-Credentials", "true")
		}

		//options 请求直接返回
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
```

当编写好中间件后，接下来就是中间件的注册了，对于全局中间件其注册位置在 `bootstrap/route.go` 中：
```go
...
//registerMiddleware 注册全局中间件
func registerMiddleware(router *gin.Engine) {
	router.Use(middleware.MountApp())
	router.Use(middleware.Catch())
	router.Use(middleware.Cors())
	router.Use(middleware.AccessLog())
	router.Use(middleware.Translations())
}
...
```

而对于应用在路由组或者个接口中的中间件，一般在路由所在的文件 `route/xxx.go` 中注册就行了。 

在这里我们给 api 路由组添加一个局部中间件：
```go
func RegisterApiRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	
	api.Use(middleware.LimitRouteAndIp("2-M"))
	{
		...
	}

	return r
}
```

## 限流
限流中间件 `limiter.go` 提供三类限流手段:
- 全局限流器 : 限制单个用户访问能访问系统几次。
- 单个路由组/接口限流 : 即所有人访问路由组/接口总的次数。
- 对某个ip访问某接口进行限流 ： 即每人限制访问该接口几次。

此外还支持`单机版限流` 和 `分布式限流` 的切换。

限流中间件中这几个函数如下：
```go
//LimitIp 全局限流器(限制单个用户访问能访问系统几次)
func LimitIp(format string, driver ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + format
		if err := Limiter(driver...).Check(key, format); err != nil {
			response.JsonAbort(errcode.TooManyRequests, err.Error(),nil)
			return
		}
		c.Next()
	}
}

//LimitRoute 针对某个接口限流(即所有人限制访问该接口总共几次)
func LimitRoute(format string, driver ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key    := c.FullPath()
		if err := Limiter(driver...).Check(key, format); err != nil {
			response.JsonAbort(errcode.TooManyRequests, err.Error(), nil)
			return
		}
		c.Next()
	}
}

//LimitRouteAndIp 对某个ip访问某接口进行限流(即每人限制访问该接口几次)
func LimitRouteAndIp(format string, driver ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key    := routeToKeyString(c.FullPath() + c.ClientIP())
		if err := Limiter(driver...).Check(key, format); err != nil {
			response.JsonAbort(errcode.TooManyRequests, err.Error(), nil)
			return
		}
		c.Next()
	}
}
...
```

`format` 参数也就是限流的格式，支持S(秒)、M(分钟)、H(小时)、D(天)这4种：
- 5 reqs/second: "5-S"   即每秒最多有 5 个请求
- 10 reqs/minute: "10-M" 即每分钟最多有 10 个请求
- 1000 reqs/hour: "1000-H" 即每小时最多 1000 个请求
- 2000 reqs/day: "2000-D"  即每天最多 2000 个请求

`driver` 参数则用来在`单机版限流` 和 `分布式限流` 之间切换：
- driver=1 单机版限流，参数缺省时，默认 driver 为 1。
- driver=2 分布式限流

需要注意：
> 当你选择了分布式限流时，如果系统判断 redis 不可用，则自动退化为单机版限流。

最后我们看下限流中间件的使用：
```go
api := r.Group("/api")

api.Use(middleware.LimitRouteAndIp("500-H"))
{
    api.Any("/foo", middleware.LimitRoute("2-S"),func(ctx *gin.Context) {
        response.Json(errcode.Success, "", nil)
        return
    })

    v1 := api.Group("/v1").Use(middleware.LimitRoute("10-M"))
    {
        v1.Any("/info", controller.NewUserController().Info)
    }
}
...
```

## 认证
系统自带了基于 JWT 的认证组件，并且在路由文件 `route/api.go` 中给出了关于`生成token`,`刷新token`的示例：
```go
//token 相关
tokenGroup := api.Group("/token").Use(middleware.LimitRouteAndIp("2-M"))
{
    //生成token
    tokenGroup.Any("/get", func(ctx *gin.Context) {
        user := gin.H{"name": "tcl", "age": 30}
        token := jwt.GenerateToken(user, time.Hour)
        response.Json(errcode.Success, "success", token)
        return
    })

    //刷新token
    tokenGroup.Any("/fresh", middleware.JwtAuth(), func(ctx *gin.Context) {
        token, err := jwt.RefreshToken(ctx.Request.FormValue("token"), time.Hour)
        if err != nil {
            response.Json(errcode.Fail, err.Error(), nil)
        } else {
            response.Json(errcode.Success, "success", token)
        }
        return
    })
}
```

而`token认证`一般我们放在中间件中：
```go
v1 := api.Group("/v1").Use(middleware.LimitRoute("2-M"))
{
    v1.Any("/info", middleware.JwtAuth(), controller.NewUserController().Info)
}
```

## 接口请求
一个接口的请求一般由两部分来构成: `参数验证` 和 `参数获取`。

参数验证 ：
> gin 是基于`github.com/go-playground` 来实现的。 而框架本身基于 go-playground 已封装好一个基础的验证工具 `BaseVld.go`， 以便用户很方便的对参数验证和绑定。

参数获取：
> gin 通过参数绑定就能实现参数获取，相关的教程件见 [这里](https://laravelacademy.org/post/21885) 。


我们以新建一个控制器 `application/http/controller/UserController.go` 为例：
```go
package controller

import (
	"gin-api/application/errcode"
	"gin-api/application/http/logic"
	"gin-api/application/http/validate"
	"gin-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

type Person struct {
	Name string `form:"name" json:"name" binding:"required,oneof=red green"`
	Email   string `form:"email" json:"email" binding:"required,email"`
	Address string `form:"address" json:"address" binding:"required"`
	Age     int    `form:"age" json:"age" binding:"required,gt=0,lt=120"`
}

func (u *UserController) Info(c *gin.Context) {
	var person Person
	valid, err := validate.BindAndValid(c, &person)
	if ! valid {
		response.Json(errcode.Fail, err.First(), nil)
		return
	}

	info := logic.NewUserLogic().Info()

	response.Json(
		errcode.Success,
		"success",
		gin.H{
			"p":     "",
			"m":     c.Request.Method,
			"phone": info.Phone,
			"html":  "<h1>world</h1>",
		},
	)
	return
}
```

1. 首先验证时创建好要绑定的 Struct，并且在 Struct tag 的`form` 和 `binding` 字段设置好映射的字段和验证规则。
2. 接着调用框架本身已封装好的验证函数 `validate.BindAndValid()` 就行了。
3. 最周如果验证失败，验证函数会返回对应的错误信息，如果验证通过，则相应的参数值会挂在 Struct 对象上，交由后续的逻辑处理就行了。

关于 `go-playground` 的更多见 https://github.com/go-playground/validator。

## 接口响应

框架在启动时，通过中间件 `mount.go` 将 `gin.Context` 挂载到了框架的容器 `app` 上。

```go
//MountApp 挂载 gin.Context 到 app 上
func MountApp(c *gin.Context) {
	scheme := "http://"
	if c.Request.TLS != nil {
		scheme = "https://"
	}

	app.C 		    = c
	app.host        = scheme   + c.Request.Host
	app.fullUrl     = app.host + c.Request.RequestURI
	body,_          := c.GetRawData()
	app.requestBody = body
	c.Request.Body  = ioutil.NopCloser(bytes.NewBuffer(body))
}
```

这样一来，在后续需要依赖 gin.Context 的地方，就不需要显式的层层传递了。

框架本身封装了一个非常简洁的 `response` 包，用以屏蔽 gin.Context。
```go
var ctx *gin.Context

//Json response json
func Json(code int, msg string, data interface{}) {
	if msg == "" {
		msg = errcode.CodeText(code)
	}
	ctx = app.GetInstance().C
	ctx.PureJSON(errcode.HttpCode(code), gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return
}

//JsonAbort response json and abort current request
func JsonAbort(code int, msg string, data interface{})  {
	Json(code, msg, data)
	ctx.Abort()
	return
}
```

那么对于响应客户端的请求也就很简单了：
```go
api.Any("/foo", func(ctx *gin.Context) {
    response.Json(errcode.Success, "", nil)
    return
})
```

此外框架还内置了一些约定好的 error code 用来定义特定的响应，
这些code的定义位于`application/errcode/code.go` 中。


```go
const (
	Success         = 200
	Fail            = 400
	No              = 404
	Unauthorized    = 401
	TooManyRequests = 429
	Fatal           = 500
)

var textMap = map[int]string{
	Success:         "成功",
	Fail:            "失败",
	No:              "路由不存在",
	Unauthorized:    "认证失败",
	TooManyRequests: "请求太频繁",
	Fatal:           "系统异常",
}

var httpMap = map[int]int{
	Success:         http.StatusOK,
	Fail:            http.StatusOK,
	No:              http.StatusOK,
	Fatal:           http.StatusInternalServerError,
}

//CodeText 获取错误码描述
func CodeText(code int) string {
	return textMap[code]
}

//HttpCode 获取错误码映射的http状态码
func HttpCode(code int) int {
	v, ok := httpMap[code]
	if ok == false {
		return http.StatusOK
	}
	return v
}
```

在响应请求时，如果缺省 `msg`, 通过`CodeText()`来获取对应的提示信息，通过`HttpCode()` 来获取传递给底层 gin 的 http code。


## 静态资源打包
框架在入口文件 `index.go` 中通过 embed 的方式来打包好了静态资源。
>需注意， `//go:embed` 的方式打包静态资源要求go的版本不得小于 1.16。
```go
package main

//go:embed public/*
var Assets embed.FS

//go:embed application/http/view/*
var Views embed.FS
```

所谓的静态资源通常分为两种：一种是`图片、css、js` 等资源类文件，另一类则是需要解析的 `.tmpl` 模板。

对于这两类文件，在具体的路由文件我们都做了处理：
```go
router := gin.Default()

//======================解析嵌入的静态资源========================
router.StaticFS("/public", Asset())

//=====================解析嵌入的html模板=========================
//通过 pattern 可以支持从多级目录中解析模板
templ := template.Must(template.New("").ParseFS(tmpl, "templates/*.tmpl", "templates/*/*.tmpl"))
router.SetHTMLTemplate(templ)
```


对于资源目录中的图片 "public/static/tyh.jpg"，在浏览器中这样访问：
```shell 
curl 127.0.0.1/public/static/tyh.jpg
```

对于.tmpl 模板这样渲染 ：
```go
router.GET("/index", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "title": "Embed Demo",
    })
})
router.GET("/home", func(c *gin.Context) {
	//注意，所有的模板子目录都相对于 `application/http/view`
    c.HTML(http.StatusOK, "home/index.tmpl", gin.H{
        "title": "hello",
    })
})
```

## 日志
系统针对不同的场景封装了三个日志函数，位于`pkg/logger` 包中。
```go
//RuntimeLog 记录错误日志
func RuntimeLog(text string) {
	filename := app.RuntimeLogFile()
	GetLogger(filename).Error(text)
}

//AccessLog  记录访问日志
func AccessLog(name string, text interface{}) {
	var txt string
	if str , ok := text.(string); ok {
		txt = str
	} else {
		txt = fmt.Sprintf("%+v",text)
	}

	filename := app.AccessLogFile()
	if filename == "" {
		filename = getFilename()
	}

	if len(txt) > 0 {
		GetLogger(filename).Info(name, zap.String("info", txt))
	} else {
		GetLogger(filename).Info(name)
	}
}

//Log 记录普通日志(可自定义存放日志的目录)
func Log(name string, text interface{}, logFile... string) {
	filename := getFilename(logFile...)
	GetLogger(filename).Sugar().Info(name + " : ", text)
}
```

RuntimeLog() :
> 用来记录 panic 引起的错误日志，该日志会写入 `runtime/logs/{date}.log` 以天进行滚动的日志文件中。该日志函数由系统自身调用，用户无需关注。

AccessLog()  :
> 用来记录接口访问日志，会将接口路径映射生成对应的日志目录，并在其中写入访问日志，该函数由中间件自动调用，用户无需关注。

Log() :
> 用来给用户使用的，用户可以自定义日志存储目录，默认情况下会以调用Log()方法所在的包为路径生成对应的目录，并在其中写入访问日志。 

## 数据库
系统在 GORM 封装了一个查询构造器 `application/http/model/Builder.go` ，其包含一系列辅助函数用来快速进行 CRUD 等操作。

下面给出该查询构造器的一些使用案例。

1. 首先新增模型 `application/http/model/UserModel.go`
```go
type User struct {
	ID            uint      `json:"id"`
	Username      string    `json:"username"` // 用户名
	Phone         string    `json:"phone"`    // 手机号码
}

//TableName 重写表名
func (* User) TableName() string {
	return "union_users"
}

func NewUser() *User {
	u  := &User{}
	return u
}
```

> 因为我们的数据表名不总是按照 GORM 的规则来映射 Struct 名称，因此更多时候显示定义 `TableName()` 来告诉 GORM 映射的表名。 

2. 模型创建好了后，接下来利用该模型来进行增删查改。

```go
func TestInsert(t *testing.T) {
	
	//也支持map格式插入数据
	/*user := map[string]interface{}{
		"username" : "tcl",
		"phone" : "1388888888",
	}*/

	user := User{
		Username: "tcl",
		Phone:    "1388888888",
	}

	rows := Create(&user)
	fmt.Printf("%+v\n", user)
	fmt.Println(user.ID)
	fmt.Println(rows)
}

func TestBatchInsert(t *testing.T) {
	users := []User{
		{
			Username: "tcl1",
			Phone:    "13135677653",
		},
		{
			Username: "tcl2",
			Phone:    "13135677655",
		},
	}

	rows := CreateBatch(&users, len(users)/2)
	fmt.Printf("%+v\n", users)
	for _, user := range users {
		fmt.Println(user.ID)
	}
	fmt.Println(rows)
}

func TestUpdate(t *testing.T) {
	data := User{
		Username: "tcl",
		Phone:    "13877777777",
	}
	where := "id in (3,4) and is_del=1"
	rows := Updates(data, where)
	fmt.Println(rows)
}

type Extract struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

//TestTake 获取单行数据
func TestTake(t *testing.T) {
	var user User
	var extract Extract

	where := ""
	Take(&user, where, "id desc")

	//提取字段
	helpers.Extract(user, &extract)
	fmt.Printf("%+v\n", user)
	fmt.Printf("%+v\n", extract)
	fmt.Println(helpers.JsonEncode(extract))
}

//TestFind 获取多行数据
func TestFind(t *testing.T) {
	var users []User
	var extracts []Extract

	where := "age=20"
	Find(&users, where, "id desc")

	//提取字段
	helpers.Extract(users, &extracts)
	fmt.Printf("%+v\n", users)
	fmt.Printf("%+v\n", extracts)
	fmt.Println(helpers.JsonEncode(extracts))
}

//TestFindPage 分页获取数据
func TestFindPage(t *testing.T) {
	var users []User

	where := "id > 6"
	page := 2     //当前页
	pageSize := 2 //每页数量

	rows := FindPage(&users, where, "id desc", page, pageSize)

	fmt.Printf("%d\n", rows)
	fmt.Printf("%+v\n", users)

	data := map[string]interface{}{
		"list":     users,
		"paginate": GetPaginate(),
	}
	fmt.Println(helpers.JsonEncode(data))
}

//TestColumn 获取单列
func TestColumn(t *testing.T) {
	where := "age > 20"
	phones := []string{}
	Column(&phones, "phone", where)
	fmt.Println(phones)
}

//TestDelete 删除
func TestDelete(t *testing.T) {
	where := "id = 27"
	rows := Delete(&User{}, where)
	fmt.Println(rows)
}

//TestTansAction 事务处理
func TestTansAction(t *testing.T) {
    user := User{
        Username: "tcl",
        Phone:    "1388888888",
    }
    
    //开启事务
    TransStart()
    
    rows := Create(&user)
    if rows == 0 {
        //回滚事务
        TransRollback()
        return
    }
    
    where := fmt.Sprintf("id=%d", user.ID)
    if Delete(user, where) == 0 {
        TransRollback()
        return
    }
    
    //事务提交
    TransCommit()
    fmt.Println("success")
}
```