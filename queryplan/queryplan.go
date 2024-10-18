package queryplan

import (
	"errors"
	"strings"
)

type QueryPlan struct {
	Operation     string
	Database      string
	Table         string
	CreateColumns map[string]string
	SelectColumns []string
	WhereClause   WhereClause
}

type QueryPlanError struct {
	Msg string
}

var reservedKeyWords map[string]bool = map[string]bool{
	"select":   true,
	"from":     true,
	"where":    true,
	"insert":   true,
	"into":     true,
	"database": true,
	"table":    true,
	"order":    true,
	"group":    true,
	"by":       true,
}

func New(q *string) (QueryPlan, error) {
	qp := QueryPlan{}
	smallQ := strings.ToLower(*q)
	operationName, operationError := getOperationName(&smallQ)
	if operationError != nil {
		return qp, operationError
	}
	qp.Operation = operationName

	switch qp.Operation {
	case "CREATE TABLE":
		// get table name
		tableName, tableNameError := getTableName(q, &smallQ)
		if tableNameError != nil {
			return qp, tableNameError
		}
		qp.Table = tableName

		//get cols to create
		cols, colsError := getCreateColumns(q, &smallQ)
		if colsError != nil {
			return qp, colsError
		}

		qp.CreateColumns = cols
	case "CREATE DATABASE":
		dbName, dbNameErr := getCreateDatabaseName(q, &smallQ)
		if dbNameErr != nil {
			return qp, dbNameErr
		}
		qp.Database = dbName
	case "SELECT":
		cols, colsErr := getSelectColumns(q, &smallQ)
		if colsErr != nil {
			return qp, colsErr
		}
		qp.SelectColumns = cols

		tableName, tableNameError := getTableName(q, &smallQ)
		if tableNameError != nil {
			return qp, tableNameError
		}
		qp.Table = tableName

		whereClause, whereErr := getWhereClause(q, &smallQ)
		if whereErr != nil {
			return qp, whereErr
		}

		qp.WhereClause = whereClause
	default:
		return qp, errors.New("Unable to identity operation")
	}

	return qp, nil
}

func getCreateDatabaseName(q *string, smallQ *string) (string, error) {
	smallQStr := *smallQ
	if !strings.HasPrefix(smallQStr[7:], "database") {
		return "", errors.New("Unable to parse database name")
	}
	qStr := *q
	dbName := strings.TrimSpace(qStr[16:])

	return dbName, nil
}

func getOperationName(smallQ *string) (string, error) {
	var operation string = ""
	if strings.HasPrefix(*smallQ, "create") {
		smallQStr := *smallQ
		if strings.HasPrefix(smallQStr[7:], "table") {
			operation = "CREATE TABLE"
		} else if strings.HasPrefix(smallQStr[7:], "database") {
			operation = "CREATE DATABASE"
		} else {
			return "", errors.New("Unable to parse operation subtype")
		}
	} else if strings.HasPrefix(*smallQ, "select") {
		operation = "SELECT"
	} else {
		return "", errors.New("Unable to parse operation")
	}

	return operation, nil
}

func getTableName(qpTR *string, smallQ *string) (string, error) {
	ind := strings.Index(*smallQ, "from")
	if ind == -1 {
		ind = strings.Index(*smallQ, "table")
		if ind == -1 {
			return "", errors.New("Unable to parse table name")
		} else {
			ind += 6
		}
	} else {
		ind += 5
	}
	q := *qpTR
	var tableName string = ""
	for i := ind; i < len(q) && q[i:i+1] != " "; i++ {
		tableName += q[i : i+1]
	}

	return tableName, nil
}

func getSelectColumns(qPtr *string, smallQ *string) ([]string, error) {
	selectColumns := []string{}

	endIndex := strings.Index(*smallQ, "from")
	if endIndex == -1 {
		return selectColumns, errors.New("Unable to get select column names")
	}

	qStr := *qPtr
	colString := strings.Split(strings.TrimSpace(qStr[7:endIndex]), ",")

	for i, col := range colString {
		colString[i] = strings.TrimSpace(col)
	}

	return colString, nil

}

func getCreateColumns(qPtr *string, smallQ *string) (map[string]string, error) {
	cols := map[string]string{}

	startInd := strings.Index(*smallQ, "(")
	if startInd == -1 {
		return cols, errors.New("Unable to parse columns names")
	}

	endInd := strings.Index(*smallQ, ")")
	if endInd == -1 {
		return cols, errors.New("Unable to parse columns names")
	}
	q := *qPtr
	colAndTypes := strings.Split(q[startInd+1:endInd], ",")

	for _, rec := range colAndTypes {
		colAndType := (strings.Split(strings.TrimSpace(rec), " "))
		if len(colAndType) < 2 {
			return cols, errors.New("unable to parse column names")
		}
		cols[strings.TrimSpace(colAndType[0])] = strings.TrimSpace(colAndType[1])
	}

	return cols, nil

}

func getWhereClause(qPtr *string, smallQ *string) (WhereClause, error) {
	startInd := strings.Index(*smallQ, "where")
	if startInd == -1 {
		return WhereClause{}, errors.New("Unable to parse where clause")
	}
	startInd += 6
	endInd := len(*smallQ)

	for key, _ := range reservedKeyWords {
		tmp := strings.Index(*smallQ, key)
		if tmp > startInd {
			endInd = tmp
			break
		}
	}

	whereClause, whereClauseErr := ParseWhereClause(qPtr, smallQ, startInd, endInd)
	if whereClauseErr != nil {
		return WhereClause{}, errors.New("Unable to paarse where clause")
	}
	return whereClause, nil

}
