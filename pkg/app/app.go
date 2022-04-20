package app

import (
	"embed"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/fs"
	"net/http"
	"os"
	"sync"
)

var (
	app  *Application
	once sync.Once
)

type Application struct {
	Router        *gin.Engine
	EmbedAsset    embed.FS
	EmbedTemplate embed.FS
	C             *gin.Context
	DB            *gorm.DB
	//Cache          cache.Store
	Logger         *zap.Logger
	TraceId        string
	host           string
	fullUrl        string
	requestBody    []byte
	Path           Path
	PathSeparator  string
	LineSeparator  string
	RuntimeLogFile string
	AccessLogFile  string
}

type Path struct {
	Root      string
	App       string
	Bootstrap string
	Config    string
	Route     string
	Public    string
	Static    string
	Runtime   string
	View      string
	Log       string
}

type Option func(app *Application)

func WithAsset(asset embed.FS) Option {
	return func(app *Application) {
		app.EmbedAsset = asset
	}
}

func WithView(view embed.FS) Option {
	return func(app *Application) {
		app.EmbedTemplate = view
	}
}

//New 实例化 Application
func New(options ...Option) *Application {
	once.Do(func() {
		app = &Application{
			TraceId:       "", //TODO
		}

		rootPath, _ := os.Getwd()
		buildPath(rootPath)
	})

	for _, option := range options {
		option(app)
	}

	return app
}

//GetInstance 返回 Application
func GetInstance() *Application {
	return app
}

//EmbedAsset 返回嵌入的静态资源文件
func EmbedAsset() http.FileSystem {
	stripped, err := fs.Sub(app.EmbedAsset, "public")
	if err != nil {
		panic(err)
	}
	return http.FS(stripped)
}

//EmbedTemplate 返回嵌入的模板文件
func EmbedTemplate() embed.FS {
	return app.EmbedTemplate
}

