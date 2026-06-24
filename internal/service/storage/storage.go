package storage

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	migrationsRoot = "manifest/migrations"
	sqliteType     = "sqlite"
	pgsqlType      = "pgsql"
)

var migrationVersionAliasMap = map[string][]string{
	"000001_core_control_plane":         {"202606160002_core_control_plane"},
	"000002_remove_device_code":         {"202606230001_remove_device_code"},
	"000003_nullable_audit_timestamps":  {"202606230002_nullable_audit_timestamps"},
}

// Init 初始化数据库目录并执行启动迁移。
func Init(ctx context.Context) error {
	db := g.DB()
	if err := ensureSqliteDir(db.GetConfig().Link); err != nil {
		return err
	}
	return Migrate(ctx, db)
}

// Migrate 按当前数据库类型执行尚未应用的迁移。
func Migrate(ctx context.Context, db gdb.DB) error {
	dbType, err := databaseType(db.GetConfig())
	if err != nil {
		return err
	}
	if err := createMigrationTable(ctx, db, dbType); err != nil {
		return err
	}

	files, err := migrationFiles(dbType)
	if err != nil {
		return err
	}

	for _, file := range files {
		version := migrationVersion(file)
		applied, err := migrationApplied(ctx, db, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := applyMigration(ctx, db, version, file); err != nil {
			return err
		}
		g.Log().Infof(ctx, "数据库迁移已应用: %s", version)
	}
	return nil
}

func databaseType(config *gdb.ConfigNode) (string, error) {
	dbType := strings.ToLower(strings.TrimSpace(config.Type))
	if dbType == "" && config.Link != "" {
		dbType = strings.ToLower(strings.TrimSpace(strings.SplitN(config.Link, "::", 2)[0]))
	}
	switch dbType {
	case sqliteType, pgsqlType:
		return dbType, nil
	default:
		return "", gerror.Newf("不支持的数据库类型: %s", dbType)
	}
}

func ensureSqliteDir(link string) error {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(link)), sqliteType+"::") {
		return nil
	}
	start := strings.Index(link, "file(")
	if start < 0 {
		return nil
	}
	start += len("file(")
	end := strings.Index(link[start:], ")")
	if end < 0 {
		return gerror.New("SQLite 数据库路径格式无效")
	}
	dbPath := link[start : start+end]
	if dbPath == "" {
		return gerror.New("SQLite 数据库路径不能为空")
	}
	return gfile.Mkdir(gfile.Dir(dbPath))
}

func migrationFiles(dbType string) ([]string, error) {
	root, err := projectRoot()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, migrationsRoot, dbType)
	if !gfile.IsDir(dir) {
		return nil, gerror.Newf("数据库迁移目录不存在: %s", dir)
	}
	files, err := gfile.ScanDirFile(dir, "*.up.sql")
	if err != nil {
		return nil, gerror.Wrapf(err, "读取数据库迁移目录失败: %s", dir)
	}
	sort.Strings(files)
	return files, nil
}

func projectRoot() (string, error) {
	dir := gfile.Pwd()
	for {
		if gfile.IsFile(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", gerror.Wrap(err, "读取当前工作目录失败")
	}
	return "", gerror.Newf("无法定位项目根目录: %s", wd)
}

func migrationVersion(path string) string {
	return CanonicalMigrationVersion(strings.TrimSuffix(gfile.Basename(path), ".up.sql"))
}

// CanonicalMigrationVersion 返回迁移版本的规范名称。
func CanonicalMigrationVersion(version string) string {
	for canonical, aliases := range migrationVersionAliasMap {
		if version == canonical {
			return canonical
		}
		for _, alias := range aliases {
			if version == alias {
				return canonical
			}
		}
	}
	return version
}

// MigrationVersionAliases 返回迁移版本所有兼容名称。
func MigrationVersionAliases(version string) []string {
	canonical := CanonicalMigrationVersion(version)
	aliases := []string{canonical}
	if legacyAliases, ok := migrationVersionAliasMap[canonical]; ok {
		aliases = append(aliases, legacyAliases...)
	}
	return aliases
}

func createMigrationTable(ctx context.Context, db gdb.DB, dbType string) error {
	var sql string
	switch dbType {
	case sqliteType:
		sql = "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)"
	case pgsqlType:
		sql = "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW())"
	default:
		return gerror.Newf("不支持的数据库类型: %s", dbType)
	}
	_, err := db.Exec(ctx, sql)
	return gerror.Wrap(err, "创建数据库迁移记录表失败")
}

func migrationApplied(ctx context.Context, db gdb.DB, version string) (bool, error) {
	count, err := db.Model("schema_migrations").Ctx(ctx).
		WhereIn("version", MigrationVersionAliases(version)).
		Count()
	if err != nil {
		return false, gerror.Wrapf(err, "读取数据库迁移状态失败: %s", version)
	}
	return count > 0, nil
}

func applyMigration(ctx context.Context, db gdb.DB, version string, file string) error {
	sqlText := strings.TrimSpace(gfile.GetContents(file))
	if sqlText == "" {
		return gerror.Newf("数据库迁移文件为空: %s", file)
	}
	return db.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		for _, statement := range splitStatements(sqlText) {
			if !hasExecutableSQL(statement) {
				continue
			}
			if _, err := db.Exec(ctx, statement); err != nil {
				return gerror.Wrapf(err, "执行数据库迁移失败: %s", version)
			}
		}
		if _, err := db.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			return gerror.Wrapf(err, "写入数据库迁移记录失败: %s", version)
		}
		return nil
	})
}

func splitStatements(sqlText string) []string {
	var (
		statements []string
		builder    strings.Builder
		inSingle   bool
		escaped    bool
	)
	for _, char := range sqlText {
		builder.WriteRune(char)
		if char == '\'' && !escaped {
			inSingle = !inSingle
		}
		if char == ';' && !inSingle {
			statements = append(statements, strings.TrimSpace(builder.String()))
			builder.Reset()
		}
		escaped = char == '\\' && !escaped
		if char != '\\' {
			escaped = false
		}
	}
	if strings.TrimSpace(builder.String()) != "" {
		statements = append(statements, strings.TrimSpace(builder.String()))
	}
	return statements
}

func hasExecutableSQL(statement string) bool {
	for _, line := range strings.Split(statement, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		return true
	}
	return false
}
