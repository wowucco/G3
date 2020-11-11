package entity

type Menu struct {
	Id           int
	Name 	     string
	Description  string
	Image        string
	Meta 		 Meta
}

type MenuItem struct {
	Menu
	Parent *ParentMenuItem
	Children []*ChildrenMenuItem
	HasChildren  bool
	HasParent  	 bool
}

type ParentMenuItem struct {
	Menu
	HasParent bool
	Parent *ParentMenuItem
}

type ChildrenMenuItem struct {
	Menu
	HasChildren bool
	Children []*ChildrenMenuItem
}