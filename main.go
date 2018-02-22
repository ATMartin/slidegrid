package main

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
  db *gorm.DB
  err error
)

type Tile struct {
  ID uint `json:"id"`
  Type string `form:"type" json:"type"`
  Content string `form:"content" json:"content"`
}

// Tile CRUD

func GetTiles(c *gin.Context) {
  var tiles []Tile
  if err := db.Find(&tiles).Error; err != nil {
    c.AbortWithStatus(404)
    fmt.Println(err)
  } else {
    c.JSON(200, tiles)
  }
}

func CreateTile(c *gin.Context) {
  var tile Tile
  c.Bind(&tile)

  tile.Type = "Text"
  db.Create(&tile)
  c.Redirect(http.StatusFound, "/index")
}

func DestroyTile(c *gin.Context) {
  var tile Tile
  id := c.Params.ByName("id")
  db.Where("id = ?", id).Delete(&tile)
  c.Redirect(http.StatusFound, "/index")
}

func main() {
  db, err = gorm.Open("sqlite3", "./tiles.db");
  if (err != nil) { fmt.Println(err) }
  defer db.Close();

  db.AutoMigrate(&Tile{})

  router := gin.Default()

  router.LoadHTMLGlob("templates/*")

  router.GET("/api/tiles", GetTiles)
  router.POST("/api/tiles", CreateTile)
  router.DELETE("/api/tiles/:id", DestroyTile)
  router.GET("/api/tiles/:id/delete", DestroyTile)

  router.GET("/index", func(c *gin.Context) {
    var tiles []Tile
    if err := db.Find(&tiles).Error; err != nil {
      c.AbortWithStatus(404)
      fmt.Println(err)
    } else {
      c.HTML(http.StatusOK, "index.tmpl", gin.H{ "tiles": tiles })
    }
  })

  router.Run(":8080")
}
