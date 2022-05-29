package controllers

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/Koubae/goweb/src/mvc/controllers/routes"

	"github.com/gin-gonic/gin"
)

const (
	yellowNotFill = "\033[33;20m"
	blueNotFill   = "\033[32;20m"
	green         = "\033[97;42m"
	white         = "\033[90;47m"
	yellow        = "\033[90;43m"
	red           = "\033[97;41m"
	blue          = "\033[97;44m"
	magenta       = "\033[97;45m"
	cyan          = "\033[97;46m"
	reset         = "\033[0m"
)

var (
	_, b, _, _   = runtime.Caller(0)
	basepath     = filepath.Dir(b)
	rootPath     = filepath.Join(basepath, "../../../")
	publicAssets = filepath.Join(basepath, "../../../public")
	pathAssets   = filepath.Join(basepath, "../../../public/static")
)

func Routes() *gin.Engine {
	gin.ForceConsoleColor()
	fmt.Println(publicAssets)
	fmt.Println(pathAssets)
	fmt.Println(rootPath)
	router := gin.New()

	router.LoadHTMLGlob(filepath.Join(basepath, "../views/templates/**/*"))
	// Load Assets
	router.Static("/assets", pathAssets)

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		var (
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor  = param.ResetColor()
			errorColor  = red
		)
		return fmt.Sprintf("%s%v%s%s|%1s-%1v|%s%s[%3d]%s%s{ %-7s }%s%s=>%s%#v\n%s%s%s",
			yellowNotFill, param.TimeStamp.Format("2006/01/02-15:04:05"), resetColor,
			blueNotFill, param.Latency,
			param.ClientIP, resetColor,

			statusColor, param.StatusCode, resetColor,

			methodColor, param.Method, resetColor,

			cyan, resetColor, param.Path,

			errorColor, param.ErrorMessage, resetColor,
		)
	}))
	router.Use(gin.Recovery())

	return routes.SetupRouter(router)
}
