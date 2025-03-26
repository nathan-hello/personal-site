package utils

import (
	"bytes"
	cr "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"image"
	"image/jpeg"
	"math/rand"
	"sync"

	"image/color"

	"github.com/fogleman/gg"
)


type Captcha struct {
    Text        string
    Image       image.Image
    Compression int
    DarkMode    bool
    Id          string
    Error       string
}


var (
	BasicChars = []rune("23456789ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz")
    LightBg = color.RGBA{224, 238, 253, 255}
    DarkBg  = color.RGBA{18, 18, 18, 255}
)

var LightBasicColors = []color.RGBA{
	{214, 14, 50, 255},
	{240, 181, 41, 255},
	{176, 203, 40, 255},
	{105, 137, 194, 255},
	{242, 140, 71, 255},
}

var DarkBasicColors = []color.RGBA{
	{251, 188, 5, 255},
	{116, 192, 255, 255},
	{255, 224, 133, 255},
	{198, 215, 97, 255},
	{247, 185, 168, 255},
}

func getColor(darkMode bool) color.Color {
	if darkMode {
		return DarkBasicColors[rand.Intn(len(DarkBasicColors))]
	}
	return LightBasicColors[rand.Intn(len(LightBasicColors))]
}

func (c *Captcha) ToBase64() string {
    var buf bytes.Buffer
    jpeg.Encode(&buf, c.Image, &jpeg.Options{Quality: c.Compression})
    return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

type CaptchaBuilder struct {
	text        string
	width       int
	height      int
	darkMode    bool
	complexity  int
	compression int
}

func NewCaptchaBuilder() *CaptchaBuilder {
	return &CaptchaBuilder{
		width:       130,
		height:      40,
		darkMode:    false,
		complexity:  1,
		compression: 40,
	}
}

func (b *CaptchaBuilder) Text(text string) *CaptchaBuilder {
	b.text = text
	return b
}

func (b *CaptchaBuilder) Length(length int) *CaptchaBuilder {
	txt := make([]rune, length)
	for i := 0; i < length; i++ {
		txt[i] = BasicChars[rand.Intn(len(BasicChars))]
	}
	b.text = string(txt)
	return b
}

func (b *CaptchaBuilder) Width(w int) *CaptchaBuilder {
	b.width = w
	return b
}

func (b *CaptchaBuilder) Height(h int) *CaptchaBuilder {
	b.height = h
	return b
}

func (b *CaptchaBuilder) DarkMode(dm bool) *CaptchaBuilder {
	b.darkMode = dm
	return b
}

func (b *CaptchaBuilder) Complexity(c int) *CaptchaBuilder {
	if c < 1 {
		c = 1
	} else if c > 10 {
		c = 10
	}
	b.complexity = c
	return b
}

func (b *CaptchaBuilder) Compression(c int) *CaptchaBuilder {
	b.compression = c
	return b
}

func (b *CaptchaBuilder) build() *Captcha {
	if b.text == "" {
		b.Length(5)
	}
	dc := gg.NewContext(b.width, b.height)
	if b.darkMode {
		dc.SetColor(DarkBg)
	} else {
		dc.SetColor(LightBg)
	}
	dc.Clear()

	charCount := len(b.text)
	cSpacing := float64(b.width-10) / float64(charCount)
	y := float64(b.height)/2 - 15

	var fontSize float64
	if charCount <= 3 {
		fontSize = 42
	} else if charCount <= 5 {
		fontSize = 32
	} else {
		fontSize = 24
	}

	if err := dc.LoadFontFace("public/fonts/Terminus/TerminusTTF.ttf", fontSize); err != nil {
		// fmt.Println("Warning: could not load font, using default settings:", err)
	}

	for i, char := range b.text {
		dc.SetColor(getColor(b.darkMode))
		x := 5 + float64(i)*cSpacing
		angle := rand.Float64()*0.2 - 0.1 
        dc.Push()
		dc.RotateAbout(angle, x, y)
		dc.DrawStringAnchored(string(char), x, y, 0, 0.8)
		dc.Pop()
	}

	for i := 0; i < 2; i++ {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		x1, y1 := rand.Float64()*float64(b.width), rand.Float64()*float64(b.height)
		x2, y2 := rand.Float64()*float64(b.width), rand.Float64()*float64(b.height)
		dc.SetLineWidth(1)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	for i := 0; i < 2; i++ {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		x := rand.Float64() * float64(b.width)
		y := rand.Float64() * float64(b.height)
		r := 5 + rand.Float64()*5
		dc.DrawEllipse(x, y, r, r)
		dc.Stroke()
	}

	if b.complexity > 1 {
		for i := 0; i < 100*(b.complexity-1); i++ {
			dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
			dc.DrawPoint(rand.Float64()*float64(b.width), rand.Float64()*float64(b.height), 1)
			dc.Fill()
		}
	}

	return &Captcha{
		Text:        b.text,
		Image:       dc.Image(),
		Compression: b.compression,
		DarkMode:    b.darkMode,
	}
}

func GenerateCaptcha() *Captcha {
    captcha := NewCaptchaBuilder().
		Length(5).
		Width(180).
		Height(56).
		DarkMode(false).
		Complexity(3).
		Compression(40).
		build()
	id := make([]byte, 16)
	cr.Read(id)
    captcha.Id = hex.EncodeToString(id)
    return captcha
}

type storeEntry struct {
    Solution string
    Error    string
}
type captchaStore struct {
	sync.RWMutex
	store map[string]storeEntry
}

var GLOBAL_CAPTCHA_STORE captchaStore = captchaStore{
	store: make(map[string]storeEntry),
}

func (c *captchaStore) SetCaptcha(id, solution string) {
	c.Lock()
	defer c.Unlock()
	c.store[id] = storeEntry{Solution: solution, Error: ""}
}

func (c *captchaStore) GetCaptcha(id string) (storeEntry, bool) {
	c.RLock()
	defer c.RUnlock()
	entry, ok := c.store[id]
	return entry, ok
}

func (c *captchaStore) UpdateCaptchaError(id, errMsg string) {
	c.Lock()
	defer c.Unlock()
	if entry, ok := c.store[id]; ok {
		entry.Error = errMsg
		c.store[id] = entry
	}
}

func (c *captchaStore) DeleteCaptcha(id string) {
	c.Lock()
	defer c.Unlock()
	delete(c.store, id)
}
