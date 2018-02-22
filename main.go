package main

import (
  "fmt"
  "net/http"
  "os"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  _ "github.com/jinzhu/gorm/dialects/postgres"
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

type DisplayData struct {
  TileCollection []Tile `json:"tileCollection"`
  BaseColor string `json:"baseColor"`
  DisplayTitle string `json:"displayTitle"`
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

func GetDisplayData(c *gin.Context) {
  var displaydata DisplayData
  var tiles []Tile
  if err := db.Find(&tiles).Error; err != nil {
    c.AbortWithStatus(404)
    fmt.Println(err)
  } else {
    displaydata.TileCollection = tiles
    displaydata.BaseColor = "#009999"
    displaydata.DisplayTitle = ""
    c.JSON(200, displaydata)
  }
}

func main() {
  //db, err = gorm.Open("sqlite3", "./tiles.db");
  db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
  if (err != nil) { fmt.Println(err) }
  defer db.Close();

  db.AutoMigrate(&Tile{})

  router := gin.Default()

  router.LoadHTMLGlob("templates/*")

  router.GET("/api/tiles", GetTiles)
  router.POST("/api/tiles", CreateTile)
  router.DELETE("/api/tiles/:id", DestroyTile)
  router.GET("/api/tiles/:id/delete", DestroyTile)
  router.Static("/public", "./public")

  router.GET("/api/displaydata", GetDisplayData)

  router.GET("/index", func(c *gin.Context) {
    var tiles []Tile
    if err := db.Find(&tiles).Error; err != nil {
      c.AbortWithStatus(404)
      fmt.Println(err)
    } else {
      c.HTML(http.StatusOK, "index.tmpl", gin.H{ "tiles": tiles })
    }
  })

  router.Run(":" + os.Getenv("PORT"))
}
