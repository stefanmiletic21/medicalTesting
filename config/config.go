package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig(fileName string, additionalDirs []string) (string, error) {
	viper.SetConfigName(fileName)

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	for _, dir := range additionalDirs {
		viper.AddConfigPath(dir)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return "", err
	}

	configFile := viper.ConfigFileUsed()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
	})

	return configFile, nil
}

func getConfigString(key string) string {
	return viper.GetString(key)
}

func getConfigBool(key string) bool {
	return viper.GetBool(key)
}

func getConfigInt(key string) int {
	return viper.GetInt(key)
}

func getConfigFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func getConfigStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func getConfigDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func getConfigTime(key string) time.Time {
	return viper.GetTime(key)
}
