package sample

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init()  {
	rand.Seed(time.Now().Unix())
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func RandomStr(n int) string {
	var res strings.Builder
	k := len(alphabet)
	for i:=0; i < n; i ++ {
			res.WriteByte(alphabet[rand.Intn(k)])
	}
	return res.String()
}
