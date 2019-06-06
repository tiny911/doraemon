package localtracing

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
	"sync/atomic"
	"time"
)

const defaultHostname = "localhost"

var (
	objectIdCounter uint32 = 0
	machineId       []byte = loadMachineId()
)

func genId(seeds ...interface{}) string {
	return NewObjectId().Hex()
}

func loadMachineId() []byte {
	var (
		err      error
		sum      [3]byte
		hostname string
	)
	machineId := sum[:]

	hostname, err = os.Hostname()
	if err != nil {
		hostname = defaultHostname
		log.Printf("loadMachineId_Get_HostHame_failed, errmsg:%s.", err)
	}

	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(machineId, hw.Sum(nil))

	return machineId
}

type ObjectId string

// 4byte 时间，
// 3byte 机器ID
// 2byte pid
// 3byte 自增ID
func NewObjectId() ObjectId {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return ObjectId(b[:])
}

// 返回16进制对应的字符串
func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}
