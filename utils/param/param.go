package param

import (
	"github.com/Ericwyn/v2sub/utils/log"
	"os"
)

func AssistParamLength(param []string, length int) {
	if len(param) < length {
		log.E("param length error, len(params) need greater than ", length, " args now : ", param)
		os.Exit(-1)
	}
}
