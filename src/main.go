package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"site/chat"
	"site/handlers"

	"site/auth"
	"site/database"
)

func main() {
	e := echo.New()

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Ошибка при загрузке файла .env:", err)
	//}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DBUSERNAME"),
		os.Getenv("DBPASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBSSLMODE"),
		os.Getenv("TIMEZONE"))
	database.Get().ConnectToDB(dsn)
	database.Get().CreateClientsTable()

	auth.SetJwtKey(os.Getenv("JWTKEY"))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Static("/", "ui/static")

	//e.GET("/test", test.WriteCookie)
	e.POST("/login", auth.TakeAuthHandler)
	e.POST("/reg", auth.TakeRegHandler)
	e.GET("/", handlers.MainHandler)
	e.GET("/login", handlers.LoginHandler)
	e.GET("/reg", handlers.RegHandler)
	e.GET("/chat", handlers.ChatTemplateHandler)
	e.GET("/chess", handlers.ChessHandler)
	e.GET("/profile", handlers.ProfileHandler)
	//e.GET("/markdown")

	//e.POST("/markdown", md.MarkdownHandler)
	e.GET("/logout", handlers.LogoutHandler)
	e.GET("/chat/web_socket", chat.ChatHandler)
	e.POST("/database/get_messages", handlers.GetMessagesHistory)
	e.POST("/feedback/add", handlers.TakeFeedback)
	e.POST("/database/get_logins", handlers.GetLogins)
	e.POST("/profile/change_photo", handlers.ChangeProfilePhoto)
	e.GET("/chat/take_logins", handlers.TakeUserLogins)

	e.POST("/test/post", func(c echo.Context) error {
		responseData := make(map[string]interface{})
		if err := c.Bind(&responseData); err != nil {
			return err
		}
		for k, v := range responseData {
			fmt.Println(k, ": ", v)
		}
		return c.String(http.StatusOK, "Profile Page")
	})

	e.Logger.Fatal(e.Start(":8080"))
	//e.Logger.Fatal(e.StartTLS(os.Getenv("STANDARTPORT"), os.Getenv("SERTFILE"), os.Getenv("KEYFILE")))
}