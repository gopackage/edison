package mraa

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func writeFile(path string, data []byte) (i int, err error) {
	fmt.Println(">>", path, string(data))
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i = 0
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}

func parseInt(input string, bits int) (int64, error) {
	str := strings.TrimSpace(input)
	if str == "" {
		return 0, nil
	}

	return strconv.ParseInt(str, 10, bits)
}
