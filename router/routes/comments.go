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
		return
	}
	if err = r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	captchaID := r.PostForm.Get("captcha-id")
	userResponse := r.PostForm.Get("captcha-response")
	commentText := r.PostForm.Get("comment-text")
	entry, ok := utils.GLOBAL_CAPTCHA_STORE.GetCaptcha(captchaID)
	captchaError := ""
	if !ok || userResponse != entry.Solution {
		captchaError = "Error: Captcha is incorrect."
	}
	if userResponse == "" {
		captchaError = "Error: Captcha is empty"
	}
	if commentText == "" {
		captchaError = "Error: Body is empty."
	}
	w.Header().Set("HX-Trigger", "x-captcha-reload")
	if captchaError != "" {
		utils.GLOBAL_CAPTCHA_STORE.UpdateCaptchaError(captchaID, captchaError)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.GLOBAL_CAPTCHA_STORE.DeleteCaptcha(captchaID)

	author := r.PostForm.Get("comment-author")
	if author == "" {
		author = "Anonymous"
	}

	text := r.PostForm.Get("comment-text")
	escaped := render.EscapeHtml(text)
	html := render.MarkdownRender([]byte(escaped))
	html = bytes.TrimPrefix(html, []byte("<br/>"))
	html = append([]byte("<div class=\"py-2\"></div>"), html...)

	n, err := db.Conn.InsertComment(r.Context(), db.InsertCommentParams{
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Author:    author,
		Text:      text,
		Html:      string(html),
		PostID:    blogId,
	})
	if err != nil {
		fmt.Printf("err apicommentsget: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	components.Comment(n.NewBlogComment(), blogId).Render(r.Context(), w)
}

func apiCommentsGet(w http.ResponseWriter, r *http.Request) {
	b, err := strconv.Atoi(r.PathValue("id"))
	blogId := int64(b)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
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
	for _, v := range comments {
		uc = append(uc, v.NewBlogComment())
	}

	var buf bytes.Buffer

	for _, v := range uc {
		components.Comment(v, blogId).Render(r.Context(), &buf)
	}

	w.Write(buf.Bytes())
}

func ApiCaptcha(w http.ResponseWriter, r *http.Request) {
	captchaID := string(r.URL.Query().Get("captcha-id"))
	if captchaID != "" {
		entry, ok := utils.GLOBAL_CAPTCHA_STORE.GetCaptcha(captchaID)
		if ok && entry.Error != "" {
			captcha := utils.GenerateCaptcha()
			utils.GLOBAL_CAPTCHA_STORE.SetCaptcha(captcha.Id, captcha.Text)
			captcha.Error = entry.Error
			utils.GLOBAL_CAPTCHA_STORE.DeleteCaptcha(captchaID)
			components.CaptchaBox(*captcha).Render(r.Context(), w)
			return
		}
	}
	captcha := utils.GenerateCaptcha()
	utils.GLOBAL_CAPTCHA_STORE.SetCaptcha(captcha.Id, captcha.Text)
	components.CaptchaBox(*captcha).Render(r.Context(), w)
}
