package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
)

func main() {
	http.HandleFunc("/", Image2Text)
	http.ListenAndServe(":8363", nil)
}

func Image2Text(w http.ResponseWriter, req *http.Request) {
	pic := req.URL.Query().Get("pic")

	if pic == "" {
		io.WriteString(w, "")
		return
	}
	template := `
<!DOCTYPE HTML>
<html>
<head>
  <meta http-equiv="content-type" content="text/html; charset=utf-8" />
  <style type="text/css" media="all">
    pre {
      white-space: pre-wrap;       /* css-3 */
      white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
      white-space: -pre-wrap;      /* Opera 4-6 */
      white-space: -o-pre-wrap;    /* Opera 7 */
      word-wrap: break-word;       /* Internet Explorer 5.5+ */
      font-family: 'Inconsolata', 'Consolas'!important;
      line-height: 0.99;
      font-size: %dpx;
    }
  </style>
</head>
<body>
  <pre>%s</pre>
</body>
</html>	
`
	htmlstr := ""
	maxLen := 100.0

	if pic == "" {
		io.WriteString(w, "")
		return
	}
	res, err := http.Get(pic)
	if err != nil {
		log.Printf("err is  %s  %v\n", pic, err)
	}
	defer res.Body.Close()
	im, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	width, height := im.Bounds().Max.X, im.Bounds().Max.Y
	rate := maxLen / math.Max(float64(width), float64(height))
	width, height = int(rate*float64(width)), int(rate*float64(height))
	im = resize.Resize(uint(width), uint(height), im, resize.Lanczos3)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			rgb := fmt.Sprintf("%v", im.At(w, h))
			reg := regexp.MustCompile(`(\d+)\s(\d+)\s(\d+)`)
			matches := reg.FindStringSubmatch(rgb)
			r, g, b := matches[1], matches[2], matches[3]
			htmlstr += fmt.Sprintf("<span style='color:rgb(%s,%s,%s);'>▇</span>", r, g, b)
		}
		htmlstr += "\n"
	}
	io.WriteString(w, fmt.Sprintf(template, 7, htmlstr))
}
