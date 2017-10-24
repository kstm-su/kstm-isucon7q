package main

import (
	"os"
	"io/ioutil"
	"github.com/labstack/echo"
)

func syncImage(c echo.Context) error {
	avatarData, _ := ioutil.ReadAll(c.Request().Body)
	ioutil.WriteFile("/home/isucon/isubata/webapp/public/icons/"+c.Param("file_name"), avatarData, os.ModePerm)
	return nil
}
