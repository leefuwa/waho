package models

import (
	sql2 "database/sql"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strings"
	"waho/comm"
	"waho/register"
)

type models struct {
	TableName string
	Data interface{}
	Fields []string
	Pk string
	Prefix string
	DefaultSoftDelete bool
	DefaultSoftDeleteName string
	tableClass string
	lastSql string
	whereFields []string
	txOrm orm.TxOrmer
	error error
}

type WhereMap map[string]interface{}
type WhereArr [][]string

func defaultPk() string {
	return "id"
}

func NewModels(data interface{}) models {
	model := models{Data: data}
	model.setTableDetail()
	model.DefaultSoftDelete = true
	model.DefaultSoftDeleteName = "is_delete"
	return model
}

func MakeWhereMap() (where WhereMap) {
	where = make(WhereMap)

	return
}
func MakeWhereArr() (where WhereArr) {
	where = make(WhereArr, 0)
	return
}

func NewTxOrmer() (orm.TxOrmer, error) {
	txOrm, err := orm.NewOrm().Begin()
	return txOrm, err
}

func GetWhereStructMap() map[string]interface{} {
	where := make(map[string]interface{})

	return where
}

func GetWhereStructArr() [][]string {
	where := make([][]string, 0)

	return where
}

func (model *models) StartTrans(ormer orm.TxOrmer) *models {
	model.txOrm = ormer
	return model
}

func (model *models) Commit() error {
	err := model.txOrm.Commit()
	return err
}

func (model *models) Rollback() error {
	err := model.txOrm.Rollback()
	return err
}

func (model *models) GetTableName() string {
	return model.TableName
}

func (model *models) GetLastSql() string {
	return model.lastSql
}

func (model *models) SetWhereField(fields []string) *models {
	model.whereFields = fields
	return model
}

// 基础where
//   params := [][]string{0:[]string{"id", ">", "1"},{0:[]string{"name", "like", "%test%"}}}
//     OR params := map[string]interface{}{"id": 1, "name": "test",}
//   sql, where, diff, error = .GetWhere(params)
func (model *models) GetWhere(params interface{}) (string, []string, map[string][]string, error) {
	// 设置params条件
	paramsKey := reflect.TypeOf(params)
	paramsKeyType := paramsKey.String()
	paramsVal := reflect.ValueOf(params)
	var arr2Where [][]string
	if paramsKeyType == "[][]string" || strings.Contains(paramsKeyType, "WhereArr")  {
		paramsValLen := paramsVal.Len()
		for i := 0; i < paramsValLen; i++ {
			var paramsValTemp []string
			valueTemp := reflect.ValueOf(paramsVal.Index(i).Interface())
			valueTempLen := valueTemp.Len()
			for vi := 0; vi < valueTempLen; vi++ {
				paramsValStr, _ := comm.GetValueString(valueTemp.Index(vi).Interface())
				paramsValTemp = append(paramsValTemp, paramsValStr)
			}
			arr2Where = append(arr2Where, paramsValTemp)
		}
	} else if paramsKeyType == "map[string]interface {}" || strings.Contains(paramsKeyType, "WhereMap") {
		for _, mapkey := range paramsVal.MapKeys() {
			temp, _ := comm.GetValueString(paramsVal.MapIndex(mapkey).Interface())
			paramsValTypeString := reflect.TypeOf(paramsVal.MapIndex(mapkey).Interface()).String()
			if paramsValTypeString == "[]interface {}" || paramsValTypeString == "[]string"{
				tempVal := reflect.ValueOf(paramsVal.MapIndex(mapkey).Interface())
				tempValLen := tempVal.Len()
				if tempValLen == 1 {
					arr2Where = append(arr2Where, []string{
						tempVal.Index(0).Interface().(string),
					})
				} else if tempValLen == 2 {
					tempValLen2Val_0, _ := comm.GetValueString(tempVal.Index(0).Interface())
					tempValLen2Val_1, _ := comm.GetValueString(tempVal.Index(1).Interface())
					arr2Where = append(arr2Where, []string{
						mapkey.String(), tempValLen2Val_0, tempValLen2Val_1,
					})
				} else {
					tempValLen3Val_0, _ := comm.GetValueString(tempVal.Index(0).Interface())
					tempValLen3Val_1, _ := comm.GetValueString(tempVal.Index(1).Interface())
					tempValLen3Val_2, _ := comm.GetValueString(tempVal.Index(2).Interface())
					arr2Where = append(arr2Where, []string{
						tempValLen3Val_0, tempValLen3Val_1, tempValLen3Val_2,
					})
				}
			} else {
				arr2Where = append(arr2Where, []string{
					mapkey.String(), "=" , temp,
				})
			}
		}
	}  else if paramsKeyType == "int" || paramsKeyType == "string"{
		id, _ := comm.GetValueString(paramsVal.Interface())
		arr2Where = append(arr2Where, []string{
			model.Pk, "=" , id,
		})
	} else if strings.Contains(paramsKeyType, "models"){
		if len(model.whereFields) == 0 {
			id, _ := comm.GetValueString(paramsVal.FieldByName(comm.CamelString(model.Pk)).Interface())
			arr2Where = append(arr2Where, []string{
				model.Pk, "=" , id,
			})
		} else {
			for i := 0; i < paramsKey.NumField(); i++ {
				key := paramsKey.Field(i).Tag.Get("json")
				if key == "" {
					key = comm.SnakeString(paramsKey.Field(i).Name)
				}
				tempVal, _ := comm.GetValueString(paramsVal.Field(i).Interface())
				arr2Where = append(arr2Where, []string{
					key, "=" , tempVal,
				})
			}
		}
	} else {
		return "", nil, nil, errors.New("Base.GetWhere 获取参数条件存在问题,只能传递[][]string,map[string]interface{},int,string")
	}

	if model.DefaultSoftDelete && model.isField(model.DefaultSoftDeleteName){
		arr2Where = append(arr2Where, []string{
			model.DefaultSoftDeleteName, "=" , "0",
		})
	}

	sql, where, diff, err := model.setBaseSqlWhere(arr2Where)

	return sql, where, diff, err
}

func (model *models) GetInfo(params interface{}) (bool, error) {
	var err error
	if params == nil {
		params, err = model.GetFieldValue(model.Pk)

		if err != nil {
			return false, err
		}
	}


	whereSql, where, diff, err := model.GetWhere(params)
	if err != nil {
		return false, err
	}

	var orderBy string
	orderByLen := len(diff["order_by"])
	if orderByLen > 0 {
		if orderByLen == 1 {
			orderBy = " ORDER BY " + diff["order_by"][0] + " DESC"
		} else {
			orderBy = " ORDER BY " + diff["order_by"][0] + " " + diff["order_by"][1]
		}
	} else {
		if model.isField("o") {
			orderBy = " ORDER BY o asc"
		} else {
			orderBy = " ORDER BY " + model.Pk + " DESC"
		}
	}

	sql := "SELECT * FROM " + model.TableName + " WHERE TRUE " + whereSql + orderBy + " LIMIT 1"
	model.lastSql = replaceBindParamsSql(sql, where)

	modelOrm := orm.NewOrm()
	err = modelOrm.Raw(sql, where).QueryRow(model.Data)
	log.Info("Get Info Sql :  " , model.lastSql)

	var exist bool
	if err != nil {
		if err != orm.ErrNoRows {
			log.Panic("Warning Err . ",model.GetLastSql())
			return false, err
		}
		exist = false
	} else {
		exist = true
	}

	return exist, nil
}

func (model *models) GetCount(params interface{}) (int64, error) {
	if params == nil{
		params = [][]string{}
	}

	whereSql, where, _, err := model.GetWhere(params)
	if err != nil {
		return 0, err
	}

	sql := "SELECT count(*) FROM " + model.TableName + " WHERE TRUE " + whereSql
	model.lastSql = replaceBindParamsSql(sql, where)
	log.Info("Get Count Sql :  " , model.lastSql)

	var count int64
	modelOrm := orm.NewOrm()
	err = modelOrm.Raw(sql, where).QueryRow(&count)

	if err != nil {
		log.Panic("Warning Err . ",model.GetLastSql())
		return 0, err
	}

	return count, nil
}

func (model *models) GetList(params interface{}, page int, size int) (int64, error) {
	if params == nil{
		params = [][]string{}
	}

	whereSql, where, diff, err := model.GetWhere(params)
	if err != nil {
		return 0, err
	}

	var orderBy string
	orderByLen := len(diff["order_by"])
	if orderByLen > 0 {
		if orderByLen == 1 {
			orderBy = " ORDER BY " + diff["order_by"][0] + " DESC"
		} else {
			orderBy = " ORDER BY " + diff["order_by"][0] + " " + diff["order_by"][1]
		}
	} else {
		if model.isField("o") {
			orderBy = " ORDER BY o asc"
		} else {
			orderBy = " ORDER BY " + model.Pk + " DESC"
		}
	}

	count, err := model.GetCount(params)
	if err != nil || count == 0{
		return 0, nil
	}

	// 总页数
	pageMax := int(math.Ceil(float64(comm.Max(count, 1)) / float64(size)))

	// 当前页
	page = comm.Min(page, pageMax)

	// 当前条数
	startNum := (page - 1) * size

	sql := "SELECT * FROM " + model.TableName + " WHERE TRUE " + whereSql + orderBy + " limit " + comm.ToSting(startNum) + "," + comm.ToSting(size)
	model.lastSql = replaceBindParamsSql(sql, where)
	log.Info("Get List Sql :  " , model.lastSql)

	modelOrm := orm.NewOrm()
	_, err = modelOrm.Raw(sql, where).QueryRows(model.Data)

	if err != nil {
		log.Panic("Warning Err . ",model.GetLastSql())
		return 0, err
	}

	return count, nil
}

func (model *models) AddObj(save ...interface{}) (int64, error) {
	// 如果传值为1个
	saveLen := len(save)
	var txOrm orm.TxOrmer
	if saveLen == 1 {
		saveTypeString := reflect.TypeOf(save[0]).String()
		if strings.Contains(saveTypeString, "models") {
			if saveTypeString[:1] != "*" {
				log.Panic("请用传递类型")
				return 0, errors.New("请用传递类型！")
			}
			model.Data = save[0]
		} else if strings.Contains(saveTypeString, "TxOrmer") {
			txOrm = save[0].(orm.TxOrmer)
		}
	} else if saveLen > 1 {
		saveTypeString := reflect.TypeOf(save[0]).String()
		if strings.Contains(saveTypeString, "models") {
			if saveTypeString[:1] != "*" {
				log.Panic("请用传递类型")
				return 0, errors.New("请用传递类型！")
			}
			model.Data = save[0]
		}

		txOrm = save[1].(orm.TxOrmer)
	}

	if txOrm == nil {
		txOrm = model.txOrm
	}

	pkValStr, _ := model.GetFieldValue(model.Pk)
	if pkValStr == "0" {
		model.SetFieldValue(model.Pk, comm.GetSnowflakeId())
	}

	var id int64
	var err error
	if txOrm != nil {
		id, err = txOrm.Insert(model.Data)
	} else {
		thisOrm := orm.NewOrm()
		id, err = thisOrm.Insert(model.Data)
	}

	return id, err
}

func (model *models) Add(save map[string]interface{}, txOrmer ...orm.TxOrmer) (int64, error) {
	var sqlKey,sqlWhere string
	var sqlValue []string
	saveVal := reflect.ValueOf(save)
	var isExistencePk bool
	for _, mapkey := range saveVal.MapKeys() {
		if !isExistencePk && mapkey.String() == model.Pk {
			isExistencePk = true
		}
		sqlKey += "," + mapkey.String()
		sqlWhere += ",?"
		sqlValue = append(sqlValue, comm.ToSting(saveVal.MapIndex(mapkey).Interface()))
	}
	var txOrm orm.TxOrmer
	if len(txOrmer) > 0 {
		txOrm = txOrmer[0].(orm.TxOrmer)
	}

	if sqlKey == "" || sqlWhere == "" || len(sqlValue) == 0 {
		return 0, errors.New("Base.Add 获取参数存在问题！")
	}
	if !isExistencePk {
		sqlKey += "," + model.Pk
		sqlWhere += ",?"
		sqlValue = append(sqlValue, config.ToString(comm.GetSnowflakeId()))
	}
	sqlKey = sqlKey[1:]
	sqlWhere = sqlWhere[1:]
	sql := "INSERT INTO " + model.TableName + " (" + sqlKey + ") VALUES (" + sqlWhere + ")"
	model.lastSql = replaceBindParamsSql(sql, sqlValue)

	if txOrm == nil {
		txOrm = model.txOrm
	}

	var result sql2.Result
	var err error
	if txOrm != nil {
		result, err = txOrm.Raw(sql, sqlValue).Exec()
	} else {
		thisOrm := orm.NewOrm()
		result, err = thisOrm.Raw(sql, sqlValue).Exec()
	}

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()

	return id, err
}

func (model *models) UpdateObj(save ...interface{}) (int64, error) {
	// 如果传值为1个
	saveLen := len(save)
	var txOrm orm.TxOrmer
	if saveLen == 1 {
		saveTypeString := reflect.TypeOf(save[0]).String()
		if strings.Contains(saveTypeString, "models") {
			if saveTypeString[:1] != "*" {
				log.Panic("请用传递类型")
				return 0, errors.New("请用传递类型！")
			}
			model.Data = save[0]
		} else if strings.Contains(saveTypeString, "TxOrmer") {
			txOrm = save[0].(orm.TxOrmer)
		}
	} else if saveLen > 1 {
		saveTypeString := reflect.TypeOf(save[0]).String()
		if strings.Contains(saveTypeString, "models") {
			if saveTypeString[:1] != "*" {
				log.Panic("请用传递类型")
				return 0, errors.New("请用传递类型！")
			}
			model.Data = save[0]
		}

		txOrm = save[1].(orm.TxOrmer)
	}

	if txOrm == nil {
		txOrm = model.txOrm
	}

	var count int64
	var err error
	if txOrm != nil {
		count, err = txOrm.Update(model.Data)
	} else {
		thisOrm := orm.NewOrm()
		count, err = thisOrm.Update(model.Data)
	}

	return count, err
}

func (model *models) Update(save ...map[string]interface{}) (int64, error)  {
	var sqlKeyWhere string
	var sqlValue []string
	saveLen := reflect.ValueOf(save).Len()
	if saveLen == 0 {
		return 0, errors.New("Base.Update 参数传递存在问题，传递的参数只能1和2个！")
	} else if saveLen == 1 {
		pkVal, err := model.GetFieldValue(model.Pk)
		if err != nil {
			return 0, err
		}
		sqlKeyWhere += " AND " + model.Pk + " = ?"
		sqlValue = append(sqlValue, pkVal)
	} else {
		sql, where, _, err := model.GetWhere(save[1])
		if err != nil {
			return 0, err
		}
		sqlKeyWhere = sql
		sqlValue = append(sqlValue, where...)
	}

	if sqlKeyWhere == "" {
		return 0, errors.New("Base.Update 参数传递存在问题，where条件是空值")
	}

	var updateKey string
	var updateValue []string
	saveVal := reflect.ValueOf(save[0])
	for _, mapKey := range saveVal.MapKeys() {
		updateKey += "," + mapKey.String() + " = ?"
		updateValue = append(updateValue, comm.ToSting(saveVal.MapIndex(mapKey).Interface()))
	}

	if updateKey == "" || len(updateValue) == 0{
		return 0, errors.New("Base.Update 参数传递存在问题，保存值是空值")
	}

	updateKey = updateKey[1:]
	sqlKeyWhere = sqlKeyWhere[4:]

	sql := "UPDATE " + model.GetTableName() + " SET " + updateKey + " WHERE " + sqlKeyWhere
	where := append(updateValue, sqlValue...)
	model.lastSql = replaceBindParamsSql(sql, where)
	log.Info("Update Sql :  " , model.lastSql)

	var result sql2.Result
	var err error
	if model.txOrm != nil {
		result, err = model.txOrm.Raw(sql, where).Exec()
	} else {
		thisOrm := orm.NewOrm()
		result, err = thisOrm.Raw(sql, where).Exec()
	}

	if err != nil {
		log.Error("Update Error：", model.GetLastSql())
		return 0, err
	}

	count, err := result.RowsAffected()

	if err != nil {
		log.Error("Update Error：", model.GetLastSql())
		return 0, err
	}

	return count, nil
	
}

func replaceBindParamsSql(sql string, params []string) string {
	for _, v := range params{
		sql = strings.Replace(sql, "?", "\""+v+"\"", 1)
	}

	return sql
}

func (model *models) setBaseSqlWhere(params [][]string) (string, []string, map[string][]string, error) {
	var sql string
	var where []string
	diff := make(map[string][]string)

	for _, _v := range params{
		vLen := len(_v)

		if vLen == 0 || vLen > 3 {
			return "", nil, nil, errors.New("Base.setBaseSqlWhere 参数值存在问题，Len只能是1/2/3")
		}

		if vLen == 1{
			sql += " AND " + _v[0]
			continue
		}

		if !model.isWhereField(_v[0]) {
			diff[_v[0]] = _v[1:]
			continue
		}

		if vLen == 2 {
			sql += " AND " + _v[0] + " = ?"
			where = append(where, _v[1])
		} else if vLen == 3 {
			lowerComparison := strings.ToLower(_v[1])
			if lowerComparison == "in" || lowerComparison == "not in"{
				value := strings.Split(_v[2],",")
				strBindSet := strings.Join(comm.ArrayFill("?", len(value)),",")
				sql += " AND " + _v[0] + " " + _v[1] + "(" + strBindSet + ")"
				where = append(where, value...)
			} else {
				sql += " AND " + _v[0] + " " + _v[1] + " ?"
				where = append(where, _v[2])
			}
		}
	}


	return sql, where, diff , nil
}

func (model *models) isWhereField(field string) bool {
	if len(model.whereFields) > 0 {
		for i := 0; i < len(model.whereFields); i++ {
			if field == model.whereFields[i] {
				return true
			}
		}

		return false
	} else {
		return model.isField(field)
	}
}

func (model *models) isField(field string) bool {
	for i := 0; i < len(model.Fields); i++ {
		if field == model.Fields[i] {
			return true
		}
	}

	return false
}

func (model *models) setTableDetail()  {
	model.setDbPrefix()
	model.setTalesName()
	model.setFields()
}

func (model *models) setDbPrefix() {
	model.Prefix = GetDbPrefix()
}

func (model *models) setTalesName() {
	table := model.Data
	key := reflect.TypeOf(table)
	if key.Kind() == reflect.Ptr {
		key = key.Elem()
	}

	//fmt.Println("key.String()",key.String())
	tableStr := key.Name()
	if tableStr == "" {
		tableStr = strings.Replace(key.String(), "[]models.", "", 1)
	}

	model.tableClass = tableStr
	model.TableName = model.Prefix + comm.SnakeString(tableStr)
}

func (model *models) setFields()  {
	table := model.Data
	var fields interface{}
	tableType := reflect.TypeOf(table).Elem().String()[:2]
	if tableType == "[]" {
		fields = reflect.New(reflect.TypeOf(table).Elem().Elem()).Interface()
	} else {
		fields = table
	}


	var fieldsKey []string
	var fieldPk string
	var filedIsId bool
	fieldsType := reflect.TypeOf(fields).Elem()
	for i := 0; i < fieldsType.NumField(); i++ {
		key := fieldsType.Field(i).Tag.Get("json")
		if key == "" {
			key = comm.SnakeString(fieldsType.Field(i).Name)
		}
		if key == defaultPk() {
			filedIsId = true
		}
		fieldsKey = append(fieldsKey, key)

		ormArr := strings.Split( fieldsType.Field(i).Tag.Get("orm"), ";")
		if fieldPk == "" {
			for _, ormArrV := range ormArr{
				if ormArrV == "pk" {
					fieldPk = fieldsType.Field(i).Tag.Get("json")
					if fieldPk == "" {
						key = comm.SnakeString(fieldsType.Field(i).Name)
					}
					break
				}
			}
		}
	}

	if fieldPk == "" {
		if !filedIsId {
			errors.New(model.tableClass + "尚未设置主键")
		}

		fieldPk = defaultPk()
	}


	model.Pk = fieldPk
	model.Fields = fieldsKey
}

func (model *models) GetFieldValue(field string) (string, error) {
	table := model.Data
	var fields interface{}
	tableType := reflect.TypeOf(table).Elem().String()[:2]
	if tableType == "[]" {
		fields = reflect.New(reflect.TypeOf(table).Elem().Elem()).Interface()
	} else {
		fields = table
	}

	fieldsType := reflect.TypeOf(fields).Elem()
	fieldsVal := reflect.ValueOf(fields).Elem()
	for i := 0; i < fieldsType.NumField(); i++ {
		key := fieldsType.Field(i).Tag.Get("json")
		if key == field {
			return comm.ToSting(fieldsVal.Field(i).Interface()), nil
		}
	}

	return "", errors.New("Base.GetFieldValue 获取Val失败，当前不存在【" + field + "】或尚未设置JsonKey")
}

func (model *models) SetFieldValue(key string, val interface{})  {

	table := model.Data
	var fields interface{}
	tableType := reflect.TypeOf(table).Elem().String()[:2]
	if tableType == "[]" {
		fields = reflect.New(reflect.TypeOf(table).Elem().Elem()).Interface()
	} else {
		fields = table
	}

	fieldsType := reflect.TypeOf(fields).Elem()
	fieldsVal := reflect.ValueOf(fields).Elem()
	for i := 0; i < fieldsType.NumField(); i++ {
		if key == fieldsType.Field(i).Tag.Get("json") {
			fieldKey := fieldsVal.Field(i).Type().String()
			if fieldKey == "int" || fieldKey == "int64" {
				fieldsVal.Field(i).SetInt(comm.ToInt64(val))
			} else if fieldKey == "float64" {
				fieldsVal.Field(i).SetFloat(val.(float64))
			} else if fieldKey == "string" {
				fieldsVal.Field(i).SetString(comm.ToSting(val))
			}
			break
		}
	}
}

func GetDbPrefix() string {
	return register.GetConf.Section("mysql").GetDefault("dbprefix", "")
}
