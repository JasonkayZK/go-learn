package server

import (
	"encoding/json"
	"fmt"
	"github.com/jasonkayzk/distributed-id-generator/config"
	"github.com/jasonkayzk/distributed-id-generator/core"
	"net"
	"net/http"
	"strconv"
	"time"
)

type allocResponse struct {
	RespCode int    `json:"resp_code"`
	Msg      string `json:"msg"`
	Id       int    `json:"id"`
}

type healthResponse struct {
	RespCode int    `json:"resp_code"`
	Msg      string `json:"msg"`
	Left     int    `json:"left"`
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	var (
		appTag string
	)
	healthResp := healthResponse{}
	err := r.ParseForm()
	if err != nil {
		goto RESP
	}
	if appTag = r.Form.Get("app_tag"); appTag == "" {
		err = fmt.Errorf("need app_tag param")
		goto RESP
	}
	healthResp.Left = core.GlobalIdAllocator.LeftCount(appTag)
	if healthResp.Left == 0 {
		err = fmt.Errorf("no available id")
		goto RESP
	}

RESP:
	if err != nil {
		healthResp.RespCode = -1
		healthResp.Msg = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
	} else {
		healthResp.Msg = "success"
	}
	if bytes, err := json.Marshal(&healthResp); err == nil {
		_, _ = w.Write(bytes)
	} else {
		w.WriteHeader(500)
		healthResp.Msg = fmt.Sprintf("%v", err)
	}
}

func handleAlloc(w http.ResponseWriter, r *http.Request) {
	var (
		resp   = allocResponse{}
		err    error
		bytes  []byte
		appTag string
	)

	if err = r.ParseForm(); err != nil {
		goto RESP
	}
	if appTag = r.Form.Get("app_tag"); appTag == "" {
		err = fmt.Errorf("need app_tag param")
		goto RESP
	}

	for {
		if resp.Id, err = core.GlobalIdAllocator.NextId(appTag); err != nil {
			goto RESP
		}
		if resp.Id != 0 { // 跳过ID=0, 一般业务不支持ID=0
			break
		}
	}
RESP:
	if err != nil {
		resp.RespCode = -1
		resp.Msg = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
	} else {
		resp.Msg = "success"
	}
	if bytes, err = json.Marshal(&resp); err == nil {
		_, _ = w.Write(bytes)
	} else {
		w.WriteHeader(500)
		resp.Msg = fmt.Sprintf("%v", err)
	}
}

func StartServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/alloc", handleAlloc)
	mux.HandleFunc("/health", handleHealth)

	srv := &http.Server{
		ReadTimeout:  time.Duration(config.AppConfig.HttpReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.AppConfig.HttpWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.AppConfig.HttpPort))
	if err != nil {
		return err
	}

	fmt.Printf("server started at: localhoost:%d\n", config.AppConfig.HttpPort)
	return srv.Serve(listener)
}
