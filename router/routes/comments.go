package routes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

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
	if err = r.ParseMultipartForm(6*1024*1024); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

    // TODO: integrate NSFW check?
    imageID := db.Image{}
    image, imageInfo, err := r.FormFile("comment-image")
    if err != http.ErrMissingFile {
        imageBuf := make([]byte, imageInfo.Size)
        size, err := image.Read(imageBuf)
        if err != nil || int64(size) != imageInfo.Size {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        imageb64 := base64.StdEncoding.EncodeToString(imageBuf)
        ext := strings.TrimPrefix(filepath.Ext(imageInfo.Filename), ".")

        imageID, err = db.Conn.InsertIntoImage(r.Context(), db.InsertIntoImageParams{
            Image:  imageb64,
            Size:   imageInfo.Size,
            Ext:    ext,
        })
        if err != nil {
		    fmt.Printf("err apiimageinsert: %s", err)
		    w.WriteHeader(http.StatusInternalServerError)
		    return
        }
    }

	captchaID := r.FormValue("captcha-id")
	userResponse := r.FormValue("captcha-response")
	commentText := r.FormValue("comment-text")
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

	escaped := render.EscapeHtml(commentText)
	html := render.MarkdownRender([]byte(escaped))
	html = bytes.TrimPrefix(html, []byte("<br/>"))
	html = append([]byte("<div class=\"py-2\"></div>"), html...)

	comment, err := db.Conn.InsertComment(r.Context(), db.InsertCommentParams{
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Author:    author,
		Text:      commentText,
		Html:      string(html),
		PostID:    blogId,
        ImageID:   &imageID.ID,
	})
	if err != nil {
		fmt.Printf("err apicommentsget: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	components.Comment(comment.NewBlogComment(), blogId).Render(r.Context(), w)
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

func ApiCommentImage(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(strings.Split(r.PathValue("id"), ".")[0])
	imageID := int64(i)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
    image, err := db.Conn.SelectFromImage(context.Background(), imageID)
    if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
    }
    imageBuf := make([]byte, image.Size)
    n, err := base64.StdEncoding.Decode(imageBuf, []byte(image.Image))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    ext := "." + image.Ext
    mimeType := mime.TypeByExtension(ext)
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }

    w.Header().Set("Content-Type", mimeType)
    w.Header().Set("Content-Length", strconv.Itoa(n))
    w.WriteHeader(http.StatusOK)
    w.Write(imageBuf[:n])
}

func ApiCommentsDelete(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    adminPass := utils.Env().ADMIN_PASS
    passAttempt := r.FormValue("delete-password")
    err = bcrypt.CompareHashAndPassword([]byte(adminPass), []byte(passAttempt))
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    selectedComments := r.Form["selected-comments"]
    for _, idStr := range selectedComments {
        id, err := strconv.Atoi(idStr)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        err = db.Conn.DeleteCommentById(r.Context(), int64(id))
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    w.Header().Set("HX-Refresh", "true")
}
