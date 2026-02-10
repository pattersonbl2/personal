---
title: "Contact"
url: "/contact/"
hideMeta: true
---

Get in touch. Send a message and I'll respond when I can.

<form action="https://api.ark31.info/api/contact" method="POST" class="contact-form">
  <p style="position:absolute;left:-9999px" aria-hidden="true">
    <label for="website">Website</label>
    <input type="text" id="website" name="website" tabindex="-1" autocomplete="off">
  </p>
  <p>
    <label for="name">Name</label><br>
    <input type="text" id="name" name="name" required>
  </p>
  <p>
    <label for="email">Email</label><br>
    <input type="email" id="email" name="email" required>
  </p>
  <p>
    <label for="message">Message</label><br>
    <textarea id="message" name="message" rows="5" required></textarea>
  </p>
  <p>
    <button type="submit">Send</button>
  </p>
</form>
