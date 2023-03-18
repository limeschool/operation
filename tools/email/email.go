package email

import (
	"bytes"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
)

type Interface interface {
	Send(email, subject, context string) error
	SendAll(emails []string, subject, context string) error
}

type emailSender struct {
	Templates map[string]string `json:"templates" mapstructure:"templates"`
	Length    int               `json:"length" mapstructure:"length"`
	User      string            `json:"user" mapstructure:"user"`
	Host      string            `json:"host" mapstructure:"host"`
	Password  string            `json:"password" mapstructure:"password"`
	Title     string            `json:"title" mapstructure:"title"`
	Qrcode    string            `json:"qrcode" mapstructure:"qrcode"`
	Phone     string            `json:"phone" mapstructure:"phone"`
	Email     string            `json:"email" mapstructure:"email"`
	Desc      string            `json:"desc" mapstructure:"desc"`
	Logo      string            `json:"logo" mapstructure:"logo"`
	Footer    string            `json:"footer" mapstructure:"footer"`
}

type email struct {
	length   int
	user     string
	host     string
	password string
	Template string
	Title    string
	Qrcode   string
	Phone    string
	Email    string
	Content  string
	Desc     string
	Logo     string
	Context  string
	Footer   string
}

var (
	Sender emailSender
)

const (
	DefaultTemplate = "default"
)

func InitEmail(v *viper.Viper) {
	if v.UnmarshalKey("email", &Sender) != nil {
		panic("邮箱发送器初始化失败")
	}
	// 读取模板文件到内存
	templates := map[string]string{}
	for name, fileName := range Sender.Templates {
		file, err := os.Open(fileName)
		if err != nil {
			panic("邮箱模板初始化失败:" + err.Error())
		}
		key, err := ioutil.ReadAll(file)
		if err != nil {
			panic("读取邮箱模板失败:" + err.Error())
		}
		templates[name] = string(key)
	}
	Sender.Templates = templates
}

func (s *emailSender) New(tp string) Interface {
	return &email{
		length:   s.Length,
		user:     s.User,
		host:     s.Host,
		password: s.Password,
		Template: s.Templates[tp],
		Title:    s.Title,
		Qrcode:   s.Qrcode,
		Phone:    s.Phone,
		Email:    s.Email,
		Desc:     s.Desc,
		Logo:     s.Logo,
		Footer:   s.Footer,
	}
}

func (e email) parseTemplate(context string) string {
	n := template.New("")
	t, err := n.Parse(e.Template)
	if err != nil {
		return context
	}

	e.Content = context
	html := bytes.NewBuffer([]byte(""))
	if err = t.Execute(html, e); err != nil {
		return context
	}
	return html.String()
}

func (e *email) Send(email string, subject, context string) error {
	context = e.parseTemplate(context)
	hp := strings.Split(e.host, ":")
	auth := smtp.PlainAuth("", e.user, e.password, hp[0])
	ct := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + email + "\r\nFrom: " + e.Title + "\r\nSubject: " + subject + "\r\n" + ct + "\r\n\r\n" + context)
	return smtp.SendMail(e.host, auth, e.user, []string{email}, msg)
}

func (e *email) SendAll(emails []string, subject, context string) error {
	context = e.parseTemplate(context)
	hp := strings.Split(e.host, ":")
	auth := smtp.PlainAuth("", e.user, e.password, hp[0])
	ct := "Content-Type: text/html; charset=UTF-8"
	to := strings.Join(emails, ";")
	msg := []byte("To: " + to + "\r\nFrom: " + e.Title + ">\r\nSubject: " + subject + "\r\n" + ct + "\r\n\r\n" + context)
	return smtp.SendMail(e.host, auth, e.user, emails, msg)
}
