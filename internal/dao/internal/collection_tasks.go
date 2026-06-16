// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CollectionTasksDao is the data access object for the table collection_tasks.
type CollectionTasksDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  CollectionTasksColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// CollectionTasksColumns defines and stores column names for the table collection_tasks.
type CollectionTasksColumns struct {
	Id              string //
	Name            string //
	DeviceId        string //
	SouthPluginId   string //
	ReportMode      string //
	Status          string //
	RuleHitRate     string //
	Rate            string //
	LastCollectedAt string //
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// collectionTasksColumns holds the columns for the table collection_tasks.
var collectionTasksColumns = CollectionTasksColumns{
	Id:              "id",
	Name:            "name",
	DeviceId:        "device_id",
	SouthPluginId:   "south_plugin_id",
	ReportMode:      "report_mode",
	Status:          "status",
	RuleHitRate:     "rule_hit_rate",
	Rate:            "rate",
	LastCollectedAt: "last_collected_at",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewCollectionTasksDao creates and returns a new DAO object for table data access.
func NewCollectionTasksDao(handlers ...gdb.ModelHandler) *CollectionTasksDao {
	return &CollectionTasksDao{
		group:    "default",
		table:    "collection_tasks",
		columns:  collectionTasksColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CollectionTasksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CollectionTasksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CollectionTasksDao) Columns() CollectionTasksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CollectionTasksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CollectionTasksDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CollectionTasksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
