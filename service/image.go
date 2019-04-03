package service

import (
	"image"
	"image-process/models"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/nfnt/resize"
)

/*
ImageService that handle user upload end process image
*/
type ImageService interface {
	UploadFile(file multipart.File, imageInfo *models.Image) (map[string]string, error)
	ProcessImage(imageInfo *models.Image) (map[string]string, error)
}

/*
DefaultImageService is the default implementation for ImageService
*/
type DefaultImageService struct{}

/*
UploadFile handle the request from file
*/
func (*DefaultImageService) UploadFile(file multipart.File, imageInfo *models.Image) (map[string]string, error) {
	imageInfo.Name = filename(imageInfo)
	imageInfo.Loc = os.Getenv("FILE_LOC")
	if strings.HasSuffix(os.Getenv("FILE_LOC"), "/") {
		imageInfo.Loc += imageInfo.Name
	} else {
		imageInfo.Loc += "/" + imageInfo.Name
	}
	f, err := os.OpenFile(imageInfo.Loc, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	io.Copy(f, file)
	uid, err := models.GenerateNewImage(imageInfo)
	return map[string]string{"id": uid}, err
}

/*
ProcessImage process image into less sized and avatar sized pic
*/
func (*DefaultImageService) ProcessImage(imageInfo *models.Image) (map[string]string, error) {
	result := make(map[string]string, 2)
	headID, err := saveResizeImage("head", 200, imageInfo)
	if err != nil {
		return nil, err
	}
	avatarID, aErr := saveResizeImage("avatar", 120, imageInfo)
	if aErr != nil {
		return nil, aErr
	}
	result["head"] = headID
	result["avatar"] = avatarID
	return result, nil
}

/*
ProcessImageWithConfig process image according to its size
*/
func (*DefaultImageService) ProcessImageWithConfig(imageID string,
	configs []models.ImageProcessConfig) (map[string]string, error) {
	imageIDs := make(map[string]string)
	processImage := models.FindImageById(imageID)
	for _, config := range configs {
		newImageID, err := saveResizeImage(config.Name, config.Width, processImage)
		if err != nil {
			return nil, err
		}
		imageIDs[config.Name] = newImageID
	}
	return imageIDs, nil
}

func filename(image *models.Image) string {
	imageNameSplit := strings.Split(image.Name, ".")
	return image.Type + "-" + image.Relatedid + "." + imageNameSplit[len(imageNameSplit)-1]
}

func saveResizeImage(prefix string, size uint, imageInfo *models.Image) (string, error) {
	baseLoc := os.Getenv("FILE_LOC")
	if !strings.HasSuffix(baseLoc, "/") {
		baseLoc += "/"
	}
	headImageName := prefix + "-" + imageInfo.Name
	headImageLoc := baseLoc + headImageName
	resizeImage(imageInfo.Loc, headImageLoc, size)
	headImage := new(models.Image)
	headImage.Name = headImageName
	headImage.Relatedid = imageInfo.Relatedid
	headImage.Loc = headImageLoc
	headImage.Type = prefix
	models.GenerateNewImage(headImage)
	return headImage.Id, nil
}

func resizeImage(from, to string, size uint) error {
	// open "test.jpg"
	var err error

	file, err := os.Open(from)
	if err != nil {
		beego.Error("resizing", err)
		return err
	}

	var img image.Image
	// decode jpeg into image.Image
	if strings.HasSuffix(from, "jpg") || strings.HasSuffix(from, "jpeg") {
		img, err = jpeg.Decode(file)
		if err != nil {
			beego.Error("resizing2", err)
			return err
		}
	} else {
		img, err = png.Decode(file)
		if err != nil {
			beego.Error("resizing2", err)
			return err
		}
	}

	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(size, 0, img, resize.Lanczos3)

	out, err := os.Create(to)
	if err != nil {
		beego.Error("resizing3", err)
		return err
	}
	defer out.Close()

	// write new image to file
	png.Encode(out, m)
	return nil
}
