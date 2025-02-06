package configs

// @Title        configs.go
// @Description
// @Create       XdpCs 2025-02-01 下午8:12
// @Update       XdpCs 2025-02-01 下午8:12
var Conf = new(Config)

type Config struct {
	App   ApplicationConfig `mapstructure:"app"`
	Log   LoggerConfig      `mapstructure:"log"`
	DB    DBConfig          `mapstructure:"database"`
	Redis RedisConfig       `mapstructure:"redis"`
}

type ApplicationConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Env         string `mapstructure:"env"`
	LogfilePath string `mapstructure:"logfilePath"`
}
type LoggerConfig struct {
	Level    int8   `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Director string `mapstructure:"director"`
	ShowLine bool   `mapstructure:"show-line"`
}

type DBConfig struct {
	Driver      string `mapstructure:"driver"`
	AutoMigrate bool   `mapstructure:"migrate"`
	Dsn         string `mapstructure:"dsn"`
}
type RedisConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}
