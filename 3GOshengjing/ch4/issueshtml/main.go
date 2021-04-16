// Issueshtml prints an HTML table of issues matching the search terms.

package main

import (
	"GoNote/3GOshengjing/ch4/github"
	"log"
	"os"
)

//!+template
import "html/template"

var issueList = template.Must(template.New("issuelist").Parse(`
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-allign: left'>
	<th>#</th>
	<th>State</th>
	<th>User</th>
	<th>Title</th>
</tr>
{{range .Items}}
<tr>
	<td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
	<td>{{.State}}</td>
	<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
	<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
	</tr>
	{{end}}
	</table>
	`))

//!- template

//!+
func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := issueList.Execute(os.Stdout, result); err != nil { //调用模板
		log.Fatal(err)
	}
}

//!-

// $ go build .../issueshtml
// $ ./issueshtml.exe repo:golang/go commenter:gopherbot json encoder >issues.html //会输出一个html文档

// $ ./issueshtml repo:golang/go 3133 10535 >issues2.html
//输出含有<和&的issue
//# State User Title
//10535 open dvyukov x/net/html: void element <link> has child nodes
//3133 closed ukai html/template: escape xmldesc as &lt;?xml
