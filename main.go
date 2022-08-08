package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	socks5 "github.com/armon/go-socks5"
)

var (
	l    = flag.String("l", "18889", "端口 - 监听的端口")
	auth = flag.String("a", "", "密钥 - 生成访问密钥")
	dns  = flag.String("s", "", "DNS - 设置DNS")
	help = flag.Bool("h", false, "帮助 - 查看所有的帮助")
	d    = flag.Bool("d", false, "后台 - 后台运行系统")
)

func forkDaemon() int {
	args := os.Args
	os.Setenv("__Daemon", "true")
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	pid, _ := syscall.ForkExec(os.Args[0], args, procAttr)
	return pid
}

func init() {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *d && os.Getenv("__Daemon") != "true" {
		fmt.Println("后台运行 [PID]", forkDaemon())
		os.Exit(0)
	}
}

func usage() {
	fmt.Printf("以下为命令提示\n")
	flag.PrintDefaults()
}

func main() {
	cred := socks5.StaticCredentials{}
	d := socks5.DNSResolver{}
	if *dns != "" {
		d = socks5.DNSResolver{}
		ctx := context.Background()
		d.Resolve(ctx, *dns)
	}
	if *auth != "" {
		auths := strings.Split(*auth, ":")
		cred = socks5.StaticCredentials{auths[0]: auths[1]}
	}
	conf := &socks5.Config{Credentials: cred, Resolver: d}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}
	err = server.ListenAndServe("tcp", ":"+*l)
	if err != nil {
		panic(err)
	}
}
