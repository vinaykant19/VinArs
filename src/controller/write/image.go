package write

import (
	"encoding/json"
	"fmt"
	"../../service"
	response "../../../lib/core/response"
	"github.com/google/uuid"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (c *WriteController) UploadImage () {
	maxFileSize,_ := strconv.ParseInt(c.Configuration.FileUpload_MaxSize_MB, 10, 64)
	maxFileSize = maxFileSize * 1024

	if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
		fmt.Println("Could not parse multipart form: %v\n", err)
		c.Result, c.StatusCode = response.ErrorCodeResponse("FileUploadLimit", c.Language)
		return
	}

	itemType := c.Request.PostFormValue("type")
	file, fileHeader, err := c.Request.FormFile("uploadImg")
	if err != nil {
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidFile", c.Language)
		return
	}
	fmt.Println(itemType)
	defer file.Close()

	fileSize := fileHeader.Size
	fmt.Printf("File size (bytes): %v\n", fileSize)
	if fileSize > (maxFileSize * 1024){
		c.Result, c.StatusCode = response.ErrorCodeResponse("FileMaxSizeLimit", c.Language)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidFile", c.Language)
		return
	}
	fileType := http.DetectContentType(fileBytes)
	if fileType != "image/jpeg" && fileType != "image/jpg" &&
		fileType != "image/gif" && fileType != "image/png" {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidFileType", c.Language)
		return
	}
	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidFileType", c.Language)
		return
	}
	uid := uuid.New()
	fileName := uid.String()+fileEndings[0]
	newPath := filepath.Join(c.Configuration.FileUpload_Path, fileName)
	newFile, err := os.Create(newPath)
	if err != nil {
		fmt.Println( "CANT_WRITE_FILE")
		c.Result, c.StatusCode = response.ErrorCodeResponse("InternalError", c.Language)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		fmt.Println( "CANT_WRITE_FILE")
		c.Result, c.StatusCode = response.ErrorCodeResponse("InternalError", c.Language)
		return
	}

	iRes := map[string]string{
		"ImagePath": fileName,
	}
	var data []interface{}
	data = append(data, iRes)
	c.Result = response.ResponseWithData(data, "Image uploaded")
	return
}

func (c *WriteController) AddImageProduct (labId int64) {
	var l service.Image
	//check permission

	err := json.NewDecoder(c.Request.Body).Decode(&l)
	if err != nil {
		c.Result, c.StatusCode = response.ErrorCodeResponse("InvalidInput", c.Language)
		return
	}

	err, imageId := l.AddProductImage(c.DB)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}
	if imageId != 0 {
		lRes := map[string]string{
			"imageId": strconv.FormatInt(imageId, 10),
		}
		var data []interface{}
		data = append(data, lRes)
		c.Result = response.ResponseWithData(data, c.Language.GetMessageTextValue("ProductImageAdded"))
		return
	}
}

func (c *WriteController) DeleteProductImage (imageId  int64,  labId int64) {
	var img service.Image
	//check permission
	//@todo
	img.Id = imageId
	img.ConnectId = labId
	err := img.DeleteProductImage(c.DB)
	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}
	c.Result = response.ResponseWithoutData("SUCCESS", c.Language.GetMessageTextValue("ProductImageRemoved"))
	return
}
