// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PluginsDao is the data access object for the table plugins.
type PluginsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PluginsColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PluginsColumns defines and stores column names for the table plugins.
type PluginsColumns struct {
	Id            string //
	Name          string //
	Type          string //
	Runtime       string //
	Protocol      string //
	Status        string //
	ActiveVersion string //
	Source        string //
	Description   string //
	Enabled       string //
	InstalledAt   string //
	UpdatedAt     string //
	DeletedAt     string //
}

// pluginsColumns holds the columns for the table plugins.
var pluginsColumns = PluginsColumns{
	Id:            "id",
	Name:          "name",
	Type:          "type",
	Runtime:       "runtime",
	Protocol:      "protocol",
	Status:        "status",
	ActiveVersion: "active_version",
	Source:        "source",
	Description:   "description",
	Enabled:       "enabled",
	InstalledAt:   "installed_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewPluginsDao creates and returns a new DAO object for table data access.
func NewPluginsDao(handlers ...gdb.ModelHandler) *PluginsDao {
	return &PluginsDao{
		group:    "default",
		table:    "plugins",
		columns:  pluginsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PluginsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PluginsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PluginsDao) Columns() PluginsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PluginsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PluginsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *PluginsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
