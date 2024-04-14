package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type accessTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Device   string `json:"device"`
}

// getApiAccessToken retrieves the access token from an external URL based on the request body.
//
// It takes a Fiber context object as a parameter.
// Returns an error.
func getApiAccessToken(c *fiber.Ctx) error {
	apiURL := viper.GetString("external.url")

	var requestBody accessTokenRequest
	var accessTokenResponse AccessTokenResponse

	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.Device == "" || requestBody.Email == "" || requestBody.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email, password and device are required",
		})
	}

	reqBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1/get-access-token", apiURL), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&accessTokenResponse)
	if err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Failed to get access token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"access_token": accessTokenResponse.AccessToken,
	})
}
