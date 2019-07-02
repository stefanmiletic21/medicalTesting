package config

import (
	"time"

	"github.com/spf13/viper"
)

var (
	GetGeneralIsTestingMode = getGeneralIsTestingMode

	GetHTTPReadTimeout         = getHTTPReadTimeout
	GetHTTPWriteTimeout        = getHTTPWriteTimeout
	GetSSLCertificatePath      = getSSLCertificatePath
	GetSSLKeystorePath         = getSSLKeystorePath
	GetHTTPServerAddress       = getHTTPServerAddress
	GetHTTPServerAddressSecure = getHTTPServerAddressSecure

	GetDatabaseConnectionString      = getDatabaseConnectionString
	GetDatabaseMaxIdleConnections    = getDatabaseMaxIdleConnections
	GetDatabaseMaxOpenConnections    = getDatabaseMaxOpenConnections
	GetDatabaseConnectionMaxLifetime = getDatabaseConnectionMaxLifetime

	GetQuestionStartString   = getQuestionStartString
	GetQuestionEndString     = getQuestionEndString
	GetQuestionTypeString    = getQuestionTypeString
	GetQuestionTextString    = getQuestionTextString
	GetQuestionAnswersString = getQuestionAnswersString

	GetQuestionTypeNamesFreeText      = getQuestionTypeNamesFreeText
	GetQuestionTypeNamesFreeNumerical = getQuestionTypeNamesFreeNumerical
	GetQuestionTypeNamesRadioGroup    = getQuestionTypeNamesRadioGroup
	GetQuestionTypeNamesCheckbox      = getQuestionTypeNamesCheckbox
)

func getGeneralIsTestingMode() bool {
	return getConfigBool("general.is_testing_mode")
}

func getHTTPReadTimeout() time.Duration {
	return getConfigDuration("http.http_read_timeout")
}

func getHTTPWriteTimeout() time.Duration {
	return getConfigDuration("http.http_write_timeout")
}

func getSSLCertificatePath() string {
	return getConfigString("http.ssl_certificate_path")
}

func getSSLKeystorePath() string {
	return getConfigString("http.ssl_keystore_path")
}

func getHTTPServerAddress() string {
	return getConfigString("http.http_server_address")
}

func getHTTPServerAddressSecure() string {
	return viper.GetString("http.http_server_address_secure")
}

func getDatabaseConnectionString() string {
	return getConfigString("db.connection")
}

func getDatabaseMaxIdleConnections() int {
	return getConfigInt("db.max_idle_connections")
}

func getDatabaseMaxOpenConnections() int {
	return getConfigInt("db.max_open_connections")
}

func getDatabaseConnectionMaxLifetime() time.Duration {
	return getConfigDuration("db.max_lifetime")
}

func getQuestionStartString() string {
	return getConfigString("questions.start")
}

func getQuestionEndString() string {
	return getConfigString("questions.end")
}

func getQuestionTypeString() string {
	return getConfigString("questions.type")
}

func getQuestionTextString() string {
	return getConfigString("questions.text")
}

func getQuestionAnswersString() string {
	return getConfigString("questions.answers")
}

func getQuestionTypeNamesFreeText() string {
	return getConfigString("question_type_names.free_text")
}

func getQuestionTypeNamesFreeNumerical() string {
	return getConfigString("question_type_names.free_numerical")
}

func getQuestionTypeNamesRadioGroup() string {
	return getConfigString("question_type_names.radio_group")
}

func getQuestionTypeNamesCheckbox() string {
	return getConfigString("question_type_names.checkbox")
}
