// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PluginVersionsDao is the data access object for the table plugin_versions.
type PluginVersionsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  PluginVersionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// PluginVersionsColumns defines and stores column names for the table plugin_versions.
type PluginVersionsColumns struct {
	Id               string //
	PluginId         string //
	Version          string //
	PackagePath      string //
	ManifestJson     string //
	PermissionsJson  string //
	CapabilitiesJson string //
	Checksum         string //
	Active           string //
	InstalledAt      string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// pluginVersionsColumns holds the columns for the table plugin_versions.
var pluginVersionsColumns = PluginVersionsColumns{
	Id:               "id",
	PluginId:         "plugin_id",
	Version:          "version",
	PackagePath:      "package_path",
	ManifestJson:     "manifest_json",
	PermissionsJson:  "permissions_json",
	CapabilitiesJson: "capabilities_json",
	Checksum:         "checksum",
	Active:           "active",
	InstalledAt:      "installed_at",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewPluginVersionsDao creates and returns a new DAO object for table data access.
func NewPluginVersionsDao(handlers ...gdb.ModelHandler) *PluginVersionsDao {
	return &PluginVersionsDao{
		group:    "default",
		table:    "plugin_versions",
		columns:  pluginVersionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PluginVersionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PluginVersionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PluginVersionsDao) Columns() PluginVersionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PluginVersionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PluginVersionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PluginVersionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
