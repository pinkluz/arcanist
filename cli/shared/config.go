package shared

import (
	"github.com/spf13/viper"
)

// SetupConfig sets up the config that will be read from the filesystem
func SetupConfig(file string) {

	if file == "" {
		viper.SetConfigName(".garc")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/garc/")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(file)
	}

	// TODO: config isn't actually used at all so removing this for now
	//
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	fmt.Printf("error loading config: %s \n", err)
	// 	os.Exit(1)
	// }

}
