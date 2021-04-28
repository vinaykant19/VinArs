package read
import (
	"../../service"
	"fmt"
	response "../../../lib/core/response"
)

type imgSuccessWithData struct {
	Images []service.Image
}

func (c *ReadController) ProductImages (productId int64) {
	var i service.Image

	result, err := i.GetAllImagesByProductId(c.DB, c.Language, productId)

	if err != nil {
		//log error
		fmt.Println(err)
		c.Result, c.StatusCode = response.ErrorResponse(err, c.Language)
		return
	}

	data := formatImageResponse(result)
	c.Result = response.ResponseWithData(data, "")
	return
}

func formatImageResponse(result []service.Image) []interface{} {
	pRes := imgSuccessWithData{
		Images: result,
	}
	var data []interface{}
	return append(data, pRes)
}
