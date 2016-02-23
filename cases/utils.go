package bender

import (
	"bufio"
	"io"
	"log"
	"os"
)

func CheckExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	}
	// exists but there is error
	return true, err
}

// Parameter `flag ...int` when flag = []int{1} will indicate a `xxxx.set` type, and the
// first part need to keep not drop.
func SplitMultiYamlToSingle(fp string, flag ...int) [][]byte {

	r, e := CheckExists(fp)
	if !r || e != nil {
		log.Fatalf("%s is not exists or error happens", fp)
	}
	f, err := os.Open(fp)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var total [][]byte
	var tmp []byte
	for {
		line, err := reader.ReadBytes('\n')
		if string(line) == "---\n" || err == io.EOF {
			v := []byte{}
			v = append(v, tmp...)
			total = append(total, v)
			tmp = []byte{}
		} else {
			tmp = append(tmp, line...)
		}
		if err != nil {
			if err != io.EOF {
				log.Fatalln(err)
			}
			break
		}
	}
	if len(total[0]) == 0 {
		if len(flag) != 0 && flag[0] == 1 {
			total = total[1:]
		} else {
			total = total[2:]
		}

	} else {
		total = total[1:]
	}
	return total
}

type SimpleSet map[string]bool

func (this SimpleSet) Add(key string) {
	if _, ok := this[key]; !ok {
		this[key] = true
	} else {
		log.Printf("Key [%s] already exists", key)
	}
}

func (this SimpleSet) Del(key string) {
	panic("Not Implement")
}

func (this SimpleSet) AllKeys() []string {
	rv := make([]string, 0, 50)
	for k, _ := range this {
		rv = append(rv, k)
	}
	return rv
}
