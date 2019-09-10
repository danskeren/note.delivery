package tmpl

const IndexHTML = `
{{ define "content" }}

<form class="content" action="/" method="post">
  <div class="note">
    <textarea name="note" rows="10" placeholder="Be Water, My Friend" required>{{ .NoteContent }}</textarea>
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