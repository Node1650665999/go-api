package main

import (
	"embed"
	"fmt"
	"gin-api/application/cmd"
	"gin-api/bootstrap"
	_ "gin-api/config"
	"gin-api/pkg/app"
	"gin-api/pkg/config"
	"github.com/spf13/cobra"
	//"gin-api/application/cmd/make"
	"gin-api/pkg/console"
	"os"
)

//go:embed public/*
var Assets embed.FS

//go:embed application/http/view/*
var Views embed.FS

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   config.Get("app.name"),
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// 配置初始化，依赖命令行 --env 参数
			config.InitConfig(cmd.Env)

			// 应用初始化
			app.New(app.WithAsset(Assets), app.WithView(Views))

			//加载其他组件
			bootstrap.Setup()
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
	)

	// 注册默认运行的命令
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册选项
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}





