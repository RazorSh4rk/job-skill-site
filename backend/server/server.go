package server

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"razorsh4rk.github.io/jobsite/data"
	"razorsh4rk.github.io/jobsite/db"
)

func setup(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT"},
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/names", func(ctx *gin.Context) {
		allDBS := db.GetCollectionNames()
		log.Println(allDBS)
		ctx.JSON(http.StatusOK, gin.H{
			"names": allDBS,
		})
	})
	router.GET("/records/:collection", func(ctx *gin.Context) {
		collection := ctx.Params.ByName("collection")
		allDBS := db.GetCollectionNames()
		if !slices.Contains(allDBS, collection) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"records": []db.InterpretedListing{},
			})
			return
		}
		all := db.GetCollection(collection)
		ctx.JSON(http.StatusOK, gin.H{
			"records": all,
		})
	})
	router.GET("/heatmap/:collection", func(ctx *gin.Context) {
		collection := ctx.Params.ByName("collection")
		allDBS := db.GetCollectionNames()
		if !slices.Contains(allDBS, collection) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"records": []db.InterpretedListing{},
			})
			return
		}
		all := db.GetCollection(collection)
		hMap := data.GetSkillHeatmap(all)
		ctx.JSON(http.StatusOK, gin.H{
			"heatmap": hMap,
		})
	})
}

func Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	setup(router)
	router.Run()
}
