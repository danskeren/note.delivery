package tmpl

const RateLimitHTML = `
{{ define "head" }}
<meta name="robots" content="noindex">
{{ end }}

{{ define "content" }}
<div>
  <h1>Rate Limit!</h1>
  <p>
    To protect the server against abuse, you're only permitted to make 500 requests per hour. Your understanding is much appreciated.
  </p>
</div>
{{ end }}
`
