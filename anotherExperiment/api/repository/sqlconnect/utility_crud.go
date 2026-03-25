package sqlconnect

import (
	"api/utils"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func AddNewRow(model interface{}, tablename string) (error, error) {
	db, err := ConnectDb()

	if err != nil {
		return nil, utils.ErrorHandler(err, "db conn error")
	}

	defer db.Close()

	stmt, err := db.Prepare(GenerateInsertQuery(tablename, model))

	if err != nil {
		return nil, utils.ErrorHandler(err, "SQL prep statement err")
	}
	defer stmt.Close()
	values := getStructValues(model)

	fmt.Println("Args:", values)
	res, err := stmt.Exec(values...)

	if err != nil {
		//return nil, utils.ErrorHandler(err, "db insertion error")
		if strings.Contains(err.Error(), "Duplicate entry") {
			log.Printf("Skipping duplicate entry: %v", values[0])
			return nil, nil
		}
		return nil, utils.ErrorHandler(err, "db insertion error")
	}

	fmt.Println(res.RowsAffected())
	//lastid, err := res.LastInsertId()

	// if err != nil {
	// 	return nil, utils.ErrorHandler(err, "err getting last id")
	// }

	// fmt.Printf("job successful!  job insertd with id %d ", lastid)

	return nil, nil

}

// generic inserter
func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	paramindex := 1 // postgres way of doing things
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("dbTag", dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")

		//if dbTag != "" && dbTag != "job_id" {
		if dbTag != "" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			//placeholders += "?" mysqlway
			placeholders += fmt.Sprintf("$%d", paramindex)

			paramindex += 1
		}
	}
	//mysql way
	// fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	// return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)

	//postgresway
	fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func getStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()
	values := []interface{}{}
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")

		if dbTag != "" {
			fmt.Println("Processing ", modelValue.Field(i).Interface())
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	log.Println("Values", values)
	return values
}
