package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/warlockz/ase-service/configs"
	"github.com/warlockz/ase-service/types"
	aebmes "github.com/warlockz/ase-service/utils/aebmes"
)

func EncryptFile(c *gin.Context) {
	// single file
	file, _ := c.FormFile("fileData")
	log.Println(file.Filename)
	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, "upload/"+newFileName)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save the file",
		})
		return
	}
	// // Create a file in the current directory to store the uploaded file
	// dst, err := os.Create(file.Filename)
	// defer dst.Close()
	// if err != nil {
	// 	c.String(http.StatusBadRequest, fmt.Sprintf("create file err: %s", err.Error()))
	// 	return
	// }

	// // Copy the uploaded file to the created file
	// if _, err = io.Copy(dst, file); err != nil {
	// 	c.String(http.StatusBadRequest, fmt.Sprintf("copy file err: %s", err.Error()))
	// 	return
	// }
	// fileDest := "upload/" + file.Filename
	var setDataParameter types.SetData

	c.ShouldBind((&setDataParameter))

	keywordArray := strings.Split(setDataParameter.Keywords, ",")

	aebmes.EncryptFile(keywordArray, newFileName)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func GetFile(c *gin.Context) {
	var setDataParameter types.SetData

	c.BindJSON((&setDataParameter))

	keywordArray := strings.Split(setDataParameter.Keywords, ",")
	configOperatorID := configs.EnvHederaOperatorID()
	fileName, err := aebmes.GenerateTrapdoor(keywordArray, configOperatorID)
	if err != nil {
		c.AbortWithError(400, err)
	}
	fileBytes, err := ioutil.ReadFile(fileName)

	c.Header("Content-Type", "application/octet-stream")
	//c.
	c.Data(http.StatusOK, "application/octet-stream", fileBytes)

}
