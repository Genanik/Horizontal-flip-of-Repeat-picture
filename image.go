package main

import (
	"C"
	"fmt"
	"os"

	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
}

//export ConvertPic
func ConvertPic() int {
	err := convertImg("dstImg", "srcImg")
	if err != nil {
		return -1
	}
	return 0
}

func convertImg(dst, src string) error {
	// 读取源图片
	img, typ, err := readImg(src)
	if err != nil {
		return err
	}
	if typ == "gif" {
		g, err := readGif(src)
		if err != nil {
			return fmt.Errorf("无法读取Gif图片，%v", err)
		}
		vFlipGIF(g)
		return writeGif(dst, g)
	}
	// Debugf("图片处理", "成功解析一张%s图片", typ)
	// 翻转后写入新图片
	return writeImg(dst, typ, vFlip(img))
}

func readImg(src string) (image.Image, string, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, "", fmt.Errorf("无法打开图片，%v", err)
	}
	defer f.Close()
	img, typ, err := image.Decode(f)
	if err != nil {
		return nil, "", fmt.Errorf("无法解析图片，%v", err)
	}
	return img, typ, nil
}

func writeImg(dst, typ string, img image.Image) error {
	df, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建图片文件，%v", err)
	}
	defer df.Close()

	switch typ {
	case "png":
		err = png.Encode(df, img)
	case "jpeg":
		err = jpeg.Encode(df, img, nil)
	default:
		err = fmt.Errorf("未知格式: %v", typ)
	}

	if err != nil {
		fmt.Errorf("无法编码图片，%v", err)
	}
	return nil
}

func readGif(src string) (*gif.GIF, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("无法打开图片，%v", err)
	}
	defer f.Close()
	return gif.DecodeAll(f)
}

func writeGif(dst string, g *gif.GIF) error {
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建图片文件，%v", err)
	}
	defer f.Close()
	if err = gif.EncodeAll(f, g); err != nil {
		return fmt.Errorf("无法编码Gif图片，%v", err)
	}
	return nil
}

// 左右翻转
func vFlip(m image.Image) image.Image {
	mb := m.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, mb.Dx(), mb.Dy()))
	for x := mb.Min.X; x < mb.Max.X; x++ {
		for y := mb.Min.Y; y < mb.Max.Y; y++ {
			//  设置像素点
			dst.Set(mb.Max.X-x-1, y, m.At(x, y))
		}
	}
	return dst
}

func vFlipGIF(img *gif.GIF) {
	p := img.Image[0].Rect.Max.Sub(img.Image[0].Rect.Min)
	for i := 0; i < len(img.Image); i++ {
		m := img.Image[i]
		mb := m.Bounds()
		dst := image.NewPaletted(image.Rect(
			mb.Max.X,
			p.Y-mb.Max.Y,
			mb.Min.X,
			p.Y-mb.Min.Y,
		), m.Palette)
		for x := mb.Min.X; x < mb.Max.X; x++ {
			for y := mb.Min.Y; y < mb.Max.Y; y++ {
				// 设置像素点，此调换了Y坐标以达到垂直翻转的目的
				dst.Set(p.X-x-1, y, m.At(x, y))
			}
		}
		img.Image[i] = dst
	}
}
