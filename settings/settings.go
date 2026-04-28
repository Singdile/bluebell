package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Appconfig       Appconfig       `mapstructure:"app"`
	Logconfig       Logconfig       `mapstructure:"log"`
	Databaseconfig  Databaseconfig  `mapstructure:"mysql"`
	Redisconfig     Redisconfig     `mapstructure:"redis"`
	Snowflakeconfig Snowflakeconfig `mapstructure:"snowflake"`
}

type Appconfig struct {
	Appname string `mapstructure:"appname"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type Logconfig struct {
	Filename string `mapstructure:"filename"`
	Level    string `mapstructure:"level"`
}

type Databaseconfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Dbname   string `mapstructure:"dbname"`
}

type Redisconfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Poolsize int    `mapstructure:"poolsize"`
	Protocol int    `mapstructure:"protocol"`
	Db       int    `mapstructure:"db"`
}

type Snowflakeconfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

// 设置一个全局的配置变量
var Conf *Config

// 初始化配置文件
func Init() (err error) {
	// 配置查找
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // 搜索当前路径
	viper.AddConfigPath("./settings/")

	// 读取
	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
	}

	//必须要执行这一步，为Conf分配空间,否则原本是nil,没有分配对应的内存空间
	Conf = new(Config)

	//反序列化
	if err = viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("fatal error ummarshal config: %v", err)
	}

	//实时检查配置文件是否变动
	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("config file changed\n")

		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("热更新配置失败: %v\n", err)
		} else {
			fmt.Printf("配置同步到全局变量 Conf\n")
		}
	})

	return nil
}
