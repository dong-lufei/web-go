package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"html/template"
	"regexp"
	"errors"
)

// 定义Page为一个结构体
type Page struct {
	Title string
	Body  []byte
}


// 模板缓存创建一个名为 的全局变量templates，并用 对其进行初始化ParseFiles。
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// 创建一个全局变量来存储我们的验证表达式：
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// 使用validPath 表达式来验证路径并提取页面标题
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	// 如果标题无效，该函数将向 HTTP 连接写入“404 Not Found”错误，并向处理程序返回错误
	if m == nil {
			http.NotFound(w, r)
			return "", errors.New("invalid Page Title")
	}
	// 如果标题有效，它将与nil 错误值一起返回
	return m[2], nil // 标题是第二个子表达式
}

// save方法来解决这个问题Page 接收者p是一个指向的指针Page。它不接受任何参数，并返回一个类型的值error
func (p *Page) save() error { 
	filename := p.Title + ".txt" 
	// 0600作为第三个参数传递给 的八进制整数字面WriteFile量表示该文件应创建为仅对当前用户具有读写权限
	return os.WriteFile(filename, p.Body, 0600) 
}

// 从 title 参数构造文件名，将文件的内容读入一个新变量body，并返回一个指向Page由正确的标题和正文值构造的文字的指针
func loadPage(title string) (*Page, error) { 
	filename := title + ".txt" 
	body, err := os.ReadFile(filename) 
	if err != nil { 
			return nil, err 
	} 
	return &Page{Title: title, Body: body}, nil
}

// 渲染模板
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p) 
	if err != nil { 
			http.Error(w, err.Error() , http.StatusInternalServerError) 
	} 

	// t, err := template.ParseFiles(tmpl + ".html") 
	// if err != nil { 
	// 		http.Error(w, err.Error(), http.StatusInternalServerError) 
	// 		return 
	// } 
	// err = t.Execute(w, p) 
	// if err != nil { 
	// 		http.Error(w, err.Error(), http.StatusInternalServerError) 
	// } 
}

var flag=false
var path=""

func handler(w http.ResponseWriter, r *http.Request) {
	if !flag {
		// r.URL.Path是请求 URL 的路径组件
		path=r.URL.Path[1:]
		fmt.Println( "http://localhost:8080/"+path )
	}
	flag=true
	fmt.Fprintf(w, "你好，我喜欢 %s!",path )
}

// 处理以“/view/”为前缀的 URL
func viewHandler(w http.ResponseWriter, r *http.Request, title string) { 
	// title, err := getTitle(w, r) 
	// if err != nil { 
	// 		return 
	// } 

	// title := r.URL.Path[len("/view/"):] 
	p, err := loadPage(title)
	if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
	}
	// fmt.Fprintf(w, " <h1>%s</h1><div>%s</div>", p.Title, p.Body) 

	// t, _ := template.ParseFiles("view.html")
	// t.Execute(w, p)

	renderTemplate(w, "view" , p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r) 
	// if err != nil { 
	// 		return 
	// } 

	// title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
			p = &Page{Title: title}
	}
	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// 		"<form action=\"/save/%s\" method=\"POST\">"+
	// 		"<textarea name=\"body\">%s</textarea><br>"+
	// 		"<input type=\"submit\" value=\"Save\">"+
	// 		"</form>",
	// 		p.Title, p.Title, p.Body)

	// 使用html/template库的template.ParseFiles函数读取内容 edit.html并返回一个*template.Template
	// t, _ := template.ParseFiles("edit.html")
	// // 执行模板，将生成的 HTML 写入http.ResponseWriter. .Title和带点的.Body标识符指的是 p.Title和p.Body
	// t.Execute(w, p)

	renderTemplate(w, "edit", p) 
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r) 
	// if err != nil { 
	// 		return 
	// } 
	
	// title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 闭包title从请求路径中提取 ，并使用正则validPath表达式对其进行验证
			m := validPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
				// 如果 title无效，将向 ResponseWriter使用http.NotFound函数写入错误
					http.NotFound(w, r)
					return
			}
			// 如果title有效， fn则将使用ResponseWriter、 Request和title作为参数调用封闭的处理程序函数
			fn(w, r, m[2])
	}
}

func main() {
	p1 := &Page{Title: "testPage", Body: []byte("这是一个示例页面。\n hello world!!")}
	p1.save()
	p2, _ := loadPage("testPage")
	fmt.Println(string(p2.Body))

	fmt.Println( "http://localhost:8080/view/testPage" )

	// http.HandleFunc函数处理对 Web 根 http的所有请求
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler)) 
	http.HandleFunc("/edit/", makeHandler(editHandler)) 
	http.HandleFunc("/save/", makeHandler(saveHandler))
	
	// http.ListenAndServe函数指定它应该在任何接口 ( ":8080") 上侦听端口 8080。（暂时不要担心它的第二个参数nil。）这个函数将一直阻塞，直到程序终止。
	// log.Fatal函数记录错误
	log.Fatal(http.ListenAndServe(":8080", nil))
}
