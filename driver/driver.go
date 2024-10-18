package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TanmayBhosale/goDB/queryplan"
)

func ProcessQuery(q *string) map[string]any {
	qp, err := queryplan.New(q)
	if err != nil {
		return map[string]any{
			"error": err.Error(),
		}
	}

	if qp.Operation == "CREATE" {
		createErr := createTable(qp)
		if createErr != nil {
			return map[string]any{
				"error": createErr.Error(),
			}
		}
	}
	return map[string]any{
		"operation":  qp.Operation,
		"table":      qp.Table,
		"selectCols": qp.SelectColumns,
		"createCols": qp.CreateColumns,
		"database":   qp.Database,
		"where":      qp.WhereClause,
	}
}

func createTable(qp queryplan.QueryPlan) error {
	_, colJsonErr := json.Marshal(qp.CreateColumns)
	if colJsonErr != nil {
		fmt.Println("error!!")
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Err")
	}
	driverDir := filepath.Dir(filename)

	fmt.Println(driverDir)

	pathSlice := strings.Split(driverDir, "\\")

	dir := strings.Join(pathSlice[:len(pathSlice)-1], "/")

	dbErr := os.Mkdir(dir+"/db", 0666)

	if dbErr != nil {
		errStr := dbErr.Error()
		if !strings.Contains(errStr, "Cannot create a file when that file already exists") {
			return errors.New("Uaable to get db location")
		}
	}

	return nil
}
