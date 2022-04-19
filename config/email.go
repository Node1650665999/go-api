// Package config 站点配置信息
package config

import "gin-api/pkg/config"

func init() {
	config.Add("mail", func() map[string]interface{} {
		return map[string]interface{}{
			"host":     config.Env("MAIL_HOST", "localhost"),
			"port":     config.Env("MAIL_PORT", 1025),
			"username": config.Env("MAIL_USERNAME", ""),
			"password": config.Env("MAIL_PASSWORD", ""),
			"from":     config.Env("MAIL_FROM", ""),
			"is_ssl":   config.Env("MAIL_SSL", ""),
			"error_notice" : false,
		}
	})
}
