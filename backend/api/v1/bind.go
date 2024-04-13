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

type connectPageRequest struct {
	ParentID int `json:"parentId"`
	Relation int `json:"relation"`
	Kinship  int `json:"kinship"`
}

func connectPageToRelative(c *fiber.Ctx) error {
	apiURL := viper.GetString("external.url")
	apiAccessToken := c.Locals("api_access_token").(string)

	if apiAccessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var requestBody connectPageRequest
	if err := c.BodyParser(&requestBody); err != nil {
		zap.S().Debugln(err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if requestBody.ParentID == 0 || requestBody.Relation == 0 || requestBody.Kinship == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "parentId, relation, and kinship are required",
		})
	}

	reqBody, _ := json.Marshal(requestBody)

	zap.S().Debugln(string(reqBody))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/page/relative", apiURL), bytes.NewReader(reqBody))
	if err != nil {
		zap.S().Debugln(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiAccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Debugln(err.Error())
		return err
	}
	defer resp.Body.Close()

	zap.S().Debugln(req)
	zap.S().Debugln(fmt.Sprintf("Response Status: %d", resp.StatusCode))
	zap.S().Debugln(fmt.Sprintf("Response Status: %v", resp))

	if resp.StatusCode == http.StatusOK {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	} else {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": resp.Status,
		})
	}
}
