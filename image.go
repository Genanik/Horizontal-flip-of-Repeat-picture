package main

import (
	"C"
	"fmt"
	"image/color"
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
	//ReverseGif()
}

/*
	return code:
	-1	unexpected error
	0	successful
	1	is not gif
*/

//export ReverseGif
func ReverseGif() int {
	err, isGifImg := isGif("srcImg")
	if err != nil {
		return -1
	}
	if isGifImg {
		reverseGif("dstImg", "srcImg")
		return 0
	}
	return -1
}

//export HorizontalFilpPic
func HorizontalFilpPic() int {
	err := hConvertImg("dstImg", "srcImg")
	if err != nil {
		return -1
	}
	return 0
}

func hConvertImg(dst, src string) error {
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
		hFlipGIF(g)
		return writeGif(dst, g)
	}
	// Debugf("图片处理", "成功解析一张%s图片", typ)
	// 翻转后写入新图片
	return writeImg(dst, typ, hFlip(img))
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

// 水平翻转
func hFlip(m image.Image) image.Image {
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

func hFlipGIF(img *gif.GIF) {
	p := img.Image[0].Rect.Max.Sub(img.Image[0].Rect.Min)
	for i := 0; i < len(img.Image); i++ {
		m := img.Image[i]
		mb := m.Bounds()
		dst := image.NewPaletted(image.Rect(
			p.X-mb.Max.X,
			mb.Max.Y,
			p.X-mb.Min.X,
			mb.Min.Y,
		), m.Palette)
		for x := mb.Min.X; x < mb.Max.X; x++ {
			for y := mb.Min.Y; y < mb.Max.Y; y++ {
				// 设置像素点，此调换了X坐标以达到水平翻转的目的
				dst.Set(p.X-x-1, y, m.At(x, y))
			}
		}
		img.Image[i] = dst
	}
}

func isGif(src string) (error, bool) {
	// 读取源图片
	_, typ, err := readImg(src)
	if err != nil {
		return err, false
	}
	if typ == "gif" {
		return nil, true
	}
	return nil, false
}

func reverseGif(dst, src string) error {
	// 读取源图片
	g, err := readGif(src)
	if err != nil {
		return fmt.Errorf("无法读取Gif图片，%v", err)
	}
	rGIF(g)
	return writeGif(dst, g)
}

func rGIF(img *gif.GIF) {
	var dstPalette = make([]*image.Paletted, 0)

	// 反优化
	antiOptimizationGif(img.Image)

	for i := len(img.Image) - 1; i > -1; i-- {
		m := img.Image[i]
		mb := m.Bounds()
		dst := image.NewPaletted(image.Rect(
			mb.Max.X,
			mb.Max.Y,
			mb.Min.X,
			mb.Min.Y,
		), m.Palette)
		for x := mb.Min.X; x < mb.Max.X; x++ {
			for y := mb.Min.Y; y < mb.Max.Y; y++ {
				dst.Set(x, y, m.At(x, y))
			}
		}

		dstPalette = append(dstPalette, dst)
	}
	img.Image = dstPalette
}

func antiOptimizationGif(img []*image.Paletted) {
	for i := 1; i < len(img); i++ {
		img[i] = imgPlusImg(img[i-1], img[i])
	}

}

// 将img2重叠至img1上方并返回
func imgPlusImg(img1, img2 *image.Paletted) *image.Paletted {
	m1b := img1.Bounds()
	m2b := img2.Bounds()
	X := img1.Rect.Max.X
	Y := img1.Rect.Max.Y
	// 复制img1到dst
	dst := image.NewPaletted(image.Rect(
		X,
		Y,
		0,
		0,
	), img1.Palette)
	for x := m1b.Min.X; x < m1b.Max.X; x++ {
		for y := m1b.Min.Y; y < m1b.Max.Y; y++ {
			dst.Set(x, y, img1.At(x, y))
		}
	}
	// 复制img2到dst
	dst.Palette = append(dst.Palette)
	for x := m2b.Min.X; x < m2b.Max.X; x++ {
		for y := m2b.Min.Y; y < m2b.Max.Y; y++ {
			none := color.RGBA{R: 0, G: 0, B: 0, A: 0}
			if img2.At(x, y) != none {
				dst.Set(x, y, img2.At(x, y))
			}
		}
	}

	return dst
}
