package fileoperate

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Readfile() {
	fs, err := os.Open("go.mod")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer fs.Close()
	n, err := fs.Read([]byte(""))
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		if n < 19 {
			fmt.Println("read file is empty!")
			return
		}
		fmt.Println(n)

	}
}

func Readbuff() {
	fs, err := os.Open("go.mod")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer fs.Close()
	bf := bufio.NewReader(fs)
	//buf := make([]byte, 1024)
	for {
		n, err := bf.ReadString('\n')
		//st = strings.TrimSpace(st)
		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
		}
		if err == io.EOF {
			break
		}
		fmt.Println(n)
	}

}
func ReadIoUtil() {
	fs, err := os.Open("go.mod")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer fs.Close()
	st, err := ioutil.ReadAll(fs)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(st))

}
func WriteFile() {
	fs, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fs.Close()
	n, err := fs.WriteString("good\n")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(n)
}

func WriteBufFile() {
	fs, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fs.Close()
	buf := bufio.NewWriter(fs)
	n, err := buf.WriteString("s string")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buf.Flush()
	fmt.Println(n)
}

func WriteIoUtil() {
	li := make([]string, 0)
	var str string
	fmt.Print("please input String:\n")
	fmt.Scanln(&str)
	li = append(li, str)
	for _, value := range li {
		fmt.Println("value:", value)
		err := os.WriteFile("test.log", []byte(value), 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

}
func Istype(a any) {
	//a = 1000
	switch v := a.(type) {
	case int8:
		fmt.Println("int8:", v)
	case int32:
		fmt.Println("int32:", v)
	case string:
		fmt.Println("string:", v)
	case int:
		fmt.Println("int:", v)
	default:
		fmt.Println("error intput:", v)

	}
}

// 错误处理函数
func ErrorHanding(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
