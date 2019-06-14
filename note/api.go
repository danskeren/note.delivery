package note

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/danskeren/config"
	"github.com/danskeren/cookie"
	"github.com/danskeren/generate"
	"github.com/danskeren/note.delivery/db"
	"github.com/danskeren/note.delivery/templates"
	"github.com/dgraph-io/badger"
	"github.com/go-chi/chi"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"golang.org/x/crypto/bcrypt"
)

const (
	MAXNOTESIZE = 5242880 // 5 MB
)

type Note struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CanDelete bool   `json:"canDelete"`
	Password  []byte `json:"password"`
}

func NoteGET(w http.ResponseWriter, r *http.Request) {
	sess, _ := cookie.GetSession(r, config.File().GetString("session.key"))
	commonData := templates.ReadCommonData(w, r)
	noteBytes, err := db.BadgerDB.Get([]byte(chi.URLParam(r, "noteid")))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			commonData.MetaTitle = "404"
			templates.Render(w, "not-found.html", map[string]interface{}{
				"Common": commonData,
			})
			return
		} else {
			sess.AddFlash(err.Error())
			sess.Save(r, w)
		}
	}
	var note Note
	if err = json.Unmarshal(noteBytes, &note); err != nil {
		sess.AddFlash(err.Error())
		sess.Save(r, w)
	}
	passwordProtected := false
	if len(note.Password) > 0 {
		passwordProtected = true
	}
	templates.Render(w, "note.html", map[string]interface{}{
		"Common":            commonData,
		"NoteID":            note.ID,
		"CanDelete":         note.CanDelete,
		"Locked":            passwordProtected,
		"PasswordProtected": passwordProtected,
		"Note":              template.HTML(note.Content),
	})
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	r.Body = http.MaxBytesReader(w, r.Body, MAXNOTESIZE)
	sess, _ := cookie.GetSession(r, config.File().GetString("session.key"))
	if r.PostFormValue("note") == "" {
		sess.AddFlash("Note content cannot be empty.")
		sess.Save(r, w)
		templates.Render(w, "index.html", map[string]interface{}{
			"Common": commonData,
		})
	}
	var note Note
	var err error
	if r.PostFormValue("password") != "" {
		if note.Password, err = bcrypt.GenerateFromPassword([]byte(r.PostFormValue("password")), 12); err != nil {
			sess.AddFlash(err.Error())
			sess.Save(r, w)
			templates.Render(w, "index.html", map[string]interface{}{
				"Common":      commonData,
				"NoteContent": r.PostFormValue("note"),
			})
		}
	}
	if r.PostFormValue("allowDeletion") == "true" {
		note.CanDelete = true
	}
	for {
		id := generate.RandomLowercaseNumbers(8)
		if _, err := db.BadgerDB.Get([]byte(id)); err != badger.ErrKeyNotFound {
			continue
		}
		note.ID = id
		break
	}

	unsafeContent := blackfriday.Run([]byte(
		r.PostFormValue("note")),
		blackfriday.WithExtensions(blackfriday.HardLineBreak|blackfriday.FencedCode),
	)
	note.Content = string(bluemonday.UGCPolicy().SanitizeBytes(unsafeContent))

	noteBytes, err := json.Marshal(note)
	if err != nil {
		sess.AddFlash(err.Error())
		sess.Save(r, w)
		templates.Render(w, "index.html", map[string]interface{}{
			"Common":      commonData,
			"NoteContent": r.PostFormValue("note"),
		})
	}

	if err = db.BadgerDB.Set([]byte(note.ID), noteBytes); err != nil {
		sess.AddFlash(err.Error())
		sess.Save(r, w)
		templates.Render(w, "index.html", map[string]interface{}{
			"Common":      commonData,
			"NoteContent": r.PostFormValue("note"),
		})
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", note.ID), http.StatusSeeOther)
}

func UnlockNotePOST(w http.ResponseWriter, r *http.Request) {
	sess, _ := cookie.GetSession(r, config.File().GetString("session.key"))
	commonData := templates.ReadCommonData(w, r)
	noteBytes, err := db.BadgerDB.Get([]byte(chi.URLParam(r, "noteid")))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			commonData.MetaTitle = "404"
			templates.Render(w, "not-found.html", map[string]interface{}{
				"Common": commonData,
			})
			return
		} else {
			sess.AddFlash(err.Error())
			sess.Save(r, w)
		}
	}
	var note Note
	if err = json.Unmarshal(noteBytes, &note); err != nil {
		sess.AddFlash(err.Error())
		sess.Save(r, w)
	}
	err = bcrypt.CompareHashAndPassword(note.Password, []byte(r.PostFormValue("password")))
	if err != nil {
		fmt.Println(err)
		sess.AddFlash("Given password does not match the note password")
		sess.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/%s", chi.URLParam(r, "noteid")), http.StatusSeeOther)
		return
	}
	templates.Render(w, "note.html", map[string]interface{}{
		"Common":            commonData,
		"NoteID":            note.ID,
		"CanDelete":         note.CanDelete,
		"Locked":            false,
		"PasswordProtected": true,
		"Note":              template.HTML(note.Content),
	})
}

func DeleteNotePOST(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	sess, _ := cookie.GetSession(r, config.File().GetString("session.key"))

	noteBytes, err := db.BadgerDB.Get([]byte(chi.URLParam(r, "noteid")))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			commonData.MetaTitle = "404"
			templates.Render(w, "not-found.html", map[string]interface{}{
				"Common": commonData,
			})
			return
		} else {
			sess.AddFlash(err.Error())
			sess.Save(r, w)
		}
	}
	var note Note
	if err = json.Unmarshal(noteBytes, &note); err != nil {
		sess.AddFlash(err.Error())
		sess.Save(r, w)
	}

	if len(note.Password) > 0 {
		err = bcrypt.CompareHashAndPassword(note.Password, []byte(r.PostFormValue("confirm")))
		if err != nil {
			fmt.Println(err)
			sess.AddFlash("Given password does not match the note password")
			sess.Save(r, w)
			http.Redirect(w, r, fmt.Sprintf("/%s", chi.URLParam(r, "noteid")), http.StatusSeeOther)
			return
		}
	} else {
		if r.PostFormValue("confirm") != chi.URLParam(r, "noteid") {
			sess.AddFlash(fmt.Sprintf("Written NoteID (%s) does not match NoteID from the URL (%s)", r.PostFormValue("confirm"), chi.URLParam(r, "noteid")))
			sess.Save(r, w)
			http.Redirect(w, r, fmt.Sprintf("/%s", chi.URLParam(r, "noteid")), http.StatusSeeOther)
			return
		}
	}

	if note.CanDelete {
		if err := db.BadgerDB.Delete([]byte(chi.URLParam(r, "noteid"))); err != nil {
			sess.AddFlash(err.Error())
			sess.Save(r, w)
		}
	} else {
		sess.AddFlash("It is not allowed to delete this note.")
		sess.Save(r, w)
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", chi.URLParam(r, "noteid")), http.StatusSeeOther)
}
