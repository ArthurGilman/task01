package api

import (
	"net/http"
	"strconv"
	"task01/internal/cache"
	"task01/internal/config"
	"task01/internal/db"
	"task01/internal/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var conf cors.Config

func Init() {
	conf = cors.DefaultConfig()
	conf.AllowOrigins = []string{"*"}
	conf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	conf.AllowHeaders = []string{"Content-Type"}

	port := config.Get().HTTP.Port

	addr := ":" + port

	router := gin.Default()

	router.Use(cors.New(conf))
	router.Use(logger.LoggerMiddleware())

	router.GET("api/persons", showPersons)
	router.POST("api/towers", addPerson)
	router.PUT("api/towers", editPerson)
	router.DELETE("api/towers", deletePerson)

	router.Run(addr)
}

func showPersons(c *gin.Context) {
	pers := []db.Person{}
	var p *db.Person
	var f db.Filter
	var err error

	filter := c.Query("filter")

	switch filter {
	case "id":
		id := c.Query("id")
		p, err = cache.GetPersonsFromCacheOrDatabase(id)

		if err != nil {
			logger.InfoLog.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get data from database"})
			return
		}
		c.JSON(http.StatusOK, p)
	case "age":
		f.Start = c.Query("start")
		f.End = c.Query("end")
		f.Orderby = c.Query("orderby")
		f.Limit = c.Query("limit")
		f.Offset = c.Query("offset")

		pers, err = db.GetFiltredAge(f.Start, f.End, f.Orderby, f.Limit, f.Offset)

		if err != nil {
			logger.InfoLog.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get data from database"})
			return
		}
		c.JSON(http.StatusOK, pers)
	}

}

func addPerson(c *gin.Context) {
	var p db.Person
	var id int
	var err error

	if err = c.ShouldBindJSON(&p); err != nil {
		logger.InfoLog.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if id, err = db.AddPerson(&p); err != nil {
		logger.InfoLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while adding information to the database"})
		return
	}

	p.Id = id

	cache.CachePerson(&p)

	c.JSON(http.StatusOK, id)
}

func editPerson(c *gin.Context) {
	p := db.Person{}

	if err := c.ShouldBindJSON(&p); err != nil {
		logger.InfoLog.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := db.UpdatePerson(&p); err != nil {
		logger.InfoLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while updating information in the database"})
		return
	}
	i := strconv.Itoa(p.Id)

	cache.DeleteFromCache(i)

	c.JSON(http.StatusOK, p)
}

func deletePerson(c *gin.Context) {
	type req struct {
		Id int `json:"id"`
	}

	var r req

	if err := c.ShouldBindJSON(&r); err != nil {
		logger.InfoLog.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := db.DeletePerson(r.Id); err != nil {
		logger.InfoLog.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while deleting information from the database"})
		return
	}
	i := strconv.Itoa(r.Id)

	cache.DeleteFromCache(i)

	c.JSON(http.StatusOK, r.Id)
}
