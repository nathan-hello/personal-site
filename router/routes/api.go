package routes

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/nathan-hello/personal-site/components"
	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/utils"
)

func ApiComments(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
                apiCommentsGet(w, r)
                return
        }

        if r.Method == "POST" {
                apiCommentsPost(w, r)
                return
        }
}

func apiCommentsPost(w http.ResponseWriter, r *http.Request) {
        b, err := strconv.Atoi(r.PathValue("id"))
        blogId := int64(b)
        if err != nil {
                w.WriteHeader(http.StatusNotFound)
                return;
        }
        err = r.ParseForm()
        if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return;
        }


        author := r.PostForm.Get("comment-author")
        text := r.PostForm.Get("comment-text")
        escaped := render.EscapeHtml(text)
        html := render.MarkdownRender([]byte(escaped))
        html = bytes.TrimPrefix(html, []byte("<br/>"))
        html = append([]byte("<div class=\"py-2\"></div>"), html...)

        
        n, err := db.Conn.InsertComment(r.Context(), db.InsertCommentParams{
                CreatedAt: time.Now().UTC().Format(time.RFC3339),
                Author: author,
                Text: text,
                Html: string(html),
                PostID: blogId,
        })

        components.Comment(n.NewBlogComment(), blogId).Render(r.Context(), w)
}

func apiCommentsGet(w http.ResponseWriter, r *http.Request) {
        b, err := strconv.Atoi(r.PathValue("id"))
        blogId := int64(b)
        if err != nil {
                w.WriteHeader(http.StatusNotFound)
                return;
        }
        comments, err := db.Conn.SelectCommentsMany(r.Context(), blogId)
        if err != nil {
                if err == sql.ErrNoRows {
                        w.WriteHeader(http.StatusOK)
                        return
                }
                fmt.Printf("err apicommentsget: %s", err)
                w.WriteHeader(http.StatusInternalServerError)
                return
        }

        slices.SortFunc(comments, func(a, b db.Comment) int {
                return b.NewBlogComment().Date.Compare(a.NewBlogComment().Date)
        })
        
        uc := []utils.Comment{}  
        for _,v := range comments {
                uc = append(uc, v.NewBlogComment())
        }

        fmt.Sprintf("%#v\n", comments)
        var buf bytes.Buffer

        for _, v := range uc {
                components.Comment(v, blogId).Render(r.Context(), &buf)
        }


        w.Write(buf.Bytes())
}
