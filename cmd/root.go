/*
Copyright © 2025 keti

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "v0.0.1"
var configFile string // 配置文件路径

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "getstartwithgo",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// 命令出错时，不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
	SilenceUsage: true, // 指定调用 cmd.Execute() 时，执行的 Run 函数
	RunE: func(cmd *cobra.Command, args []string) error {
		// 创建默认的应用命令行选项
		opts := NewServerOptions()

		return run(opts)
	},
	// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./fg-apiserver param1 param2
	Args:    cobra.NoArgs,
	Version: version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.getstartwithgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 初始化配置函数，在每个命令运行时调用
	cobra.OnInitialize(onInitialize)

	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the fg-apiserver configuration file.")
}

func run(opts *ServerOptions) error {
	// 初始化 slog
	initLog()
	// 将 viper 中的配置解析到选项 opts 变量中.
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	// 对命令行选项值进行校验.
	if err := opts.Validate(); err != nil {
		return err
	}

	// fmt.Printf("Read MySQL host from Viper: %s\n", viper.GetString("mysql.host"))
	// fmt.Printf("Read MySQL username from opts: %s\n", opts.MySQLOptions.Username)
	// slog.Info("Read MySQL username from opts: %s\n", opts.MySQLOptions.Username)

	jsonData, _ := json.MarshalIndent(opts, "", "  ")
	// fmt.Println(string(jsonData))
	slog.Info(string(jsonData))
	return nil
}

// initLog 初始化全局日志实例
func initLog() {
	// 获取日志配置
	format := viper.GetString("log.format") // 日志格式，支持：json、text
	level := viper.GetString("log.level")   // 日志级别，支持：debug, info, warn, error
	output := viper.GetString("log.output") // 日志输出路径，支持：标准输出stdout和文件

	// 转换日志级别
	var slevel slog.Level
	switch level {
	case "debug":
		slevel = slog.LevelDebug
	case "info":
		slevel = slog.LevelInfo
	case "warn":
		slevel = slog.LevelWarn
	case "error":
		slevel = slog.LevelError
	default:
		slevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: slevel}

	var w io.Writer
	var err error
	// 转换日志输出路径
	switch output {
	case "":
		w = os.Stdout
	case "stdout":
		w = os.Stdout
	default:
		w, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 转换日志格式
	if err != nil {
		return
	}
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewJSONHandler(w, opts)

	}

	// 设置全局的日志实例为自定义的日志实例
	slog.SetDefault(slog.New(handler))
}
