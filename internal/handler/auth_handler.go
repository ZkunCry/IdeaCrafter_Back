package handler

import (
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
// Signup user
// @Summary Register a new user
// @Description Создает нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.User true "User info"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Router /api/auth/signup [post]
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

func (h * AuthHandler) IdentityMe(c *fiber.Ctx) error {
  user, err := h.service.Auth.IdentityMe(c.Context())
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