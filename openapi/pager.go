/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/jmespath/go-jmespath"
	"math"
	"strconv"
	"strings"
	"github.com/aliyun/aliyun-cli/cli"
	"github.com/aliyun/aliyun-cli/i18n"
)

var PagerFlag = &cli.Flag{Category: "caller",
	Name: "pager",
	Hidden: false,
	AssignedMode: cli.AssignedRepeatable,
	Aliases: []string{"--all-pages"},
	Short: i18n.T(
		"use `--pager` to merge pages for pageable APIs",
		"使用 `--pager` 在访问分页的API时合并结果分页"),
	Fields: []cli.Field{
		{Key: "", Required:false},
		{Key: "PageNumber", DefaultValue: "PageNumber", Short: i18n.T(" PageNumber", "指定PageNumber的属性")},
		{Key: "PageSize", DefaultValue: "PageSize", Short: i18n.T("PageSize", "")},
		{Key: "TotalCount", DefaultValue: "TotalCount", Short: i18n.T("TotalCount", "")},
	},
	ExcludeWith: []string {WaiterFlag.Name},
}

type Pager struct {
	PageNumberFlag string
	PageSizeFlag   string

	PageNumberExpr string
	PageSizeExpr   string
	TotalCountExpr string

	PageSize int

	totalCount        int
	currentPageNumber int
	collectionPath    string

	results []interface{}
}


func GetPager() *Pager {
	if !PagerFlag.IsAssigned() {
		return nil
	}
	pager := &Pager{}
	pager.PageNumberFlag, _ = PagerFlag.GetFieldValue("PageNumber")
	pager.PageSizeFlag, _ = PagerFlag.GetFieldValue("PageSize")
	pager.PageNumberExpr, _ = PagerFlag.GetFieldValue("PageNumber")
	pager.PageSizeExpr, _ = PagerFlag.GetFieldValue("PageSize")
	pager.TotalCountExpr, _ = PagerFlag.GetFieldValue("TotalCount")

	pager.collectionPath, _ = PagerFlag.GetFieldValue("")
	return pager
}

func (a *Pager) CallWith(invoker Invoker) (string, error) {
	for {
		resp, err := invoker.Call()
		if err != nil {
			return "", err
		}

		err = a.FeedResponse(resp.GetHttpContentString())
		if err != nil {
			return "", fmt.Errorf("call failed %s", err)
		}

		if !a.HasMore() {
			break
		}
		a.MoveNextPage(invoker.getRequest())
	}

	return a.GetResponseCollection(), nil
}

func (a *Pager) HasMore() bool {
	pages := int(math.Ceil(float64(a.totalCount) / float64(a.PageSize)))
	if a.currentPageNumber >= pages {
		return false
	} else {
		return true
	}
}

func (a *Pager) GetResponseCollection() string {
	root := make(map[string]interface{})
	current := make(map[string]interface{})
	path := a.collectionPath

	for {
		l := strings.Index(path, ".")
		if l > 0 {
			// fmt.Printf("%s %d\n", path, l)
			prefix := path[:l]
			root[prefix] = current
			path = path[l + 1:]
		} else {
			if strings.HasSuffix(path, "[]") {
				key := path[:len(path) - 2]
				current[key] = a.results
				break
			}
		}
	}

	s, err := json.Marshal(root)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func (a *Pager) FeedResponse(body string) error {
	var j interface{}
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		return fmt.Errorf("unmarshal %s", err.Error())
	}

	if total, err := jmespath.Search(a.TotalCountExpr, j); err == nil {
		a.totalCount = int(total.(float64))
	} else {
		return fmt.Errorf("jmespath: '%s' failed %s", a.TotalCountExpr, err)
	}

	if pageNumber, err := jmespath.Search(a.PageNumberExpr, j); err == nil {
		a.currentPageNumber = int(pageNumber.(float64))
	} else {
		return fmt.Errorf("jmespath: '%s' failed %s", a.PageNumberExpr, err)
	}

	if pageSize, err := jmespath.Search(a.PageSizeExpr, j); err == nil {
		a.PageSize = int(pageSize.(float64))
	} else {
		return fmt.Errorf("jmespath: '%s' failed %s", a.PageSizeExpr, err)
	}

	if a.collectionPath == "" {
		p2 := a.detectArrayPath(j)
		if p2 == "" {
			return fmt.Errorf("can't auto reconize collections path: you need add `--pager VSwitches.VSwitch[]` to assign manually")
		} else {
			a.collectionPath = p2
		}
	}

	a.mergeCollections(j)
	return nil
}

func (a *Pager) MoveNextPage(request *requests.CommonRequest) {
	a.currentPageNumber = a.currentPageNumber + 1
	// fmt.Printf("Move to page %d", a.currentPageNumber)
	request.QueryParams[a.PageNumberFlag] = strconv.Itoa(a.currentPageNumber)
}

func (a *Pager) mergeCollections(body interface{}) error {
	ar, err := jmespath.Search(a.collectionPath, body)
	if err != nil {
		return fmt.Errorf("jmespath search failed: %s", err.Error())
	} else if ar == nil {
		return fmt.Errorf("jmespath result empty: %s", a.collectionPath)
	}
	for _, i := range ar.([]interface{}) {
		a.results = append(a.results, i)
	}
	return nil
}

func (a *Pager) detectArrayPath(d interface{}) string {
	m, ok := d.(map[string]interface{})
	if !ok {
		return ""
	}
	for k, v := range m {
		// t.Logf("%v %v\n", k, v)
		if m2, ok := v.(map[string]interface{}); ok {
			for k2, v2 := range m2 {
				if _, ok := v2.([]interface{}); ok {
					return fmt.Sprintf("%s.%s[]", k, k2)
				}
			}
		}
	}
	return ""
}
