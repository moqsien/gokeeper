package kctrl

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/abiosoft/ishell/v2"
	logger "github.com/moqsien/processes/logger"
)

func (that *KCtrl) InitClient() {
	if that.KcMode == CtrlUnknown {
		that.InitSockPath()
		that.Client = &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", that.UnixSockPath)
				},
			},
		}
		that.Shell = ishell.New()
		that.KcMode = CtrlClient
	}
}

type ParamsContainer map[string]string

func (that *KCtrl) CtrlGetStr(urlPath string, params ParamsContainer) (string, error) {
	if that.KcMode == CtrlClient {
		urlPath = strings.Trim(urlPath, "/")
		paramStr := ""
		for k, v := range params {
			if len(paramStr) == 0 {
				paramStr += fmt.Sprintf("?%s=%s", k, v)
			} else {
				paramStr += fmt.Sprintf("&%s=%s", k, v)
			}
		}
		resp, _ := that.Get(fmt.Sprintf("http://%s/%s/%s", that.UnixSockName, urlPath, paramStr))
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			return "", err
		}
		return string(content), nil
	}
	return "", errors.New("not initialized")
}

func (that *KCtrl) CtrlClientStart() {
	if that.KcMode == CtrlClient {
		that.Shell.Run()
	} else if that.KcMode == CtrlSever {
		logger.Error("Alreay initialized to server!")
		return
	} else {
		logger.Error("Not initialized!")
		return
	}
}
