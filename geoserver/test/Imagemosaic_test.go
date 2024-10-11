package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func Test_uploadImageMosaic(t *testing.T) {

	geoServerURL := "http://localhost:8080/geoserver/rest/workspaces/tutorial/coveragestores/wooo/file.imagemosaic"
	username := "admin"
	password := "geoserver"
	zipFilePath := "D:\\opt\\data\\SR_50M.zip"

	// Open the zip file
	file, err := os.Open(zipFilePath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return
	}
	defer file.Close()

	// Read the file content
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading zip file:", err)
		return
	}

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", geoServerURL, bytes.NewBuffer(fileContent))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set headers and authentication
	req.Header.Set("Content-Type", "application/zip")
	req.SetBasicAuth(username, password)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", string(body))
}

// geoServerURL: GeoServer 实例的基础 URL，这通常是 GeoServer 所在服务器的地址和端口。
// workspace: GeoServer 中定义的工作区（workspace）名称，用于逻辑上组织和管理地理数据。
// storeName: 图像马赛克存储的名称，这是在特定工作区下创建的一个存储，用于保存上传的栅格数据。
// filePath: 包含图像数据（通常是 TIFF/TIF 文件和相关配置文件）的 ZIP 文件在本地系统上的路径。
func TestCreateImageMosaic(t *testing.T) {
	geoServerURL := "http://localhost:8080/geoserver"
	workspace := "tutorial"
	storeName := "bev:996"
	// windows D:\\opt\\data\\demoData.zip
	// linux /d/opt/data/demoData.zip
	filePath := "D:\\opt\\data\\demoData.zip"
	storeURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/file.imagemosaic?configure=all&coverageName=%s", geoServerURL, workspace, storeName, storeName)
	fileData, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer fileData.Close()

	client := &http.Client{}
	req, err := http.NewRequest("PUT", storeURL, fileData)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/zip")
	req.SetBasicAuth("admin", "geoserver")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Adjust status code check to allow both 201 (Created) and 200 (OK)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(fmt.Errorf("failed to create/update mosaic, status: %s, body: %s", resp.Status, string(body)))
		return
	}
	return
}

func TestCreateImageStore(t *testing.T) {

	// GeoServer 配置信息
	geoServerURL := "http://localhost:8080/geoserver/rest/workspaces/tutorial/coveragestores"
	username := "admin"
	password := "geoserver"
	storeName := "bevff"
	fileURL := "file:data/mydata/Image"
	// 确保这个路径是服务器上实际存在的路径

	// 构建 XML 数据
	xmlData := fmt.Sprintf(`
<coverageStore>
    <name>%s</name>
    <type>ImageMosaic</type>
    <enabled>true</enabled>
    <workspace>tutorial</workspace>
    <url>%s</url>
</coverageStore>
`, storeName, fileURL)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", geoServerURL, bytes.NewBufferString(xmlData))
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return
	}

	// 设置请求头部
	req.Header.Set("Content-Type", "text/xml")
	req.SetBasicAuth(username, password)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", string(body))
}

func TestCreateImageMosaicByStore(t *testing.T) {

	geoServerURL := "http://localhost:8080/geoserver"
	workspace := "tutorial"
	storeName := "bev_ff"
	coverageName := "1234oofsd"
	username := "admin"
	password := "geoserver"
	// 构造请求 URL
	coverageURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s/coverages", geoServerURL, workspace, storeName)

	// 构建 XML 数据，配置新的栅格图层
	coverageXML := fmt.Sprintf(`<coverage>
		<name>%s</name>
		<title>%s mirey</title>
		<abstract>This is a sample image mosaic</abstract>
		<enabled>true</enabled>
		<store class="coverageStore">
			<name>%s</name>
		</store>
	</coverage>`, coverageName, coverageName, storeName)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", coverageURL, bytes.NewBufferString(coverageXML))
	if err != nil {
		return
	}

	// 设置 HTTP 头部和认证信息
	req.Header.Set("Content-Type", "text/xml")
	req.SetBasicAuth(username, password)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated {
		return
	}

	return
}

func TestDeleteImageMosaic(t *testing.T) {

	// GeoServer 配置
	geoServerURL := "http://localhost:8080/geoserver"
	workspace := "tutorial"
	storeName := "TTT"
	username := "admin"
	password := "geoserver"

	// 构造请求 URL
	storeURL := fmt.Sprintf("%s/rest/workspaces/%s/coveragestores/%s?recurse=true", geoServerURL, workspace, storeName)

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 创建 DELETE 请求
	req, err := http.NewRequest("DELETE", storeURL, nil)
	if err != nil {
		log.Println(fmt.Errorf("creating request failed: %v", err))
		return
	}

	//auth := username + ":" + password
	//str := base64.StdEncoding.EncodeToString([]byte(auth))
	//
	//fmt.Printf("BasicAuth：%s\n", str)
	// 设置 HTTP 基本认证
	req.SetBasicAuth(username, password)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println(fmt.Errorf("request failed: %v", err))
		return
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Println(fmt.Errorf("failed to delete mosaic, status: %s", resp.Status))
		return
	}

	return
}
