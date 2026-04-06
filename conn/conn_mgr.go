package conn

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"github.com/Ericwyn/v2sub/utils/putil"
	"github.com/gin-gonic/gin"
	"net/http"
)

const v2rayBinPath = "/usr/local/bin/v2ray"
const runtimeDirPath = "/etc/v2sub"
const v2rayPidFileName = "v2ray.pid"
const v2rayRuntimeFileName = "runtime.json"

var pacFilePath = "/etc/v2sub/v2sub.pac"

var defaultPacText = `
// 默认全部直连模式
function FindProxyForURL(url, host) {
   return 'DIRECT';
}
`

type V2rayStatus struct {
	Running       bool   `json:"running"`
	Pid           int    `json:"pid"`
	StartedAtUnix int64  `json:"startedAtUnix"`
	BinaryPath    string `json:"binaryPath"`
	ConfigPath    string `json:"configPath"`
	Mode          string `json:"mode"`
	Message       string `json:"message"`
}

func ParseArgs(args []string) {
	param.AssistParamLength(args, 1)
	switch args[0] {
	case "start":
		status, err := StartV2ray()
		if err != nil {
			log.E(err.Error())
			os.Exit(-1)
		}
		printStatus(status)
	case "kill":
		status, err := KillV2Sub()
		if err != nil {
			log.E(err.Error())
			os.Exit(-1)
		}
		printStatus(status)
	case "status":
		printStatus(Status())
	case "restart":
		status, err := RestartV2ray()
		if err != nil {
			log.E(err.Error())
			os.Exit(-1)
		}
		printStatus(status)
	case "start-pac":
		readPacConfigFile()
		go startPacServerOnly()
		status, err := StartV2rayWithMode("pac")
		if err != nil {
			log.E(err.Error())
			os.Exit(-1)
		}
		printStatus(status)
	default:
		log.E("sub args error")
	}
}

func checkV2ray() error {
	vtoo := file.OpenFile(v2rayBinPath)
	if !vtoo.Exits() {
		return errors.New("can't find v2ray bin in " + v2rayBinPath)
	}
	return nil
}

func StartV2ray() (V2rayStatus, error) {
	return StartV2rayWithMode("default")
}

func StartV2rayWithMode(mode string) (V2rayStatus, error) {
	log.I("start v2ray ......")

	if err := checkV2ray(); err != nil {
		return StatusWithMessage(false, 0, mode, err.Error()), err
	}

	if len(conf.ServerConfigNow.ServerList) == 0 {
		err := errors.New("server list is empty")
		return StatusWithMessage(false, 0, mode, err.Error()), err
	}
	if conf.ServerConfigNow.Id < 0 || conf.ServerConfigNow.Id >= len(conf.ServerConfigNow.ServerList) {
		err := fmt.Errorf("server id out of range: %d", conf.ServerConfigNow.Id)
		return StatusWithMessage(false, 0, mode, err.Error()), err
	}

	currentStatus := Status()
	if currentStatus.Running {
		currentStatus.Message = "v2ray is already running"
		return currentStatus, errors.New(currentStatus.Message)
	}

	runConfig := conf.ServerConfigNow.ServerList[conf.ServerConfigNow.Id]
	conf.SaveDefaultServerConfig(runConfig)

	log.I("use config is :   ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ")
	log.I("========================================================================")
	log.I(
		putil.F("ID", 4),
		putil.F("别名", 50),
		putil.F("地址", 24),
		putil.F("端口", 10),
		putil.F("类型", 5),
	)
	log.I(putil.F(" "+strconv.Itoa(conf.ServerConfigNow.Id), 4),
		putil.F(runConfig.Vmess.Ps, 50),
		putil.F(runConfig.Vmess.Add, 24),
		putil.F(runConfig.Vmess.Port, 10),
		putil.F(runConfig.Vmess.Type, 5))
	log.I("========================================================================")
	log.I("v2ray config path : " + conf.GetV2rayConfigPath())

	args := []string{"-config", conf.GetV2rayConfigPath()}
	if useNewV2rayVersion() {
		args = []string{"run", "-c", conf.GetV2rayConfigPath()}
	}

	cmd, err := command.StartBackground(v2rayBinPath, args...)
	if err != nil {
		return StatusWithMessage(false, 0, mode, "start v2ray error: "+err.Error()), err
	}

	status := V2rayStatus{
		Running:       true,
		Pid:           cmd.Process.Pid,
		StartedAtUnix: time.Now().Unix(),
		BinaryPath:    v2rayBinPath,
		ConfigPath:    conf.GetV2rayConfigPath(),
		Mode:          mode,
		Message:       "v2ray started",
	}
	if err = writeRuntimeState(status); err != nil {
		_ = command.KillProcess(cmd.Process.Pid)
		return StatusWithMessage(false, 0, mode, err.Error()), err
	}

	return status, nil
}

func RestartV2ray() (V2rayStatus, error) {
	stopStatus, err := KillV2Sub()
	if err != nil && stopStatus.Message != "v2ray is not running" {
		return stopStatus, err
	}
	mode := stopStatus.Mode
	if mode == "" {
		mode = "default"
	}
	return StartV2rayWithMode(mode)
}

func Status() V2rayStatus {
	status, err := readRuntimeStatus()
	if err != nil {
		return StatusWithMessage(false, 0, "", "v2ray is not running")
	}
	if status.Pid <= 0 || !command.IsProcessAlive(status.Pid) {
		clearRuntimeState()
		return StatusWithMessage(false, 0, status.Mode, "v2ray is not running")
	}
	status.Running = true
	if status.BinaryPath == "" {
		status.BinaryPath = v2rayBinPath
	}
	if status.ConfigPath == "" {
		status.ConfigPath = conf.GetV2rayConfigPath()
	}
	if status.Message == "" {
		status.Message = "v2ray is running"
	}
	return status
}

func StatusWithMessage(running bool, pid int, mode string, message string) V2rayStatus {
	return V2rayStatus{
		Running:       running,
		Pid:           pid,
		StartedAtUnix: 0,
		BinaryPath:    v2rayBinPath,
		ConfigPath:    conf.GetV2rayConfigPath(),
		Mode:          mode,
		Message:       message,
	}
}

func KillV2Sub() (V2rayStatus, error) {
	status, err := readRuntimeStatus()
	if err != nil {
		clearRuntimeState()
		return StatusWithMessage(false, 0, "", "v2ray is not running"), nil
	}

	if status.Pid <= 0 || !command.IsProcessAlive(status.Pid) {
		clearRuntimeState()
		status.Running = false
		status.Pid = 0
		status.Message = "v2ray is not running"
		return status, nil
	}

	_ = command.TerminateProcess(status.Pid)
	for i := 0; i < 10; i++ {
		if !command.IsProcessAlive(status.Pid) {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	if command.IsProcessAlive(status.Pid) {
		if err = command.KillProcess(status.Pid); err != nil {
			status.Message = "kill v2ray error: " + err.Error()
			return status, err
		}
	}

	clearRuntimeState()
	status.Running = false
	status.Pid = 0
	status.Message = "v2ray stopped"
	return status, nil
}

func readPacConfigFile() {
	pacFile := file.OpenFile(pacFilePath)
	if pacFile.Exits() {
		read, err := pacFile.Read()
		if err != nil {
			log.E("read pac config error, use default pac config")
		} else {
			defaultPacText = string(read)
		}
	} else {
		log.E("read pac config error, pacFile in '" + pacFilePath + "' not exits" +
			", use default pac config")
	}
}

func startPacServerOnly() {
	router := gin.Default()
	router.GET("v2sub.pac", func(context *gin.Context) {
		context.File(pacFilePath)
	})

	s := &http.Server{
		Addr:           ":23333",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}

func useNewV2rayVersion() bool {
	result, err := command.RunResult("v2ray -version")
	if err != nil {
		log.I("check new v2ray version, 'v2ray -version' get error msg, use new v2ray version cmd")
		return true
	}
	log.I("check new v2ray version, 'v2ray -version' get result, ", result)
	return false
}

func printStatus(status V2rayStatus) {
	fmt.Println("running:", status.Running)
	fmt.Println("pid:", status.Pid)
	fmt.Println("startedAtUnix:", status.StartedAtUnix)
	fmt.Println("binaryPath:", status.BinaryPath)
	fmt.Println("configPath:", status.ConfigPath)
	fmt.Println("mode:", status.Mode)
	fmt.Println("message:", status.Message)
}

func getPidFilePath() string {
	return filepath.Join(runtimeDirPath, v2rayPidFileName)
}

func getRuntimeFilePath() string {
	return filepath.Join(runtimeDirPath, v2rayRuntimeFileName)
}

func ensureRuntimeDir() error {
	return os.MkdirAll(runtimeDirPath, 0755)
}

func writeRuntimeState(status V2rayStatus) error {
	if err := ensureRuntimeDir(); err != nil {
		return err
	}
	if err := os.WriteFile(getPidFilePath(), []byte(strconv.Itoa(status.Pid)), 0644); err != nil {
		return err
	}
	body, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(getRuntimeFilePath(), body, 0644); err != nil {
		return err
	}
	return nil
}

func readRuntimeStatus() (V2rayStatus, error) {
	pid, err := readPid()
	if err != nil {
		return V2rayStatus{}, err
	}
	status := V2rayStatus{
		Pid:        pid,
		BinaryPath: v2rayBinPath,
		ConfigPath: conf.GetV2rayConfigPath(),
	}
	runtimeBytes, readErr := os.ReadFile(getRuntimeFilePath())
	if readErr == nil && len(runtimeBytes) > 0 {
		_ = json.Unmarshal(runtimeBytes, &status)
		status.Pid = pid
	}
	return status, nil
}

func readPid() (int, error) {
	pidBytes, err := os.ReadFile(getPidFilePath())
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil || pid <= 0 {
		return 0, errors.New("invalid pid file")
	}
	return pid, nil
}

func clearRuntimeState() {
	_ = os.Remove(getPidFilePath())
	_ = os.Remove(getRuntimeFilePath())
}

func isProcessNotExistError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH)
}
