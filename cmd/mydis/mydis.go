package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dolfly/mydis/pkg/storage/db"
	"github.com/spf13/viper"
	"github.com/tidwall/redcon"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "mydis"
	app.Usage = "my redis server"

	app.Flags = []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "conf"},
	}

	app.Before = func(c *cli.Context) error {
		conf := c.Path("config")
		viper.SetConfigName("mydis")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(conf)
		viper.AddConfigPath(".")
		viper.AddConfigPath("${HOME}")

		viper.SetEnvPrefix("MYDIS")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		viper.WatchConfig()
		return viper.ReadInConfig()
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println(viper.AllSettings())
		address := viper.GetString("address")
		fmt.Println(address)
		driver := viper.GetString("driver")
		fmt.Println(driver)
		sources := viper.GetStringSlice("sources")
		fmt.Println(sources)
		s, err := db.New(driver, sources...)
		if err != nil {
			log.Fatal(err)
			return err
		}
		err = redcon.ListenAndServe(address, s.Handler, s.Accept, s.Closed)
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Printf("started server at %s", address)
		return nil
	}
	app.Run(os.Args)
}
