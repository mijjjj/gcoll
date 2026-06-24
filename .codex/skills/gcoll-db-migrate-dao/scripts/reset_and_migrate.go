package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"gopkg.in/yaml.v3"

	storagesvc "github.com/mijjjj/gcoll/internal/service/storage"
)

const (
	migrationsRoot = "manifest/migrations"
	sqliteType     = "sqlite"
	pgsqlType      = "pgsql"
)

// configFile 描述最小化数据库配置结构。
type configFile struct {
	Database map[string]gdb.ConfigNode `yaml:"database"`
}

func main() {
	var (
		configPath string
		group      string
	)
	flag.StringVar(&configPath, "config", "manifest/config/config.yaml", "数据库配置文件路径")
	flag.StringVar(&group, "group", "default", "数据库配置分组")
	flag.Parse()

	root, err := projectRoot()
	if err != nil {
		fail(err)
	}
	configPath, err = resolveConfigPath(root, configPath)
	if err != nil {
		fail(err)
	}

	node, err := loadDatabaseConfig(configPath, group, root)
	if err != nil {
		fail(err)
	}
	if err := ensureSqliteDir(node.Link); err != nil {
		fail(err)
	}

	db, err := gdb.New(node)
	if err != nil {
		fail(gerror.Wrap(err, "创建数据库连接失败"))
	}
	defer db.Close(context.Background())

	ctx := context.Background()
	dbType, err := databaseType(&node)
	if err != nil {
		fail(err)
	}
	if err := createMigrationTable(ctx, db, dbType); err != nil {
		fail(err)
	}
	if err := resetMigrations(ctx, db, root, dbType); err != nil {
		fail(err)
	}
	fmt.Println("数据库迁移已回退到 0。")
	if err := storagesvc.Migrate(ctx, db); err != nil {
		fail(gerror.Wrap(err, "重新应用数据库迁移失败"))
	}
	fmt.Println("数据库迁移已重新应用到最新版本。")
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", gerror.Wrap(err, "读取当前工作目录失败")
	}
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
	return "", gerror.New("无法定位项目根目录")
}

func resolveConfigPath(root string, configPath string) (string, error) {
	if strings.TrimSpace(configPath) == "" {
		return "", gerror.New("数据库配置文件路径不能为空")
	}
	if filepath.IsAbs(configPath) {
		return configPath, nil
	}
	return filepath.Join(root, filepath.Clean(configPath)), nil
}

func loadDatabaseConfig(configPath string, group string, root string) (gdb.ConfigNode, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return gdb.ConfigNode{}, gerror.Wrapf(err, "读取数据库配置文件失败: %s", configPath)
	}
	var cfg configFile
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return gdb.ConfigNode{}, gerror.Wrapf(err, "解析数据库配置文件失败: %s", configPath)
	}
	node, ok := cfg.Database[group]
	if !ok {
		return gdb.ConfigNode{}, gerror.Newf("数据库配置分组不存在: %s", group)
	}
	node.Link = normalizeSQLiteLink(node.Link, root)
	return node, nil
}

func normalizeSQLiteLink(link string, root string) string {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(link)), sqliteType+"::") {
		return link
	}
	start := strings.Index(link, "file(")
	if start < 0 {
		return link
	}
	start += len("file(")
	end := strings.Index(link[start:], ")")
	if end < 0 {
		return link
	}
	rawPath := link[start : start+end]
	if rawPath == "" || filepath.IsAbs(rawPath) {
		return link
	}
	absPath := filepath.Join(root, filepath.Clean(rawPath))
	return strings.Replace(link, rawPath, filepath.ToSlash(absPath), 1)
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

func resetMigrations(ctx context.Context, db gdb.DB, root string, dbType string) error {
	appliedVersions, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}
	if len(appliedVersions) == 0 {
		return nil
	}
	sort.Sort(sort.Reverse(sort.StringSlice(appliedVersions)))
	for _, version := range appliedVersions {
		downFile := filepath.Join(root, migrationsRoot, dbType, version+".down.sql")
		if !gfile.IsFile(downFile) {
			return gerror.Newf("缺少回滚迁移文件: %s", downFile)
		}
		if err := executeMigrationFile(ctx, db, downFile); err != nil {
			return gerror.Wrapf(err, "执行回滚迁移失败: %s", version)
		}
		if _, err := db.Model("schema_migrations").Ctx(ctx).
			WhereIn("version", storagesvc.MigrationVersionAliases(version)).
			Delete(); err != nil {
			return gerror.Wrapf(err, "删除迁移记录失败: %s", version)
		}
		fmt.Printf("已回滚迁移: %s\n", version)
	}
	return nil
}

func appliedVersions(ctx context.Context, db gdb.DB) ([]string, error) {
	array, err := db.GetArray(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, gerror.Wrap(err, "读取已应用迁移失败")
	}
	versions := make([]string, 0, len(array))
	seen := make(map[string]struct{}, len(array))
	for _, item := range array {
		version := storagesvc.CanonicalMigrationVersion(strings.TrimSpace(item.String()))
		if version != "" {
			if _, ok := seen[version]; ok {
				continue
			}
			seen[version] = struct{}{}
			versions = append(versions, version)
		}
	}
	return versions, nil
}

func executeMigrationFile(ctx context.Context, db gdb.DB, path string) error {
	sqlText := strings.TrimSpace(gfile.GetContents(path))
	if sqlText == "" {
		return gerror.Newf("迁移文件为空: %s", path)
	}
	return db.Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		for _, statement := range splitStatements(sqlText) {
			if !hasExecutableSQL(statement) {
				continue
			}
			if _, err := db.Exec(ctx, statement); err != nil {
				return gerror.Wrapf(err, "执行迁移 SQL 失败: %s", path)
			}
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
