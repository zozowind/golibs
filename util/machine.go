package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

// ReadMachineID generates machine id and puts it into the machineId global variable.
// If this function fails to get the hostname, it will cause a runtime error.
func ReadMachineID() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Errorf("cannot get hostname: %v", err))
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	return hex.EncodeToString(hw.Sum(nil))[16:24]
}
