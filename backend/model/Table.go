package model

type Table struct {
	ID string `bson: "id"`
	Team bool
	Owner string //Is either a user or a team
	Name string
	Columns []Column
	Items []int
}

type Column struct {
	Index int
	Name string
	Type string
	Options []string
	Items []string
}

// func (table *Table) addRow(map[string]any) {
// 	for item := range table.Columns {
		
// 	}
// }