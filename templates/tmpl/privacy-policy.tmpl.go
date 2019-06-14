package tmpl

const PrivacyPolicyHTML = `
{{ define "content" }}
<h1>Privacy Policy</h1>
<p>
  I care deeply about user privacy, which is why I don't store any user-identifiable information. The only information that's stored on the server is the note data. The note data consist of the following:
</p>
<ul>
  <li>ID (randomly generated)</li>
  <li>Content (the note someone wrote)</li>
  <li>'Allow Deletion' boolean</li>
  <li>Hashed Password (bcrypt with 12 rounds)</li>
</ul>
<p>
  If 'Allow Deletion' has been checked when the note was created (it is checked by default), then the note data can be deleted by any person who has access to the URL. When a note is deleted, it is deleted completely. This means that a deleted note cannot be recovered.
</p>
<p>
  In order to protect Note.Delivery against abuse and attacks, I use a <a href="https://github.com/ulule/limiter" target="_blank">rate limiter package</a>, that will store your IP in memory for 1 hour. I also use CloudFlare, which will <a href="https://blog.cloudflare.com/what-cloudflare-logs/" target="_blank">discard access logs within 4 hours</a>. I do <strong>not</strong> log anything myself, as I feel confident that the rate limiter together with CloudFlare can mitigate any attacks that may occur. If you wish to self-host Note.Delivery, then you can easily remove the rate limiter package from the source code, as well as choose not to use CloudFlare if you so desire.
</p>
<p>
  Note.Delivery also use AES-256 encrypted cookies for the sole purpose of displaying potential error messages. You can easily delete and block the cookies without breaking anything (besides the fact that potential error messages won't be displayed if cookies are being blocked).
</p>
<p>
  If you encounter any problems, or have any questions, then feel free to <a href="https://github.com/danskeren/note.delivery/issues/new" target="_blank">open an issue on GitHub</a>.
</p>
{{ end }}
`
