package httpsrv

import (
	"gopkg.in/goyy/goyy.v0/comm/log"
)

var logger = log.New("[comm-httpsrv]")

func SetPriority(value int) {
	logger.SetPriority(value)
}
