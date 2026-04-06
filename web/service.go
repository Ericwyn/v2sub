package web

import (
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/gin-contrib/sessions"
	"strings"
	"time"
)

var runFlag = false

func checkLogin(session sessions.Session) bool {
	hadLogin := session.Get("hadLogin")
	if loginFlag, ok := hadLogin.(bool); !ok || !loginFlag {
		return false
	}
	loginTimeValue := session.Get("loginTime")
	loginTime, ok := loginTimeValue.(int64)
	if !ok {
		return false
	}
	return loginTime > 0 && time.Now().Unix()-loginTime <= 60*60*2
}

func v2subConnKill() {
	runLog = append(runLog, "命令执行: "+v2subBinPath+" -conn  kill")
	runLog = append(runLog, "")
	_ = command.RunSyncForResultCb(func(s string) {
		s = strings.Replace(s, "\u0000", "", -1)
		s = strings.Replace(s, "\t", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		log.I(s)
		runLog = append(runLog, "result: "+s)
	}, v2subBinPath, "-conn", "kill")
}

func v2subConnStart() {
	runFlag = true
	lastStartTimeUnix = time.Now().Unix()
	v2subConnKill()
	runLog = append(runLog, "命令执行: "+v2subBinPath+" -conn start")
	runLog = append(runLog, "")
	_ = command.RunSyncForResultCb(func(s string) {
		s = strings.Replace(s, "\u0000", "", -1)
		s = strings.Replace(s, "\t", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		runLog = append(runLog, s)
	}, v2subBinPath, "-conn", "start")

	runFlag = false
}

func v2subSubUpdateAll() (string, error) {
	cmd := v2subBinPath + " -sub updateall"
	return runCmdAndLog(cmd)
}

func v2subSerSet(id string) (string, error) {
	cmd := v2subBinPath + " -ser set " + id
	return runCmdAndLog(cmd)
}

func v2subSerSetX() (string, error) {
	cmd := v2subBinPath + " -ser setx"
	return runCmdAndLog(cmd)
}

func v2subConfList() (string, error) {
	cmd := v2subBinPath + " -conf list"
	return runCmdAndLog(cmd)
}

func v2subConfHttpPort(port string) (string, error) {
	cmd := v2subBinPath + " -conf hport " + port
	return runCmdAndLog(cmd)
}

func v2subConfSocksPort(port string) (string, error) {
	cmd := v2subBinPath + " -conf sport " + port
	return runCmdAndLog(cmd)
}

func v2subConfLocalConnect(trueOfFalse string) (string, error) {
	cmd := v2subBinPath + " -conf lconn " + trueOfFalse
	return runCmdAndLog(cmd)
}

func runCmdAndLog(cmd string) (string, error) {
	runLog = append(runLog, "命令执行: "+cmd)
	runLog = append(runLog, "")
	result, err := command.RunResult(cmd)
	runLog = append(runLog, ", result: "+result)
	if err != nil {
		runLog = append(runLog, err.Error())
	}
	return result, err
}
