package queryplan

import "fmt"

type WhereClause struct {
	Operator string
	Field    string
	Value    any
	Children []WhereClause
}

// only and, or and = allowed. Need to build better logic to parse ()
// and other clause like <=, >= and in clause.

var operations map[string]bool = map[string]bool{
	"AND": true,
	"OR":  true,
	"and": true,
	"or":  true,
}

var operators map[string]bool = map[string]bool{
	"=": true,
}

func ParseWhereClause(qPtr *string, smallQ *string, startInd int, endInd int) (WhereClause, error) {
	token := ""
	q := *qPtr
	var root *WhereClause
	cur := WhereClause{}
	for i := startInd; i < endInd; i++ {
		if q[i:i+1] == " " {
			fmt.Println(token)
			_, okOptr := operators[token]
			_, okOpn := operations[token]

			if okOpn {
				if token == "AND" || token == "and" {
					if root != nil {
						cur.Children = append(cur.Children, *root)
					}
					root = &cur
				} else if token == "OR" || token == "or" {
					if root != nil {
						if len((*root).Children) < 2 {
							parent := WhereClause{}
							parent.Children = append(parent.Children, *root, cur)
							root = &parent
						} else {
							(*root).Children = append(root.Children, cur)
						}
					} else {
						cur.Children = append(cur.Children, *root)
						root = &cur
					}
				}
				cur = WhereClause{}
			} else if okOptr {
				cur.Operator = token
			} else if cur.Operator == "" {
				cur.Field = token
			} else {
				cur.Value = token
			}

			token = ""
			continue
		}
		token += q[i : i+1]
	}
	if root == nil {
		*root = WhereClause{}
	}
	return *root, nil
}
