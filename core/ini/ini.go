package ini

//配置的读取
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const ROOT = ""

/*取值对象*/
type ItemMap interface {
	Str(string) string
	Int(string) int
	Ok(string) bool
	String() string
}

//bool list
var bool_list []string = []string{"yes", "ok", "1", "true"}

type itemObj struct {
	ptr string /*global only key*/
	k   string /*my key*/
	v   string /*my val*/
}

func newItem(k string) *itemObj {
	return &itemObj{k: k}
}

func (this *itemObj) child(k string, v string) *itemObj {
	if this.k == ROOT {
		return &itemObj{ptr: k, k: k, v: v}
	}
	return &itemObj{ptr: this.k + "." + k, k: k, v: v}
}

func (this *itemObj) append(v string) {
	this.v += v
}

func (this *itemObj) Str() string {
	return this.v
}

func (this *itemObj) Int() int {
	if n, err := strconv.Atoi(this.v); err == nil {
		return n
	}
	return 0
}

func (this *itemObj) Ok() bool {
	low_val := strings.ToLower(this.v)
	for _, v := range bool_list {
		if low_val == v {
			return true
		}
	}
	return false
}

func (this *itemObj) String() string {
	return fmt.Sprintf("%s {\n    %s : %s\n};", this.ptr, this.k, this.v)
}

/*
items(key区分大小写)
*/
type itemSet map[string]*itemObj

func (items itemSet) Str(k string) string {
	if item, ok := items[k]; ok {
		return item.Str()
	}
	return ""
}

func (items itemSet) Int(k string) int {
	if item, ok := items[k]; ok {
		return item.Int()
	}
	return 0
}

func (items itemSet) Ok(k string) bool {
	if item, ok := items[k]; ok {
		return item.Ok()
	}
	return false
}

func (items itemSet) String() string {
	var str string = "[############ INI Begin #########]\n"
	for _, item := range items {
		str += item.String() + "\n"
	}
	str += "[############ INI End ############]"
	return str
}

//解析服务器配置
func Unmarshal(path string) (ItemMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	items := make(itemSet)
	item := newItem(ROOT)
	var current *itemObj
	reader := bufio.NewReader(file)
	for {
		b, _, c := reader.ReadLine()
		if c == io.EOF {
			break
		}
		size := len(b)
		if size == 0 {
			continue
		}
		str := strings.Trim(string(b), " ") // remove left and right blank
		if str[0] == '#' {                  // this is annotation
			continue
		}
		if str[0] == '[' && str[size-1] == ']' {
			item.k = str[1 : size-1] //root key
			current = nil            //root do't append val
			continue
		}
		if i := strings.Index(str, "="); i != -1 {
			current = item.child(str[:i], str[i+1:])
			items[current.ptr] = current
		} else {
			if current != nil {
				current.append(str)
			}
		}
	}
	fmt.Println(items)
	return items, nil
}
