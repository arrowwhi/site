package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"net/http"
	"site/auth"
	"site/database"
	"strings"
	"time"
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render выполняет рендеринг шаблона
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func MainHandler(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
	}

	// Данные для передачи в шаблон
	data := map[string]interface{}{
		"title":     "Главная",
		"LoginInfo": username,
	}

	err = renderBase(c, "home.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func ChatTemplateHandler(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
		if err != nil {
			log.Error(err.Error())
			return c.Redirect(http.StatusFound, "/logout")
		}
	} else {
		return c.Redirect(http.StatusFound, "/login")
	}
	logins := database.Get().GetLoginsToLine(username)
	for i, v := range logins {
		logins[i]["create_date"] = v["create_date"].(time.Time).Format("02 Jan 15:04")
	}
	if err != nil {
		log.Error(err.Error())
	}
	data := map[string]interface{}{
		"title":     "Чат",
		"LoginInfo": username,
		"chatId":    "wss://" + c.Request().Host + "/chat/web_socket",
		"logins":    logins,
	}

	err = renderBase(c, "chat.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func LoginHandler(c echo.Context) error {

	data := map[string]interface{}{
		"title": "Авторизация",
	}

	err := renderBase(c, "login.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func RegHandler(c echo.Context) error {

	data := map[string]interface{}{
		"title": "Регистрация",
	}

	err := renderBase(c, "reg.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func ChessHandler(c echo.Context) error {

	data := map[string]interface{}{
		"title": "Шахматы",
	}

	err := renderBase(c, "chess.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func ProfileHandler(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
	}
	fmt.Println("\n\n\n" + strings.TrimSpace(username) + "\n\n\n")
	user, err := database.Get().SelectClientByLogin(strings.TrimSpace(username))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if (user.ProfilePhoto == "") || (user.ProfilePhoto == "null") {
		user.ProfilePhoto = "/profiles/default.jpg"
	}

	data := map[string]interface{}{
		"title":     "Профиль",
		"LoginInfo": username,
		"FName":     user.FirstName,
		"LName":     user.LastName,
		"Photo":     user.ProfilePhoto,
	}

	err = renderBase(c, "profile.page.tmpl", data)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func LogoutHandler(c echo.Context) error {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Устанавливаем истекшее время жизни
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/")
}

func renderBase(c echo.Context, page string, data map[string]interface{}) error {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFiles("ui/html/header.partial.tmpl",
			"ui/html/footer.partial.tmpl",
			"ui/html/base.layout.tmpl",
			"ui/html/"+page)),
	}

	// Рендерим шаблон и отправляем результат клиенту
	err := renderer.Render(c.Response().Writer, page, data)
	if err != nil {
		return err
	}
	return err
}