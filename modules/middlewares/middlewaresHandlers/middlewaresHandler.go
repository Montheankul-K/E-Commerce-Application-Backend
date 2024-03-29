package middlewaresHandlers

import (
	"strings"

	"github.com/Montheankul-K/E-Commerce-Application-Backend/config"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/entities"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/middlewares/middlewaresUsecases"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/packages/authentication"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewaresHandlersErrCode string

const (
	routerCheckErr middlewaresHandlersErrCode = "middleware-001"
	JwtAuthErr     middlewaresHandlersErrCode = "middleware-002"
	paramsCheckErr middlewaresHandlersErrCode = "middleware-003"
	authorizeErr   middlewaresHandlersErrCode = "middleware-004"
	apiKeyErr      middlewaresHandlersErrCode = "middleware-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler        // cors : ดักการเข้าถึง api
	RouterCheck() fiber.Handler // path check
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",                              // allow all ip to access
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH", // allow http method
		AllowHeaders:     "",
		AllowCredentials: false, // ใช้เป็น token แทน
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := authentication.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(JwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(JwtAuthErr),
				"no permission to access",
			).Res()
		}

		// set payload (fiber cache) เพื่อเอาไปใช้งานต่อ
		// set user id, role id
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next() // เรียกใช้ handler ตัวถัดไปเมื่อทำงานเสร็จ
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId") // get value from cache
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"parameter error",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRoleId ...int) fiber.Handler { // รับ params เข้ามาได้หลายตัว โดยจะเป็น array
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int) // convert to int
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}

		roles, err := h.middlewaresUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		// convert to binary
		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		// bitwise operator
		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key") // header X-Api-Key (field มาตรฐาน)
		if _, err := authentication.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyErr),
				"api key is invalid or required",
			).Res()
		}
		return c.Next()
	}
}
