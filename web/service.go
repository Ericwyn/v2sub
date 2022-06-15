package web

import (
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"strings"
	"time"
)

var runFlag = false

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
		//fmt.Print(s)
		s = strings.Replace(s, "\u0000", "", -1)
		s = strings.Replace(s, "\t", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		//s = strings.Split(s, "\n")[0]
		runLog = append(runLog, s)
	}, v2subBinPath, "-conn", "start")

	// 阻塞住了
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
