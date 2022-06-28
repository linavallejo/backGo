package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Site struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type SiteDto struct {
	Url string `json:"url" binding:"required"`
}

type SiteResponseDto struct {
	ShortUrl string `json:"short_url"`
}

type TokenDto struct {
	Token string `uri:"token"`
}

var sites map[string]Site

func main() {
	sites = make(map[string]Site)
	r := gin.Default()
	r.POST("/shortener", shortener)
	r.GET("/:token", show)
	r.Run(":80") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func shortener(ctx *gin.Context) {
	var site SiteDto
	if err := ctx.ShouldBindJSON(&site); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rand.Seed(time.Now().UnixNano())
	token := randSeq(10)
	sites[token] = Site{token, site.Url}
	ctx.JSON(http.StatusOK, SiteResponseDto{ctx.Request.Host + "/" + token})
}

func show(ctx *gin.Context) {
	var tokenDto TokenDto
	if err := ctx.ShouldBindUri(&tokenDto); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	originalUrl := sites[tokenDto.Token].Url
	ctx.Redirect(http.StatusFound, originalUrl)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
