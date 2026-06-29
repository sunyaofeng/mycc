package env

import (
	"os"
	"strings"
)

// Environment 环境类型
type Environment string

const (
	Development Environment = "development"
	Testing     Environment = "testing"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// GetEnv 获取当前环境
func GetEnv() Environment {
	env := os.Getenv("MYCC_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		return Development
	}

	switch strings.ToLower(env) {
	case "dev", "develop", "development":
		return Development
	case "test", "testing":
		return Testing
	case "stag", "staging":
		return Staging
	case "prod", "production":
		return Production
	default:
		return Development
	}
}

// IsDevelopment 是否是开发环境
func IsDevelopment() bool {
	return GetEnv() == Development
}

// IsTesting 是否是测试环境
func IsTesting() bool {
	return GetEnv() == Testing
}

// IsStaging 是否是预发布环境
func IsStaging() bool {
	return GetEnv() == Staging
}

// IsProduction 是否是生产环境
func IsProduction() bool {
	return GetEnv() == Production
}

// String 返回环境字符串
func (e Environment) String() string {
	return string(e)
}
