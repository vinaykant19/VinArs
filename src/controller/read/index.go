package read

func (c *ReadController) Info () {
	c.Result = "{\"msg\":\"we are in index "  + " \"}"
	c.StatusCode = 200
}

func (c *ReadController) AboutUs () {
	c.Result = "{\"msg\":\"we are here\"}"
	c.StatusCode = 200
}

func (c *ReadController) Contact () {

}