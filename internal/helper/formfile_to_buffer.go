package helper

import (
	"bytes"
	"io"
	"mime/multipart"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func FormFileToBuffer(log *logrus.Logger, uploadedFile *multipart.FileHeader) (*bytes.Buffer, error) {
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png"}
	if !slices.Contains(allowedTypes, uploadedFile.Header.Get("Content-Type")) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Image file must be jpeg, jpg, or png")
	}

	fileContent, err := uploadedFile.Open()
	if err != nil {
		log.Warnf("Failed to open file content : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	defer fileContent.Close()

	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, fileContent); err != nil {
		log.Warnf("Failed to convert file content to buffer : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return buff, nil
}
