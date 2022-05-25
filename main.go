package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

//解析嘛哩嘛哩网站的html,获取视频地址
func main() {
	fmt.Println("【嘛哩嘛哩动漫网址：https://www.malimali4.com】")
	for {
		go errorHandle()
		err := resolveVideoHtml()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

//解析视频HTML代码
func resolveVideoHtml() error {
	fmt.Print("请输入嘛哩嘛哩站点需要解析的详情页地址(exit退出)：")
	userScan := bufio.NewScanner(os.Stdin)
	if !userScan.Scan() {
		return errors.New("输入错误,程序退出")
	}
	inputText := strings.TrimSpace(userScan.Text())
	if inputText == "" {
		return errors.New("输入内容为空，请重新输入")
	}
	if inputText == SystemExit {
		fmt.Println("系统退出...")
		//系统退出
		os.Exit(200)
		return nil
	}
	htmlContent, err := getHtmlPage(inputText)
	if err != nil {
		return errors.New("错误：" + err.Error())
	}
	fmt.Println("》》》》》》》》》视频解析成功》》》》》》》》》")
	fmt.Printf("====== 视频标题：%s ======\n", htmlContent.Title)
	fmt.Printf("====== 视频地址：%s ======\n", htmlContent.Url)
	fmt.Println("》》》》》》》》》视频解析成功》》》》》》》》》")
	return nil
}

// 获取详情页面
func getHtmlPage(url string) (*HtmlContent, error) {
	htmlPageResp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("读取html网页错误:" + err.Error())
	}
	//结束前释放资源
	defer htmlPageResp.Body.Close()

	//解析html，加载body类型io.Reader
	documentBody, err := goquery.NewDocumentFromReader(htmlPageResp.Body)
	if err != nil {
		return nil, errors.New("读取html网页错误:" + err.Error())
	}

	//操作Selection对象
	playerHtml := documentBody.Find("#player").Children().Eq(1).Text()

	if playerHtml == "" {
		return nil, errors.New("未解析到视频地址，请更换或重试")
	}
	//拿到#player标签下的子标签的第二个
	jsonStr := strings.Split(playerHtml, "=")[1]

	var htmlContent HtmlContent
	//反序列化字符串到具体的类型
	err = json.Unmarshal([]byte(jsonStr), &htmlContent)

	if err != nil {
		return nil, errors.New("反序列化json字符串错误:" + err.Error())
	}
	//获取标题
	title := documentBody.Find(".u-title a span").Text()
	if title != "" {
		title = strings.ReplaceAll(title, " ", "-")
	}
	htmlContent.Title = title

	return &htmlContent, nil
}

func errorHandle() {
	err := recover()
	if err != nil {
		fmt.Println("处理异常：", err)
	}
}
