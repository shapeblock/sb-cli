package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"io/ioutil"

	"github.com/spf13/viper"
)

var version = "master"

// GetVersion returns version of PusherCLI, set in ldflags.
func GetVersion() string {
	return version
}

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Can't get your home directory.")
		os.Exit(1)
	}

	return usr.HomeDir
}

func getConfigDir() string {
	return path.Join(getUserHomeDir(), ".config")
}

func GetConfigPath() string {
	return path.Join(getConfigDir(), "sb.json")
}

//Init sets the config files location and attempts to read it in.
func Init() {
	if _, err := os.Stat(getConfigDir()); os.IsNotExist(err) {
		err = os.Mkdir(getConfigDir(), os.ModeDir|0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(GetConfigPath()); os.IsNotExist(err) {
		err = ioutil.WriteFile(GetConfigPath(), []byte("{}"), 0600)
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(GetConfigPath())
	// viper.SetDefault("endpoint", "https://shapeblock.com")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
