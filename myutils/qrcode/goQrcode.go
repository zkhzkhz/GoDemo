package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/skip2/go-qrcode"
	"image/color"
	"log"
)

func main() {
	err := qrcode.WriteFile("https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx76f1f6651ca0ad28&redirect_uri=http%3A%2F%2Fhljwx186.chsmarttv.com%2Fhljwx-tv-front%2F%23%2Fplaylist&response_type=code&scope=snsapi_base&connect_redirect=1#wechat_redirect&tvinfo=192.168.0.16|48:98:CA:85:99:D1|CHiQ_55Q5T_13CA|GITV", qrcode.Highest, 1024, "./qr.png")
	if err != nil {
		fmt.Println(err)
	}

	//不想直接生成一个PNG文件存储，我们想对PNG图片做一些处理，比如缩放了，旋转了，或者网络传输了等，
	//基于此，我们可以使用Encode函数，生成一个PNG 图片的字节流，这样我们就可以进行各种处理了。
	bytes, err := qrcode.Encode("https://www/google.com", qrcode.Highest, 256)
	if err != nil {
		log.Panic("d", err)
	}
	logs.Info("dd", bytes)

	qr, err := qrcode.New("https://www.baidu.com", qrcode.Highest)
	if err != nil {
		logs.Error(err)
	}
	qr.BackgroundColor = color.RGBA{50, 205, 50, 255}
	qr.ForegroundColor = color.Black
	_ = qr.WriteFile(256, "./dwff.png")
}
