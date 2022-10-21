package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fileName := "sampleText"
	fi, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("err", err)
	}
	fileInfo, er := os.Stat(fileName)
	if er != nil {
		fmt.Println("err", er)
		return
	}
	len := fileInfo.Size()
	fmt.Println("fileInfo size is : ", len)
	nextIdx, err := retLastNum(fi, len)
	if err != nil {
		fmt.Println("Atoi, ", err)
	}

	var b []byte
	for i := nextIdx; i < nextIdx+10; i++ {
		a := fmt.Sprintf("this line is %v \n", i)
		b = []byte(a)
		_, er := fi.Write(b)
		if er != nil {
			fmt.Println(er)
			break
		}
	}
	fi.Close()
}

// return last number (startIdx, endIdx)
func retLastNum(file *os.File, len int64) (uint64, error) {
	var bufLen int64 = 15
	if len < 15 {
		fmt.Println("hi", len)
		if len == 0 {
			return 0, nil
		}
		bufLen = len
	}
	buf := make([]byte, bufLen)
	var offset int64 = -15
	var nlpos int = 0
	tmpStr := ""
	ret := ""
	checkFirstNum := false
	for nlpos == 0 {
		idx, _ := file.Seek(offset, 2)
		file.ReadAt(buf, idx)
		for i := bufLen - 1; i >= 0; i-- {
			tmpChar := string(buf[i])
			if tmpChar == "\n" || tmpChar == " " {
				if !checkFirstNum {
					continue
				}
				nlpos++
				break
			}
			if buf[i] < 48 {
				nlpos++
				break
			}
			if buf[i] > 57 {
				nlpos++
				break
			}
			checkFirstNum = true
			tmpStr = tmpChar + tmpStr
		}
		offset += offset
		ret = tmpStr + ret
		tmpStr = ""
	}
	ans, _ := strconv.ParseUint(ret, 10, 64)
	fmt.Println(ret, ans+1)
	return ans + 1, nil
}
