package handler

import (
	"fmt"
	"startup_back/internal/dto"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)
type AuthHandler struct{
	service *service.Services
}
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
  AccessToken string  `json:"-"`
}
func NewAuthHandler(service *service.Services) * AuthHandler{
	return &AuthHandler{service: service}
}
// SignUp godoc
// @Summary      Регистрация нового пользователя
// @Description  Создаёт нового пользователя и выдаёт access/refresh токены в cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserInput  true  "Информация о пользователе"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  map[string]string
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
 
  var input dto.CreateUserInput
  if err := c.BodyParser(&input); err != nil {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
   }
  
   resultUser,err:=  h.service.Auth.SignUpUser(c.Context(), input)
   if err != nil {
       return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
   }
    c.Cookie(&fiber.Cookie{
      Name:     "access_token",
      Value:    resultUser.AccessToken,
      HTTPOnly: true,
      Secure:   true,       
      SameSite: "Strict",  
      Path:     "/",
      MaxAge:   900,        
    })
    c.Cookie(&fiber.Cookie{
      Name:     "refresh_token",
      Value:    resultUser.RefreshToken,
      HTTPOnly: true,
      Secure:   true,
      SameSite: "Strict",
      Path:     "/",
      MaxAge:   60 * 60 * 24 * 30, 
    })
    response := UserResponse{
      ID:       resultUser.User.ID,
      Username: resultUser.User.Username,
      Email:    resultUser.User.Email,
    }
   return c.Status(fiber.StatusCreated).JSON(response)

}
// SignIn godoc
// @Summary      Авторизация пользователя
// @Description  Выполняет вход по email и паролю, устанавливает токены в cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.CreateUserInput  true  "Данные для входа"
// @Success      200  {object}  UserResponse
// @Failure      400  {object}  map[string]string
// @Router       /auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
  var input dto.CreateUserInput
  if err := c.BodyParser(&input); err != nil {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
   }
   resultUser,err := h.service.Auth.SignInUser(c.Context(), input.Email,input.Password)
   if err != nil {
       return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
   }
     c.Cookie(&fiber.Cookie{
      Name:     "access_token",
      Value:    resultUser.AccessToken,
      HTTPOnly: true,
      Secure:   true,       
      SameSite: "Strict",  
      Path:     "/",
      MaxAge:   900,        
    })
    c.Cookie(&fiber.Cookie{
      Name:     "refresh_token",
      Value:    resultUser.RefreshToken,
      HTTPOnly: true,
      Secure:   true,
      SameSite: "Strict",
      Path:     "/",
      MaxAge:   60 * 60 * 24 * 30, 
    })
   response := UserResponse{
    ID:       resultUser.User.ID,
    Username: resultUser.User.Username,
    Email:    resultUser.User.Email,
  }
   return c.Status(fiber.StatusOK).JSON(response)
}
// IdentityMe godoc
// @Summary      Проверка авторизованного пользователя
// @Description  Возвращает информацию о текущем пользователе по access_token из cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  UserResponse
// @Failure      401  {object}  map[string]string
// @Router       /auth/me [get]
func (h * AuthHandler) IdentityMe(c *fiber.Ctx) error {
  user, err := h.service.Auth.IdentityMe(c.Context(), c.Cookies("access_token"))
  fmt.Println(user)
  if err != nil {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
  }
  response := UserResponse{
    ID:       user.User.ID,
    Username: user.User.Username,
    Email:    user.User.Email,
  }
  return c.Status(fiber.StatusOK).JSON(response)
}
// Refresh godoc
// @Summary      Обновление токена доступа
// @Description  Обновляет access_token по refresh_token из cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/refresh [post]
func (h * AuthHandler) Refresh(c *fiber.Ctx) error {
  accessToken, err := h.service.Auth.RefreshToken(c.Context(), c.Cookies("refresh_token"))
  
  if err != nil {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
  }
  c.Cookie(&fiber.Cookie{
    Name:     "access_token",
    Value:    accessToken,
    HTTPOnly: true,
    Secure:   true,       
    SameSite: "Strict",  
    Path:     "/",
    MaxAge:   900,        
  })
  return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": accessToken})
}
// Logout godoc
// @Summary      Выход из аккаунта
// @Description  Удаляет access_token и refresh_token из cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/logout [post]
func (h * AuthHandler) LogOut(c *fiber.Ctx) error {
  c.ClearCookie("access_token")
  c.ClearCookie("refresh_token")
  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}