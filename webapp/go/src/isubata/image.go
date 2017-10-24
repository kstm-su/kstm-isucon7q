package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"bytes"
	"compress/gzip"

	"github.com/labstack/echo"
)

func makeGzip(body []byte) ([]byte, error) {
    var b bytes.Buffer
    err := func() error {
        gw := gzip.NewWriter(&b)
        defer gw.Close()

        if _, err := gw.Write(body); err != nil {
            return err
        }
        return nil
    }()
    return b.Bytes(), err
}

func postProfile(c echo.Context) error {
	self, err := ensureLogin(c)
	if self == nil {
		return err
	}

	avatarName := ""
	var avatarData []byte

	if fh, err := c.FormFile("avatar_icon"); err == http.ErrMissingFile {
		// no file upload
	} else if err != nil {
		return err
	} else {
		dotPos := strings.LastIndexByte(fh.Filename, '.')
		if dotPos < 0 {
			return ErrBadReqeust
		}
		ext := fh.Filename[dotPos:]
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			break
		default:
			return ErrBadReqeust
		}

		file, err := fh.Open()
		if err != nil {
			return err
		}
		avatarData, _ = ioutil.ReadAll(file)
		file.Close()

		if len(avatarData) > avatarMaxBytes {
			return ErrBadReqeust
		}

		avatarName = fmt.Sprintf("%x%s", sha1.Sum(avatarData), ext)
	}

	if avatarName != "" && len(avatarData) > 0 {
		//_, err := db.Exec("INSERT INTO image (name, data) VALUES (?, ?)", avatarName, avatarData)
		avatarGzipData, err := makeGzip(avatarData)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile("/home/isucon/isubata/webapp/public/icons/"+avatarName+".gz", avatarGzipData, os.ModePerm)
		if err != nil {
			return err
		}
		req, err := http.NewRequest(
			"POST",
			"http://"+os.Getenv("ANOTHER_SERVER")+"/imagesync/"+avatarName+".gz",
			bytes.NewReader(avatarGzipData),
		)
		if err != nil {
			return err
		}
		client := &http.Client{}
		go client.Do(req)
		//_, err = client.Do(req)
		//if err != nil {
		//	return err
		//}


		_, err = db.Exec("UPDATE user SET avatar_icon = ? WHERE id = ?", avatarName, self.ID)
		if err != nil {
			return err
		}
	}

	if name := c.FormValue("display_name"); name != "" {
		_, err := db.Exec("UPDATE user SET display_name = ? WHERE id = ?", name, self.ID)
		if err != nil {
			return err
		}
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
