package manager

import (
	"github.com/this-is-h/tikuAdapter-vercel/api/configs"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/registry"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/service"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/logger"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/ratelimit"
	"gorm.io/gorm"
)

var defaultManager Manager

// Manager 所有的组件都注册到这里
type Manager struct {
	db        *gorm.DB
	config    configs.Config
	ipLimiter *ratelimit.IPRateLimiter
	es        service.Elasticsearch
}

// RegistryManagerInterface manager interface
type RegistryManagerInterface interface {
	CloseManager() error
	GetDB() *gorm.DB
	GetIPLimiter() *ratelimit.IPRateLimiter
	GetConfig() configs.Config
	GetEs() *service.Elasticsearch
}

// GetManager get manager
func GetManager() Manager {
	return defaultManager
}

// CreateManager create manager
func CreateManager() Manager {
	config := registry.Config()
	logger.SetupGinLog()
	db := registry.RegisterDB(config)
	defaultManager = Manager{
		db:        db,
		config:    config,
		ipLimiter: registry.Limit(config),
		es:        registry.RegisterEs(config),
	}
	return defaultManager
}

// CloseManager close manager
func (m Manager) CloseManager() error {
	return registry.CloseDB()
}

// GetDB get db
func (m Manager) GetDB() *gorm.DB {
	return m.db
}

// GetIPLimiter get ip limiter
func (m Manager) GetIPLimiter() *ratelimit.IPRateLimiter {
	return m.ipLimiter
}

// GetConfig get config
func (m Manager) GetConfig() configs.Config {
	return m.config
}

// GetEs get es
func (m Manager) GetEs() service.Elasticsearch {
	return m.es
}
