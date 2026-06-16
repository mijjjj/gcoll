package i18n

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

const (
	// DefaultLanguage 是系统默认语言。
	DefaultLanguage = "zh-CN"
	// EnglishLanguage 是英文语言标识。
	EnglishLanguage = "en"

	langQueryKey  = "lang"
	langHeaderKey = "Accept-Language"
)

var supportedLanguages = map[string]string{
	"zh":    DefaultLanguage,
	"zh-cn": DefaultLanguage,
	"cn":    DefaultLanguage,
	"en":    EnglishLanguage,
	"en-us": EnglishLanguage,
	"en-gb": EnglishLanguage,
}

// Init 初始化 GoFrame i18n 默认配置。
func Init(ctx context.Context) {
	if path := resolveI18nPath(); path != "" {
		_ = g.I18n().SetPath(path)
	}
	g.I18n().SetLanguage(DefaultLanguage)
	_ = ctx
}

// resolveI18nPath 从当前目录向上查找 i18n 资源目录。
func resolveI18nPath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		candidate := filepath.Join(workingDir, "manifest", "i18n")
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			return candidate
		}
		parent := filepath.Dir(workingDir)
		if parent == workingDir {
			return ""
		}
		workingDir = parent
	}
}

// NormalizeLanguage 规范化外部传入语言。
func NormalizeLanguage(language string) string {
	normalized := strings.ToLower(strings.TrimSpace(language))
	if normalized == "" {
		return DefaultLanguage
	}
	if index := strings.Index(normalized, ","); index >= 0 {
		normalized = normalized[:index]
	}
	if index := strings.Index(normalized, ";"); index >= 0 {
		normalized = normalized[:index]
	}
	if matched, ok := supportedLanguages[normalized]; ok {
		return matched
	}
	return DefaultLanguage
}

// WithLanguage 把语言写入上下文。
func WithLanguage(ctx context.Context, language string) context.Context {
	return gi18n.WithLanguage(ctx, NormalizeLanguage(language))
}

// LanguageFromRequest 从请求参数和请求头中读取语言。
func LanguageFromRequest(queryLanguage string, headerLanguage string) string {
	if queryLanguage != "" {
		return NormalizeLanguage(queryLanguage)
	}
	return NormalizeLanguage(headerLanguage)
}

// T 翻译指定文案。
func T(ctx context.Context, key string) string {
	return g.I18n().T(ctx, key)
}

// QueryKey 返回语言查询参数名。
func QueryKey() string {
	return langQueryKey
}

// HeaderKey 返回语言请求头名。
func HeaderKey() string {
	return langHeaderKey
}
