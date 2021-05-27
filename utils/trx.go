package utils

import (
	"hash/adler32"
	"strconv"
	"strings"
	"time"
)

func GenerateTrxCode(trxType string) string {
	return strings.ToUpper(trxType[:3]) + "-" + strconv.FormatUint(uint64(adler32.Checksum([]byte(strconv.Itoa(int(time.Now().UnixNano()))))), 16)
}
