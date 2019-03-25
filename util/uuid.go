package util

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
)

var requestCounter uint32
var machineID = ReadMachineID()

//UID Get UID
func UID() string {
	id := uuid.NewV4()
	return strings.Replace(id.String(), "-", "", -1)
}

//UIDToUUID uid to UUID struct
func UIDToUUID(uid string) (*uuid.UUID, error) {
	if len(uid) == 32 {
		str := fmt.Sprintf("%s-%s-%s-%s-%s", uid[0:8], uid[8:12], uid[12:16], uid[16:20], uid[20:32])
		r, err := uuid.FromString(str)
		if nil != err {
			return nil, err
		}
		return &r, nil
	}
	return nil, fmt.Errorf("wrong length")
}

//NewRequestID 请求id
func NewRequestID() string {
	formatTime := time.Now().Format("20060102150405.999")
	formatTime = strings.Replace(formatTime, ".", "", 1)
	formatTime += strings.Repeat("0", 17-len(formatTime))
	orderNum := atomic.AddUint32(&requestCounter, 1) % 10000000
	return fmt.Sprintf("%s%8s%07d", formatTime, machineID, orderNum)
}
