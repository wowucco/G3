package psql

import "database/sql"

type Meta struct {

}

type MenuItem struct {
	Id 			int64	`db:"id"`
	Name 		string	`db:"name"`
	Description sql.NullString	`db:"description"`
	Image 		sql.NullString	`db:"image"`
	Meta 		Meta
	Left 		int64	`db:"lft"`
	Right 		int64	`db:"rgt"`
	Depth 		int64	`db:"depth"`
}

func (m MenuItem) GetId() int64  {
	return m.Id
}

func (m MenuItem) GetName() string {
	return m.Name
}

func (m MenuItem) GetLevel() int64 {
	return m.Depth
}

func (m MenuItem) GetLeft() int64 {
	return m.Left
}

func (m MenuItem) GetRight() int64 {
	return m.Right
}

func (m *MenuItem) SetLevel(level int64) {
	m.Depth = level
}

func (m *MenuItem) SetLeft(left int64) {
	m.Left = left
}

func (m *MenuItem) SetRight(right int64) {
	m.Right = right
}
