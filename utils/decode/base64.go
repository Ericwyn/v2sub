package decode

import (
	"encoding/base64"
	"github.com/Ericwyn/v2sub/utils/log"
)

func Base64Decode(base64Str string) string {
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	decodestr := string(decoded)
	if err == nil {
		return decodestr
	} else {
		log.E("base64 decode fail")
		log.E(base64Str)
		return ""
	}
}

func VmessBase64Decode(vmessBase64Str string) string {
	decoded, err := base64.RawURLEncoding.DecodeString(vmessBase64Str)
	decodestr := string(decoded)
	if err == nil {
		return decodestr
	} else {
		log.E("vmess base64 decode fail")
		log.E(vmessBase64Str)
		return ""
	}
}
