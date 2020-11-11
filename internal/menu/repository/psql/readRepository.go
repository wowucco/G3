package psql

import (
	"context"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/wowucco/G3/internal/entity"
	"github.com/wowucco/nestedset"
)

const tableNameGroup = "shop_group"

type MenuReadRepository struct {
	db *dbx.DB
}

func NewMenuReadRepository(db *dbx.DB) *MenuReadRepository {
	return &MenuReadRepository{
		db: db,
	}
}

func (r MenuReadRepository) RootMenuItemWithDepth(ctx context.Context, depth int) (*entity.MenuItem, error) {

	var (
		rows []MenuItem
	)

	err := r.db.Select("g.*").
		From(tableWithAlias(tableNameGroup, "g")).
		OrderBy("lft asc").
		All(&rows)

	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("menu items not found")
	}

	ns := toNSBranch(rows)

	menu := rows[0]

	node := ns.FindById(menu.Id)

	if node == nil {
		return nil, fmt.Errorf("menu with id %d not found", menu.Id)
	}

	menuItem := entity.MenuItem{
		Menu: entity.Menu{
			Id:          int(menu.Id),
			Name:        menu.Name,
		},
		Parent:      nil,
		HasParent:   false,
	}

	if menu.Description.Valid {
		menuItem.Description = menu.Description.String
	}

	if menu.Image.Valid {
		menuItem.Image = menu.Image.String
	}

	if len(ns.Branch(node)) > 1 {
		menuItem.HasChildren = true
	}

	menuItemMapWithIdKey := make(map[int64]*MenuItem, len(rows))

	for _, v := range rows {
		menuItemMapWithIdKey[v.GetId()] = &v
	}

	if menuItem.HasChildren && depth > 0 {
		menuItem.Children = getNodeChildrenFromNS(node, ns, int64(depth), menuItemMapWithIdKey)
	}

	return &menuItem, nil
}

func (r MenuReadRepository) MenuItemWithDepthById(ctx context.Context, id, depth int, needParent bool) (*entity.MenuItem, error) {

	exist, err := r.Exist(ctx, id)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("menu %d not found", id)
	}

	var (
		rows []MenuItem
		menu MenuItem
	)

	err = r.db.Select("g.*").
		From(tableWithAlias(tableNameGroup, "g")).
		OrderBy("lft asc").
		All(&rows)

	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("menu items not found")
	}

	ns := toNSBranch(rows)

	idInt64 := int64(id)

	menuItemMapWithIdKey := make(map[int64]*MenuItem, len(rows))

	for _, v := range rows {
		menuItemMapWithIdKey[v.GetId()] = &v

		if v.GetId() == idInt64 {
			menu = v
		}
	}

	node := ns.FindById(menu.Id)

	if node == nil {
		return nil, fmt.Errorf("menu with id %d not found", menu.Id)
	}

	menuItem := entity.MenuItem{
		Menu: entity.Menu{
			Id:          int(menu.Id),
			Name:        menu.Name,
		},
		Parent:      nil,
	}

	if menu.Description.Valid {
		menuItem.Description = menu.Description.String
	}

	if menu.Image.Valid {
		menuItem.Image = menu.Image.String
	}

	if len(ns.Branch(node)) > 1 {
		menuItem.HasChildren = true
	}

	if menuItem.HasChildren && depth > 0 {
		menuItem.Children = getNodeChildrenFromNS(node, ns, int64(depth), menuItemMapWithIdKey)
	}

//HasParent:   false,
	if needParent {
		nodeParent := getParentFromNS(node, ns, menuItemMapWithIdKey)

		if nodeParent != nil {
			menuItem.HasParent = true
			menuItem.Parent = nodeParent
		}
	} else {
		nodeParent := ns.Parent(node)

		if nodeParent != nil {
			menuItem.HasParent = true
		}
	}

	return &menuItem, nil
}

func (r MenuReadRepository) Exist(ctx context.Context, id int) (bool, error) {

	var count int

	err := r.db.Select("COUNT(*)").From(tableWithAlias(tableNameGroup, "g")).
		Where(dbx.NewExp("g.id={:id}", dbx.Params{"id": id})).
		Row(&count)

	return count > 0, err
}

func tableWithAlias(tableName, alias string) string {
	return tableName + " " + alias
}

func toNSBranch(rows []MenuItem) *nestedset.NestedSet  {

	var (
		parent nestedset.NodeInterface
		last nestedset.NodeInterface
	)

	rootNode := &rows[0]

	ns := nestedset.NewNestedSet(rootNode)

	parent = rootNode
	last   = rootNode

	for _, m := range rows[1:] {

		menu := m

		if menu.GetLevel() > last.GetLevel() {
			_ = ns.Add(&menu, last)
			parent = last
		} else if menu.GetLevel() == last.GetLevel() {
			_ = ns.Add(&menu, parent)
		} else {
			for menu.GetLevel() <= parent.GetLevel() {
				parent = ns.Parent(parent)
			}
			_ = ns.Add(&menu, parent)
		}

		last = &menu
	}

	return ns
}

func getNodeChildrenFromNS(node nestedset.NodeInterface, ns *nestedset.NestedSet, depth int64, rowsWithIdKey map[int64]*MenuItem) []*entity.ChildrenMenuItem {

	var (
		children []*entity.ChildrenMenuItem
		maxDepth, currentDepth int64
	)

	maxDepth = node.GetLevel() + depth
	currentDepth = node.GetLevel() + 1

	for _, v := range ns.Branch(node) {

		if v.GetId() == node.GetId() || v.GetLevel() != currentDepth || v.GetLevel() > maxDepth {
			continue
		}

		item := entity.ChildrenMenuItem{
			Menu:     entity.Menu{
				Id:          int(v.GetId()),
				Name:        v.GetName(),
			},
			HasChildren: v.GetRight() > v.GetLeft() + 1,
		}


		if val, ok := rowsWithIdKey[v.GetId()]; ok {
			if val.Description.Valid {
				item.Description = val.Description.String
			}

			if val.Image.Valid {
				item.Image = val.Image.String
			}
		}

		if item.HasChildren && v.GetLevel() < maxDepth {
			item.Children = getNodeChildrenFromNS(v, ns, depth - 1, rowsWithIdKey)
		}

		children = append(children, &item)
	}

	return children
}

func getParentFromNS(node nestedset.NodeInterface, ns *nestedset.NestedSet, rowsWithIdKey map[int64]*MenuItem) *entity.ParentMenuItem {
	
	var (
		parent entity.ParentMenuItem
	)
	
	parentNode := ns.Parent(node)
	
	if parentNode == nil {
		return nil
	}

	parent = entity.ParentMenuItem{
		Menu:   entity.Menu{
			Id:          int(parentNode.GetId()),
			Name:        parentNode.GetName(),
		},
	}

	if val, ok := rowsWithIdKey[parentNode.GetId()]; ok {
		if val.Description.Valid {
			parent.Description = val.Description.String
		}

		if val.Image.Valid {
			parent.Image = val.Image.String
		}
	}

	parent.Parent = getParentFromNS(parentNode, ns, rowsWithIdKey)

	if parent.Parent != nil {
		parent.HasParent = true
	}

	return &parent
}
