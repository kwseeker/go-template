package main

import (
	"flag"
	"fmt"
)

/*
-word=foo1 -numb=30 true bar1
-word=foo1 -numb=30 -fork=true -svar=bar1

命令行参数格式以及解析后的数据均保存在全局变量 flag.CommandLine 中
CommandLine 是 FlagSet 类型

	type FlagSet struct {
		Usage func()			//解析错误时会调用Usage,调用Usage之后执行errorHandling
		name          string	//应用名
		parsed        bool		//是否已经解析命令行参数
		actual        map[string]*Flag	//解析后实际参数键值对
		formal        map[string]*Flag	//默认设置的参数键值对
		args          []string 	//实际传参数组
		errorHandling ErrorHandling		//解析产生错误的处理
		output        io.Writer //Usage输出的目标IO，默认输出到标准错误
	}

	Usage = {func()} flag.commandLineUsage
	name = {string} "/tmp/GoLand/___1go_build_kwseeker_top_kwseeker_go_template_basic_flag"
	parsed = {bool} true
	actual = {map[string]*flag.Flag}
	 0 = word ->
	 1 = numb ->
	formal = {map[string]*flag.Flag}
	 0 = word ->
	 1 = numb ->
	 2 = fork ->
	 3 = svar ->
	args = {[]string} len:2, cap:2
	errorHandling = {flag.ErrorHandling} ExitOnError (1)
	output = {io.Writer} nil
*/
func main() {
	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	forkPtr := flag.Bool("fork", false, "a bool")
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *forkPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}
