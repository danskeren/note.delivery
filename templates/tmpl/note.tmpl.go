package tmpl

const NoteHTML = `
{{ define "head" }}
<meta name="robots" content="noindex, nofollow">
{{ end }}

{{ define "content" }}

{{ if .Locked }}

<p>
  The note exists, but it has been password-protected. Please unlock the note to view its content.
</p>
<form class="unlock" action="/{{ .NoteID }}" method="post">
  <div class="password">
      <input name="password" type="password" placeholder="Type Password" required>
    </div>
    <div class="submit">
      <button type="submit">Unlock Note</button>
    </div>
</form>

{{ else }}

{{ .Note }}

{{ if .CanDelete }}
<div class="delete__container">
  <hr>
  <p>
    This note can be deleted by anyone who knows the URL (https://note.delivery/{{ .NoteID }}).
  </p>
  <form class="delete" action="/{{ .NoteID }}/delete" method="post">
    <div class="confirm">
      {{ if .PasswordProtected }}
      <input name="confirm" type="password" placeholder="Type Note Password" required>
      {{ else }}
      <input name="confirm" type="text" placeholder="Type Note ID ( {{ .NoteID }} )" required>
      {{ end }}
    </div>
    <div class="submit">
      <button type="submit">Delete Note</button>
    </div>
  </form>
</div>
{{ end }}

{{ end }}

{{ end }}
`