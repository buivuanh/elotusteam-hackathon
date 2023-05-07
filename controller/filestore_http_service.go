package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/buivuanh/elotusteam-hackathon/domain"
	"github.com/buivuanh/elotusteam-hackathon/infrastructure"
	"github.com/buivuanh/elotusteam-hackathon/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type FileStoreHttpService struct {
	DB            *pgxpool.Pool
	DataStorePath string
	FileRepo      infrastructure.FileInfo
}

func (f *FileStoreHttpService) Upload(c echo.Context) error {
	// Read file
	fileData, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("could not read file: %v", err))
	}
	src, err := fileData.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("could not open file: %v", err))
	}
	defer src.Close()

	// Destination
	i := strings.LastIndex(fileData.Filename, ".")
	newName := fileData.Filename[0:i] + "_" + utils.GenerateRandomNumberString(5) + "." + fileData.Filename[i+1:len(fileData.Filename)]
	filePath := fmt.Sprintf("%s/%s", f.DataStorePath, newName)
	dst, err := os.Create(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("could not create file: %v", err))
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("could not copy file: %v", err))
	}

	// Get file info
	fileInfo, err := dst.Stat()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get file info: %v", err))
	}

	// Store file info to db
	fileObj := &domain.Image{
		FilePath:     filePath,
		OriginalName: fileData.Filename,
		ContentType:  c.Get("ContentType").(string),
		ByteSize:     fileInfo.Size(),
		OwnerID:      c.Get("user_id").(int),
	}
	_, err = f.FileRepo.InsertImage(c.Request().Context(), f.DB, fileObj)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("FileRepo.InsertImage: %v", err))
	}

	return c.NoContent(http.StatusOK)
}
