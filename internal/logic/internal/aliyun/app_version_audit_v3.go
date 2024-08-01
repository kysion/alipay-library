package aliyun

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/SupenBysz/gf-admin-community/utility/idgen"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PEM_BEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEM_END   = "\n-----END RSA PRIVATE KEY-----"
)

func RsaSign(signContent string, privateKey string, hash crypto.Hash) string {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature)
}

func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	privateKey = FormatPrivateKey(privateKey)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

func FormatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PEM_BEGIN) {
		privateKey = PEM_BEGIN + privateKey
	}
	if !strings.HasSuffix(privateKey, PEM_END) {
		privateKey = privateKey + PEM_END
	}
	return privateKey
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func main() {
	//fileDir, _ := os.Getwd()
	//fileName := "WechatIMG10010.png"
	//filePath := path.Join(fileDir, fileName)

	// 1.准备数据
	gMap := gmap.New()
	gMap.Set("app_version", "0.0.1")
	gMap.Set("version_desc", "小程序首次提交审核") // 版本描述
	gMap.Set("region_type", "CHINA")

	mapJson, _ := gjson.Encode(gMap)
	// {"app_version":"0.0.1","region_type":"CHINA","version_desc":"小程序首次提交审核"}

	//var biz_data = `{"app_version": "0.1"}`
	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)

	// 2.上传文件
	filePath := "/data/kysion-files/WechatIMG10013.png"
	filePath2 := "/data/kysion-files/WechatIMG10010.png"

	file, _ := os.Open(filePath)
	defer file.Close()

	file2, _ := os.Open(filePath2)
	defer file2.Close()

	// 文件1数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("first_screen_shot"), escapeQuotes(file.Name())))

	h.Set("Content-Type", "image/png")
	part, _ := writer.CreatePart(h)
	io.Copy(part, file)

	// 文件2数据
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("second_screen_shot"), escapeQuotes(file2.Name())))

	h.Set("Content-Type", "image/png")
	part, _ = writer.CreatePart(h)
	io.Copy(part, file2)

	//  请求数据
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))

	h.Set("Content-Type", "application/json")
	part, _ = writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()

	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/audit/apply"

	AliClient, _ := NewClient(context.Background(), "2021003179623086")
	var nonce = gconv.String(idgen.NextId())

	var auth_string = "app_id=" + "2021003179681073" + ",timestamp=" + strconv.FormatInt(time.Now().Unix(), 10) + ",nonce=" + nonce + ",app_cert_sn=" + AliClient.AppCertSN + ",alipay_root_cert_sn=" + AliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"POST" + "\n" +
		path + "\n" +
		biz_data + "\n"

	fmt.Println(body.String())

	fmt.Println(content)
	var sign = RsaSign(content, AliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据
	r, _ := http.NewRequest("POST", endpoing+path, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.Header.Add("alipay-request-id", "0b426e8716715217622333545ebcaa") // 每次请求需要有唯一请求标识。当需排查问题时，可以像技术支持提供请求的唯一标识
	r.Header.Add("authorization", authorization)

	// 5.发起请求
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)

		fmt.Println((time.Now().Unix()))

	}

}
