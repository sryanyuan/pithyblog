package blog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ngaut/log"
)

func downloadHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	filename := vars["filename"]
	ctx.r.ParseForm()
	fileType := ctx.r.Form.Get("t")
	fileType = strings.ToLower(fileType)

	switch fileType {
	case "sitedata":
		{
			//	need super admin privilige
			if ctx.user.Permission < kPermission_SuperAdmin {
				ctx.RenderMessagePage("错误", "access denied", false)
				return
			}

			if len(filename) == 0 {
				ctx.RenderMessagePage("错误", "cannot find the file specific", false)
				return
			}

			//	open file
			downloadFilePath := filepath.Join("./sitedata-pack", filename)
			f, err := os.Open(downloadFilePath)
			if nil != err {
				ctx.RenderMessagePage("错误",
					fmt.Sprintf("cannot open the file (%s) specific: %v",
						downloadFilePath, err.Error()), false)
				return
			}
			defer f.Close()
			content, _ := ioutil.ReadAll(f)
			ctx.w.Header().Set("Content-Type", "application/zip")
			ctx.w.Write(content)
		}
	case "markdown_zip":
		{
			//	need super admin privilige
			if ctx.user.Permission < kPermission_SuperAdmin {
				ctx.RenderMessagePage("错误", "access denied", false)
				return
			}

			if len(filename) == 0 {
				ctx.RenderMessagePage("错误", "cannot find the file specific", false)
				return
			}

			//	open file
			f, err := os.Open("./markdown-articles/" + filename)
			if nil != err {
				ctx.RenderMessagePage("错误", "cannot open the file specific", false)
				return
			}
			defer f.Close()
			content, _ := ioutil.ReadAll(f)
			ctx.w.Header().Set("Content-Type", "application/zip")
			ctx.w.Write(content)
		}
	case "markdown":
		{
			// Get article id
			articleID, err := strconv.Atoi(ctx.r.FormValue("articleid"))
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			article, err := modelProjectArticleGet(articleID)
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			if !articleAccessible(ctx.user, article) {
				ctx.RenderMessagePage("错误", "access denied", false)
				return
			}
			fileBytes := []byte(article.ArticleContentMarkdown)
			ctx.w.Header().Set("Content-Type", "text/plain")
			ctx.w.Header().Set("Content-Disposition", "attachment;filename="+article.ArticleTitle+".md")
			//ctx.w.Header().Set("Content-Length", len(fileBytes))
			ctx.w.Write(fileBytes)
		}
	default:
		{
			log.Debugf("Invalid file type %v", fileType)
			ctx.RenderMessagePage("错误", "无效的文件索引符", false)
		}
	}
}
