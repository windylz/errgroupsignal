

package main

import (
	"github.com/kataras/iris/core/errors"
	"net/http"
	"fmt"
	"golang.org/x/sync/errgroup"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	//http配置
	srv := &http.Server{Addr: ":9001"}
	http.HandleFunc("/",index)

	//linux信号配置
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

    //errgroup
	group, _ := errgroup.WithContext(context.Background())

	//启动http
	group.Go(func() error {
		return srv.ListenAndServe()
	})

	//关闭http
	//group.Go(func() error {
	//	return srv.Shutdown(nil)
	//})

	//信号监听
	/*
	group.Go(func() error {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", s)
				fmt.Println("通过guest:control+c关闭web")
				return srv.Shutdown(nil)
				fmt.Println(s)
				return  errors.New("guest:control+c")
			case syscall.SIGUSR1:
				fmt.Println("usr1 signal", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2 signal", s)
			default:
				fmt.Println("other signal", s)
			}
		}
		return errors.New("end")
	})
	 */

	//信号监听，及通过信号关闭web
	group.Go(func() error {
		c := make(chan os.Signal)
		signal.Notify(c)
		fmt.Println("start..")
		s := <-c
		if s== syscall.SIGINT{
			fmt.Println("通过guest:control+c关闭web")
			srv.Shutdown(nil)
			return nil
		}
		return errors.New("end")
	})

	if err := group.Wait(); err != nil {
		fmt.Println("err:",err)
	}
}

func index(w http.ResponseWriter,r *http.Request)  {
   fmt.Fprintf(w,"This is index,haha")
}

