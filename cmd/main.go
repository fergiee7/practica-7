package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
//se van a usar las estructuras db.Create, db.Find, db.Delete, db.Save
// crear estructura
type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Definiendo variable para la base de datos para no crear otro documento
var db *gorm.DB

func main() {

	//Creando conexion con tabla
	dsn := "root:test@tcp(127.0.0.1:3306)/practica7?charset=utf8mb4&parseTime=True&loc=Local"
	var err error //para ir guardando los errores
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) //abriendo base de datos
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}

	//creando tabla si no existe usando la estructura existente User
	db.AutoMigrate(&User{})

	fmt.Println("Conexión exitosa y tabla creada o actualizada.")

	router := gin.Default()
	//users := []User{} //cambiar para usar base de
	//indexUser := 1

	//fmt.Println("Running App")

	//Tomar archivos de la carpeta template
	router.LoadHTMLGlob("templates/*")

	//revisar si está corriendo el servidor
	router.GET("/ping", func(c *gin.Context) { //se define una url
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//Entrando a la pagina
	router.GET("/", func(c *gin.Context) {

		users := []User{}
		db.Find(&users) // Obtener todos los usuarios

		c.HTML(200, "index.html", gin.H{
			"title":       "Main website",
			"total_users": len(users),
			"users":       users,
		})
	})

	//API URL obtener usuarios
	router.GET("/api/users", func(c *gin.Context) {
		users := []User{}
		db.Find(&users) // Obtener todos los usuarios
		c.JSON(200, users)

	})

	//CREAR USUARIO
	router.POST("/api/users", func(c *gin.Context) {
		var user User
		if c.BindJSON(&user) == nil {
			//user.Id = indexUser
			//users = append(users, user)
			//indexUser++

			// Crear un nuevo usuario en la base de datos
			result := db.Create(&user)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Error al crear usuario"})
				return
			}
			c.JSON(200, user)

		} else {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
			return
		}

	})

	//ELIMINAR USUARIO
	router.DELETE("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return

		}

		result := db.Delete(&User{}, idParsed)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar usuario"})
			return
		}
		c.JSON(200, gin.H{
			"message": "User deleted",
		})
		return

		c.JSON(201, gin.H{})

	})

	//ACTUALIZAR
	router.PUT("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return

		}

		var user User
		if c.BindJSON(&user) == nil {
			//mandando a buscar el usuario en la base de datos
			//user.Id = indexUser
			//users = append(users, user)
			//indexUser++
			var usuarioexistente User
			result := db.First(&usuarioexistente, idParsed)
			if result.Error != nil {
				c.JSON(404, gin.H{
					"error": "Usuario no encontrado",
				})
				return

			}

			//Se actualiza usuario
			usuarioexistente.Name = user.Name
			usuarioexistente.Email = user.Email
			db.Save(&usuarioexistente)

			c.JSON(200, gin.H{
				"error": "Usuario actualizado",
			})

		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})

		}

	})

	/*fmt.Println("Id a actualizar: ", id)
		for i, u := range users {
			if u.Id == idParsed {
				users[i] = user
				users[i].Id = idParsed
				c.JSON(200, users[i])
				return

			}
		}

		c.JSON(201, gin.H{})
	})*/

	router.Run(":8001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
