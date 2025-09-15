package handler

import (
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)
type UserHandler struct{
	service service.UserService
  authService service.AuthService
}
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
  AccessToken string `json:"access_token"`
}
func NewUserHandler(service service.UserService) * UserHandler{
	return &UserHandler{service:service}
}
func (h *UserHandler) Register(c *fiber.Ctx) error {

  var input service.CreateUserInput
  if err := c.BodyParser(&input); err != nil {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
   }
   resultUser,err:=  h.authService.SignUpUser(c.Context(), input)
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