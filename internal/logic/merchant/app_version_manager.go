package merchant

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kysion/alipay-library/alipay_model"
	service "github.com/kysion/alipay-library/alipay_service"
	"github.com/kysion/alipay-library/internal/logic/internal/aliyun"
	"github.com/kysion/gopay/pkg/xpem"
	"github.com/kysion/gopay/pkg/xrsa"
	"github.com/yitter/idgenerator-go/idgen"
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

// 小程序开发管理
type sAppVersion struct{}

func NewAppVersion() *sAppVersion {
	return &sAppVersion{}
}

const (
	// PEM_BEGIN ... PEM_END PKCS1格式
	PEM_BEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEM_END   = "\n-----END RSA PRIVATE KEY-----"

	// PKCS8格式
	//PEM_BEGIN = "-----BEGIN PRIVATE KEY-----\n"
	//PEM_END   = "\n-----END PRIVATE KEY-----"
)

// RsaSign 签名
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

// ParsePrivateKey 解析密钥，返回rsa格式私钥
func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	// 格式化支付宝普通应用秘钥
	priKeyFormat := xrsa.FormatAlipayPrivateKey(privateKey)

	// 编码成rsa格式私钥
	priKey, err := xpem.DecodePrivateKey([]byte(priKeyFormat))
	if err != nil {
		return nil, err
	}

	return priKey, nil
}

// FormatPrivateKey 格式化私钥
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

func (s *sAppVersion) newClient(ctx context.Context) (client *aliyun.AliPay, err error) {
	// 4.创建请求对象，设置文件数据
	id := g.RequestFromCtx(ctx).Get("appId").String()

	appId, _ := strconv.ParseInt(id, 32, 0)

	app, err := service.MerchantAppConfig().GetMerchantAppConfigByAppId(ctx, gconv.String(appId))
	if err != nil {
		return
	}

	aliClient, _ := aliyun.NewClient(context.Background(), app.AppId)

	return aliClient, nil
}

type RequestData struct {
	Method          string            `json:"method" dc:"请求方法"`
	Url             string            `json:"url" dc:"请求接口地址，endpoing + path组成"`
	Body            *bytes.Buffer     `json:"body" dc:"请求体数据，包含文件数据+请求数据"`
	Writer          *multipart.Writer `json:"writer" dc:"输入流"`
	AlipayRequestId string            `json:"alipay_request_id" dc:"每次请求需要有唯一请求标识。当需排查问题时，可以像技术支持提供请求的唯一标识"`
	Authorization   string            `json:"authorization" dc:"认证鉴权信息，设置在请求头中"`
	AppAuthToken    string            `json:"app_auth_token" dc:"代调用需要传递商家应用的认证Token，需要同时添加该请求头"`
}

func (s *sAppVersion) request(ctx context.Context, info *RequestData) *bytes.Buffer {
	// 4.创建请求对象，设置数据
	r, _ := http.NewRequest(info.Method, info.Url, info.Body)
	r.Header.Add("Content-Type", info.Writer.FormDataContentType())
	r.Header.Add("alipay-request-id", info.AlipayRequestId) // 每次请求需要有唯一请求标识。当需排查问题时，可以像技术支持提供请求的唯一标识
	r.Header.Add("authorization", info.Authorization)
	if info.AppAuthToken != "" {
		r.Header.Add("alipay-app-auth-token", info.AppAuthToken)
	}

	// 5.发起请求
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
		return nil
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
		return body

	}
	return nil
}

// SubmitVersionAudit 提交应用版本审核
func (s *sAppVersion) SubmitVersionAudit(ctx context.Context, info *alipay_model.AppVersionAuditReq, pic *alipay_model.AppVersionAuditPicReq) (*alipay_model.AppVersionAuditRes, error) {
	//fileDir, _ := os.Getwd()
	//fileName := "WechatIMG10010.png"
	//filePath := path.Join(fileDir, fileName)
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	gMap.Set("app_version", info.AppVersion)
	gMap.Set("version_desc", info.VersionDesc) // 版本描述
	gMap.Set("region_type", "CHINA")
	gMap.Set("speed_up", "true")
	gMap.Set("auto_online", "true")

	if info.TestAccout != "" {
		gMap.Set("test_accout", info.TestAccout)
	}

	if info.TestPassword != "" {
		gMap.Set("test_password", info.TestPassword)
	}

	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)

	// 2.上传文件
	//filePath := "/data/kysion-files/tinyapp-consumer/img1.png"
	//filePath2 := "/data/kysion-files/tinyapp-consumer/img2.png"
	//testFileVideoPath := "/data/kysion-files/tinyapp-consumer/testVideo.zip"

	filePath := pic.FirstScreenShotPath
	filePath2 := pic.SecondScreenShotPath
	testFileVideoPath := info.TestFileName

	file, _ := os.Open(filePath)
	defer file.Close()

	file2, _ := os.Open(filePath2)
	defer file2.Close()

	file3, _ := os.Open(testFileVideoPath)
	defer file3.Close()

	// 文件1数据 (小程序应用截图)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("first_screen_shot"), escapeQuotes(file.Name())))

	h.Set("Content-Type", "image/png")
	part, _ := writer.CreatePart(h)
	io.Copy(part, file)

	// 文件2数据 (小程序应用截图)
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("second_screen_shot"), escapeQuotes(file2.Name())))

	h.Set("Content-Type", "image/png")
	part, _ = writer.CreatePart(h)
	io.Copy(part, file2)

	// 测试录屏文件数据 (小程序测试录屏)
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("test_file_name"), escapeQuotes(file3.Name())))

	h.Set("Content-Type", "zip")
	part, _ = writer.CreatePart(h)
	io.Copy(part, file3)

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
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"POST" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	//fmt.Println(body.String())

	//fmt.Println(content)
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "POST",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	response := s.request(ctx, &reqData)

	fmt.Println(response)

	var res alipay_model.AppVersionAuditRes
	json.Unmarshal(response.Bytes(), &res)
	return &res, nil

}

// CancelVersionAudit 撤销版本审核
func (s *sAppVersion) CancelVersionAudit(ctx context.Context, version string) (*alipay_model.CancelVersionAuditRes, error) {
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	gMap.Set("app_version", version)
	gMap.Set("bundle_id", "com.alipay.alipaywallet") // 小程序投放的端参数，例如投放到支付宝钱包是支付宝端。默认支付宝端。
	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)
	//  请求数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))

	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()
	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/audit/cancel"
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	// strconv.FormatInt(time.Now().Unix(), 10) 
	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"POST" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	fmt.Println(body.String())
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "POST",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	response := s.request(ctx, &reqData)

	var res alipay_model.CancelVersionAuditRes
	json.Unmarshal(response.Bytes(), &res)

	return &res, nil
}

// CancelVersion 退回开发版本
func (s *sAppVersion) CancelVersion(ctx context.Context, version string) (*alipay_model.CancelVersionRes, error) {
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	gMap.Set("app_version", version)
	gMap.Set("bundle_id", "com.alipay.alipaywallet") // 小程序投放的端参数，例如投放到支付宝钱包是支付宝端。默认支付宝端。

	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)
	//  请求数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))

	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()
	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/audited/cancel"
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	// strconv.FormatInt(time.Now().Unix(), 10) 
	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	// 待签名内容
	var content = auth_string + "\n" +
		"POST" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	fmt.Println(body.String())
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "POST",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	response := s.request(ctx, &reqData)

	var res alipay_model.CancelVersionRes
	json.Unmarshal(response.Bytes(), &res)

	return &res, nil
}

// AppOnline 小程序上架
func (s *sAppVersion) AppOnline(ctx context.Context, version string) (*alipay_model.AppOnlineRes, error) {
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	gMap.Set("app_version", version)                 // 商家版本
	gMap.Set("bundle_id", "com.alipay.alipaywallet") // 小程序投放的端参数，例如投放到支付宝钱包是支付宝端。默认支付宝端
	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)
	//  请求数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))

	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()
	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/online"
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	// strconv.FormatInt(time.Now().Unix(), 10) 
	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"POST" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	fmt.Println(body.String())
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "POST",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	response := s.request(ctx, &reqData)

	var res alipay_model.AppOnlineRes
	json.Unmarshal(response.Bytes(), &res)
	return &res, nil
}

// AppOffline 小程序下架
func (s *sAppVersion) AppOffline(ctx context.Context, version string) (*alipay_model.AppOfflineRes, error) {
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	gMap.Set("app_version", version)                 // 商家版本
	gMap.Set("bundle_id", "com.alipay.alipaywallet") // 小程序投放的端参数，例如投放到支付宝钱包是支付宝端。默认支付宝端
	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)
	//  请求数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))

	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()
	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/offline"
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	// strconv.FormatInt(time.Now().Unix(), 10) 
	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"POST" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	fmt.Println(body.String())
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "POST",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	response := s.request(ctx, &reqData)

	var res alipay_model.AppOfflineRes
	json.Unmarshal(response.Bytes(), &res)

	return &res, nil
}

// QueryAppVersionList 小程序版本列表查询 https://openapi.alipay.com/v3/alipay/open/mini/version/list/query GET
func (s *sAppVersion) QueryAppVersionList(ctx context.Context, versionStatus string) (res *alipay_model.QueryAppVersionListRes, err error) {
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//writer.WriteField("data", biz_data)

	//  请求数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))
	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()
	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/list/query?version_status=" + versionStatus + "&bundle_id=com.alipay.alipaywallet"
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	// strconv.FormatInt(time.Now().Unix(), 10) 
	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"GET" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	fmt.Println(body.String())
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "GET",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	response := s.request(ctx, &reqData)
	var queryRes alipay_model.QueryAppVersionListRes

	fmt.Println(response)
	json.Unmarshal(response.Bytes(), &queryRes)

	return &queryRes, err
}

// GetAppVersionDetail 小程序版本详情查询 https://openapi.alipay.com/v3/alipay/open/mini/version/detail/query  GET
func (s *sAppVersion) GetAppVersionDetail(ctx context.Context, version string) (*alipay_model.QueryAppVersionDetailRes, error) {
	aliClient, _ := s.newClient(ctx)

	// 1.准备数据
	gMap := gmap.New()
	mapJson, _ := gjson.Encode(gMap)

	var biz_data = string(mapJson)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("data", biz_data)
	//  请求数据
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="data"`))

	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	part.Write([]byte(biz_data))

	writer.Close()
	// 3.签名（文件内容不参与签名）
	var endpoing = "https://openapi.alipay.com" // v3版本
	var path = "/v3/alipay/open/mini/version/detail/query?app_version=" + version + "&bundle_id=com.alipay.alipaywallet"
	url := endpoing + path

	var nonce = gconv.String(idgen.NextId())

	// strconv.FormatInt(time.Now().Unix(), 10) 
	var auth_string = "app_id=" + aliClient.ThirdConfig.AppId + ",timestamp=" + gconv.String(time.Now().UnixNano()/1000000) + ",nonce=" + nonce + ",app_cert_sn=" + aliClient.AppCertSN + ",alipay_root_cert_sn=" + aliClient.AliPayRootCertSN

	var content = auth_string + "\n" +
		"GET" + "\n" +
		path + "\n" +
		biz_data + "\n" +
		aliClient.MerchantConfig.AppAuthToken + "\n"

	fmt.Println(body.String())
	var sign = RsaSign(content, aliClient.ThirdConfig.PrivateKey, crypto.SHA256)
	var authorization = "ALIPAY-SHA256withRSA " + auth_string + ",sign=" + sign

	// 4.创建请求对象，设置文件数据

	// 5.发起请求
	reqData := RequestData{
		Method:          "GET",
		Url:             url,
		Body:            body,
		Writer:          writer,
		AlipayRequestId: nonce,
		Authorization:   authorization,
		AppAuthToken:    aliClient.MerchantConfig.AppAuthToken,
	}

	var detailRes alipay_model.QueryAppVersionDetailRes
	response := s.request(ctx, &reqData)

	json.Unmarshal(response.Bytes(), &detailRes)

	return &detailRes, nil
}
