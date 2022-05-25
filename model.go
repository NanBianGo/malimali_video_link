package main

// HtmlContent HTML内容对象
type HtmlContent struct {
	Url     string `json:"url"`
	UrlText string `json:"url_next"`
	Title   string
}

// SystemExit 系统退出标识
const SystemExit = "exit"
