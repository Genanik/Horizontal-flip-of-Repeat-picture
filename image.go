package main

import (
	"C"
	"fmt"
	"github.com/nfnt/resize"
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
	//ResizeImgToBig()
	//ResizeImgToSmall()

	//g, err := readGif("dstImg")
	//if err != nil {
	//	fmt.Errorf("无法读取Gif图片，%v", err)
	//}
	//optimizationGif(g.Image)
	//writeGif("dst", g)
}

var noneColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}

//export ResizeImgToBig
func ResizeImgToBig() int {
	resizeImgToBIG("dstImg", "srcImg")
	return 0
}

func resizeImgToBIG(dst, src string) error {
	// 读取源图片
	img, imgType, err := readImg(src)
	if err != nil {
		return fmt.Errorf("无法读取img，%v", err)
	}

	newX := uint(1.5 * float64(img.Bounds().Max.X))

	midB := resize.Resize(newX, 0, img, resize.Bilinear)

	return writeImg(dst, imgType, midB)
}

//export ResizeImgToSmall
func ResizeImgToSmall() int {
	resizeImgToSMALL("dstImg", "srcImg")
	return 0
}

func resizeImgToSMALL(dst, src string) error {
	// 读取源图片
	img, imgType, err := readImg(src)
	if err != nil {
		return fmt.Errorf("无法读取img，%v", err)
	}

	newX := uint(0.5 * float64(img.Bounds().Max.X))

	midB := resize.Resize(newX, 0, img, resize.Bilinear)

	return writeImg(dst, imgType, midB)
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

//export HorizontalFilpPic
func HorizontalFilpPic() int {
	err := hConvertImg("dstImg", "srcImg")
	if err != nil {
		return -1
	}
	return 0
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
	dstPalette := make([]*image.Paletted, 0)

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

// 反优化GIF
func antiOptimizationGif(img []*image.Paletted) {
	for i := 1; i < len(img); i++ {
		img[i] = imgPlusImg(img[i-1], img[i])
	}
}
func imgPlusImg(img1, img2 *image.Paletted) *image.Paletted {
	// 将img2重叠至img1上方并返回
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
			if img2.At(x, y) == noneColor {
				continue // 透明 继续使用上一帧内容
			}
			dst.Set(x, y, img2.At(x, y))
		}
	}
	return dst
}

// 优化GIF
func optimizationGif(img []*image.Paletted) {
	for i := len(img) - 1; i > 0; i-- {
		println(len(img), " now:", i)
		img[i] = compareImage(img[i-1], img[i])
	}
}
func compareImage(img1, img2 *image.Paletted) *image.Paletted {
	// img2与img1对比 如果不一样，使用img2的内容，一样就使用noneColor
	//*只能使用未优化的GIF
	m2b := img2.Bounds()
	X := img1.Rect.Max.X
	Y := img1.Rect.Max.Y
	// 创建画板
	dst := image.NewPaletted(image.Rect(
		X,
		Y,
		0,
		0,
	), img1.Palette)
	// 遍历img2内容
	for x := m2b.Min.X; x < m2b.Max.X; x++ {
		for y := m2b.Min.Y; y < m2b.Max.Y; y++ {
			// 判断是否相同
			state := img2.At(x, y) == img1.At(x, y)
			if state {
				// 相同
				//dst.Set(x, y, noneColor)

				continue
			}
			dst.Set(x, y, img2.At(x, y))
		}
	}
	return dst
}
