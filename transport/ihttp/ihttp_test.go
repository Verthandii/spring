package ihttp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Verthandii/spring/ilogger"
)

func TestBaidu(t *testing.T) {
	cli, err := NewClient(WithEndpoint("http://localhost:6001/ping"), WithTimeout(time.Second*10))
	if err != nil {
		t.Fatal(err)
	}
	req, _ := cli.BuildRequest(context.TODO(), "GET", "", "", nil)
	req.Header.Set("X", "123")
	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestNewServer(t *testing.T) {
	client, err := NewClient(WithEndpoint("http://localhost:8000"), WithTimeout(time.Hour*10))
	if err != nil {
		return
	}
	go func() {
		s := NewServer(Timeout(time.Hour))
		s.Use(func(c *Context) {
			ilogger.DebugwCtx(c, "Server OK")
		})
		s.Run()
	}()
	time.Sleep(time.Second)
	ctx := context.Background()
	err = client.Get(ctx, "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
}
