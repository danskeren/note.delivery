package tmpl

const IndexHTML = `
{{ define "content" }}

<form class="content" action="/" method="post">
  <div class="note">
    <textarea name="note" rows="10" placeholder="
# Markdown is supported

[Note.Delivery is open source](https://github.com/danskeren/note.delivery), licensed under AGPLv3.

Let's build a better internet, accessible without JavaScript.
" required>{{ .NoteContent }}</textarea>
  </div>
  <div class="password">
    <input name="password" type="password" placeholder="Password (optional)">
  </div>
  <div class="checkbox">
    <label>
      <span>Allow Deletion</span>
      <input name="allowDeletion" type="checkbox" value="true" checked>
      <span class="checkmark"></span>
    </label>
  </div>
  <div class="submit">
    <button type="submit">Create Note</button>
  </div>
</form>

{{ end }}
`