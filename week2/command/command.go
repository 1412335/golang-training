package command

import (
	"fmt"
	"golang-training/week2/config"
	"golang-training/week2/provider"
	"sync"

	"gorm.io/gorm"

	"github.com/urfave/cli/v2"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// cli
var InfoCommand = &cli.Command{
	Name:  "info",
	Usage: "show info",
	Action: func(c *cli.Context) error {
		// fmt.Println("added task: ", c.Args().First())
		fmt.Println("show info")
		return nil
	},
}

var ServeCommand = &cli.Command{
	Name:  "server",
	Usage: "run the service",
	Action: func(c *cli.Context) error {
		// fmt.Println("added task: ", c.Args().First())
		cfg := config.NewConfig()
		_ = provider.MustBuildResourceProvider(cfg)
		return nil
	},
}

var MigrateCommand = &cli.Command{
	Name:  "migrate",
	Usage: "migrate db",
	Action: func(c *cli.Context) error {
		cfg := config.NewConfig()
		rp := provider.MustBuildResourceProvider(cfg)
		// only for development
		// Migrate the schema
		rp.GetDB().AutoMigrate(&Product{})
		// Create
		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < 10000; j++ {
					rp.GetDB().Create(&Product{
						Code:  fmt.Sprintf("%d:%d", i, j),
						Price: uint(100 + (i+1)*j),
					})
				}
			}(i)
		}
		wg.Wait()

		// Read
		// var product Product
		// db.First(&product, 1)                 // find product with integer primary key
		// db.First(&product, "code = ?", "D42") // find product with code D42

		// // Update - update product's price to 200
		// db.Model(&product).Update("Price", 200)
		// // Update - update multiple fields
		// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
		// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

		// // Delete - delete product
		// db.Delete(&product, 1)
		return nil
	},
}
