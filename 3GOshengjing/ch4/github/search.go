//!+
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) { //IssuesSearchResult是类型，//区别 int （类型）\*int 类型\ v:=&x \ &x （地址，属于*int型） \ *v （指针所指地址的的值）
	q := url.QueryEscape(strings.Join(terms, " "))
	//用url.QueryEscape来对查询中的特殊字符进行转义操作 //如空格，冒号等
	fmt.Println(q) //repo%3Agolang%2Fgo+is%3Aopen+json+decoder
	resp, err := http.Get(IssuesURL + "?q=" + q)
	// http.Gett函数是创建HTTP请求的函数，如果获取过程没有出错，那么会在resp这个结构体中得到访问的请求结果。
	//fmt.Println(resp) //&{...}

	if err != nil {
		return nil, err
	}

	//!-
	//For long-term stability, instead of http.Get, use the
	//variant below which adds an HTTP request header indicating
	//that only version 3 of the Github API is acceptable.

	/*
		req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set(
			"Accept", "application/vnd.github.va.text-match+json")
		resp, err := http.DefaultClient.Do(req)
	*/
	//!+

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents `defer`, which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		//resp的Body字段包括一个可读的服务器响应流。
		//resp.Body.Close()是关闭这个响应流？
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		//Decode(&result) 输入流解码JSON数据给该&result地址？
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil //返回的是个&result地址？
}
