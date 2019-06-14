package tmpl

const NotFoundHTML = `
{{ define "head" }}
<meta name="robots" content="noindex">
{{ end }}

{{ define "content" }}
<div>
  <h1>404! The page you're looking for does not exist.</h1>
  <p>
    Please double-check that you typed the URL correctly. If you're certain it's the correct URL, then the note has most likely been deleted.
  </p>
  <p>
    Notes created while the 'Allow Deletion' checkbox is enabled can be deleted by anybody who knows the URL.
  </p>
  <p>
    Return to homepage to <a href="/">create a new note</a>.
  </p>
</div>
{{ end }}
`
