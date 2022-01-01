package comm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func IsFlagTrue(flag []bool) bool {
	for _, v := range flag {
		if !v {
			return false
		}
	}

	return true
}

func IsFieldNull(value interface{}) bool {
	t := reflect.TypeOf(value).String()

	if t == "string" {
		if value == "nullnullnull" {
			return true
		}
	} else if strings.Contains(t, "int") {
		if value == -998998 {
			return true
		}
	} else if strings.Contains(t, "float") {
		if value == -998998.98 {
			return true
		}
	}

	return false
}

// to string
func ToSting(value interface{}) string {
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case float64:
		return strconv.FormatFloat(value.(float64), 'E', -1, 64)
	case string:
		return value.(string)
	default:
		return ""
	}
}

func GetFloat()  {
	
}

func GetValueString(value interface{}) (string, error) {
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int)), nil
	case int64:
		return strconv.FormatInt(value.(int64), 10), nil
	case float64:
		return strconv.FormatFloat(value.(float64), 'G', -1, 64), nil
	case string:
		return value.(string), nil
	default:
		return "", errors.New("类型错误")
	}
}

// 毫秒
func Msectime() int64 {
	return time.Now().UnixNano() / 1e6
}

// to int
func ToInt(value interface{}) int {
	switch value.(type) {
	case int:
		return value.(int)
	case int64:
		return int(value.(int64))
	case float64:
		return int(value.(float64))
	case string:
		val, err := strconv.Atoi(value.(string))
		if err != nil {
			return 0
		} else {
			return val
		}
	default:
		return 0
	}
}

// to int64
func ToInt64(value interface{}) int64 {
	switch value.(type) {
	case int:
		return int64(value.(int))
	case int64:
		return value.(int64)
	case float64:
		return int64(value.(float64))
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return 0
		} else {
			return val
		}
	default:
		return 0
	}
}

func Max(x interface{},y interface{}) int {
	xInt := ToInt(x)
	yInt := ToInt(y)
	if xInt > yInt {
		return xInt
	}
	return yInt
}

func Min(x interface{},y interface{}) int {
	xInt := ToInt(x)
	yInt := ToInt(y)
	if xInt < yInt {
		return xInt
	}
	return yInt
}

func IntMaxUnsigned() int {
	return 4294967295
}

/**
 * 驼峰转蛇形 snake string
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @date 2020/7/30
 * @param s 需要转换的字符串
 * @return string
 **/
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}


/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2020/7/30
 * @param s要转换的字符串
 * @return string
 **/
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 生成N个元素的一维数组
func ArrayFill(value string, num int) []string {
	var arr []string
	for i := 0; i < num; i++ {
		arr = append(arr, value)
	}

	return arr
}

func InArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1;

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}

	return
}
func MapKeyExists(key string, mapString interface{}) bool {
	mapType := reflect.TypeOf(mapString).String()
	if mapType == "map[string]string" {
		for k, _ := range mapString.(map[string]string) {
			if key == ToSting(k) {
				return true
			}
		}
	} else {
		panic("comm/base MapKeyExists Not Type " + mapType)
	}

	return false
}

// 判断文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func ToJsonNotErr(arr interface{}) string {
	jsonByte, _ := json.Marshal(arr)
	return string(jsonByte)
}

func StructFieldsIdenticalCopy(new ,old interface{})  {
	fieldsType := reflect.TypeOf(new).Elem()
	fieldsVal := reflect.ValueOf(new).Elem()
	fieldsTypeOld := reflect.TypeOf(old)
	fieldsValOld := reflect.ValueOf(old)

	//if fieldsVal.Kind() != reflect.Ptr || fieldsVal.Kind() != reflect.Ptr {
	//	log.Panic("StructFieldsIdenticalCopy Err .")
	//}
	for i := 0; i < fieldsType.NumField(); i++ {
		key := fieldsType.Field(i).Tag.Get("json")
		for old_i := 0; old_i < fieldsTypeOld.NumField(); old_i++  {
			keyOld := fieldsTypeOld.Field(old_i).Tag.Get("json")
			if key == keyOld {
				fieldsVal.Field(i).Set(fieldsValOld.Field(old_i))
				break
			}
		}
	}
}

// 获取雪花ID
func GetSnowflakeId() int64 {
	node, _ := snowflake.NewNode(10)
	return node.Generate().Int64()
}

// 返回old数组中键值为key的列， 如果指定了可选参数index，那么old数组中的这一列的值将作为返回数组中对应值的键。
func ArrayColumn(new , old interface{}, key string, index ...string) error {
	//newType := reflect.TypeOf(new)
	newVal := reflect.ValueOf(new)
	if newVal.Kind() != reflect.Ptr {
		log.Panic("ArrayColumn Params(New) Err: ", new)
		return errors.New("ArrayColumn Params(New) Err ")
	}

	oldType := reflect.TypeOf(old)
	oldVal := reflect.ValueOf(old)
	if oldVal.Kind() != reflect.Slice && oldVal.Kind() != reflect.Array {
		log.Panic("ArrayColumn Params(Old) Err: ", old)
		return errors.New("ArrayColumn Params(Old) Err ")
	}
	if oldType.Elem().Kind() != reflect.Struct {
		log.Panic("ArrayColumn Elem Params(Old) Err: ", new)
		return errors.New("ArrayColumn Elem Params(Old) Err ")
	}

	if len(index) > 0 {
		indexKey := index[0]
		err := arrayIndexColumn(new, old, key, indexKey)
		if err != nil {
			log.Panic("ArrayColumn arrayIndexColumn  Err: ", err.Error())
		}
		return err
	}

	err := arrayColumn(new, old, key)
	if err != nil {
		log.Panic("ArrayColumn arrayColumn  Err: ", err.Error())
	}
	return err
}
func arrayColumn(new, old interface{}, key string) (err error) {
	if len(key) == 0 {
		return errors.New("key cannot not be empty")
	}

	newElemType := reflect.TypeOf(new).Elem()
	newElemTypeKind := newElemType.Kind()
	if newElemTypeKind != reflect.Map && newElemTypeKind != reflect.Slice && newElemTypeKind != reflect.Array{
		return errors.New("new must be slice")
	}

	oldType := reflect.TypeOf(old)
	oldVal := reflect.ValueOf(old)

	var columnVal reflect.Value
	newValue := reflect.ValueOf(new)
	direct := reflect.Indirect(newValue)

	for i := 0; i < oldVal.Len(); i++ {
		columnVal, err = findStructValByColumnKey(oldVal.Index(i), oldType.Elem(), key)
		if err != nil {
			return
		}
		if newElemType.Elem().Kind() != columnVal.Kind() {
			return errors.New(fmt.Sprintf("arrayColumn Err []%s", columnVal.Kind()))
		}

		direct.Set(reflect.Append(direct, columnVal))
	}
	return
}
func findStructValByColumnKey(curVal reflect.Value, elemType reflect.Type, key string) (columnVal reflect.Value, err error) {
	columnExist := false
	for i := 0; i < elemType.NumField(); i++ {
		curField := curVal.Field(i)
		if elemType.Field(i).Name == key {
			columnExist = true
			columnVal = curField
			continue
		}
	}
	if !columnExist {
		return columnVal, errors.New(fmt.Sprintf("key %s not found in %s's field", key, elemType))
	}
	return
}
func arrayIndexColumn(new, old interface{}, key, indexKey string) (err error) {
	newValue := reflect.ValueOf(new)
	if newValue.Elem().Kind() != reflect.Map {
		return errors.New("new must be map")
	}
	newElem := newValue.Type().Elem()
	if len(key) == 0 && newElem.Elem().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("new's elem expect struct, got %s", newElem.Elem().Kind()))
	}

	rv := reflect.ValueOf(old)
	oldType := reflect.TypeOf(old)
	elemType := oldType.Elem()

	var indexVal, columnVal reflect.Value
	direct := reflect.Indirect(newValue)
	mapReflect := reflect.MakeMap(newElem)
	newKey := newValue.Type().Elem().Key()

	for i := 0; i < rv.Len(); i++ {
		curVal := rv.Index(i)
		indexVal, columnVal, err = findStructValByIndexKey(curVal, elemType, indexKey, key)
		if err != nil {
			return
		}
		if newKey.Kind() != indexVal.Kind() {
			return errors.New(fmt.Sprintf("cant't convert %s to %s, your map'key must be %s", indexVal.Kind(), newKey.Kind(), indexVal.Kind()))
		}
		if len(key) == 0 {
			mapReflect.SetMapIndex(indexVal, curVal)
			direct.Set(mapReflect)
		} else {
			if newElem.Elem().Kind() != columnVal.Kind() {
				return errors.New(fmt.Sprintf("your map must be map[%s]%s", indexVal.Kind(), columnVal.Kind()))
			}
			mapReflect.SetMapIndex(indexVal, columnVal)
			direct.Set(mapReflect)
		}
	}
	return
}
func findStructValByIndexKey(curVal reflect.Value, elemType reflect.Type, indexKey, key string) (indexVal, columnVal reflect.Value, err error) {
	indexExist := false
	columnExist := false
	for i := 0; i < elemType.NumField(); i++ {
		curField := curVal.Field(i)
		if elemType.Field(i).Name == indexKey {
			switch curField.Kind() {
			case reflect.String, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Float64, reflect.Float32:
				indexExist = true
				indexVal = curField
			default:
				return indexVal, columnVal, errors.New("indexKey must be int float or string")
			}
		}
		if elemType.Field(i).Name == key {
			columnExist = true
			columnVal = curField
			continue
		}
	}
	if !indexExist {
		return indexVal, columnVal, errors.New(fmt.Sprintf("indexKey %s not found in %s's field", indexKey, elemType))
	}
	if len(key) > 0 && !columnExist {
		return indexVal, columnVal, errors.New(fmt.Sprintf("key %s not found in %s's field", key, elemType))
	}
	return
}

// 根据数组下标index获取string 防止越界
func GetArrIndexString(data []string, index int) string {
	if len(data) > index {
		return data[index]
	} else {
		return ""
	}
}