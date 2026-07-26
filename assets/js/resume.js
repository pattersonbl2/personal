(function () {
  var form = document.getElementById('resume-download-form');
  var btn = document.getElementById('resume-submit');
  var statusEl = document.getElementById('resume-status');
  if (!form || !btn) return;

  function setReady(ready) {
    btn.disabled = !ready;
    if (statusEl) {
      statusEl.hidden = ready;
    }
  }

  // Global callbacks referenced by data-* on the widget.
  window.onResumeTurnstileSuccess = function () {
    setReady(true);
  };
  window.onResumeTurnstileReset = function () {
    setReady(false);
  };

  form.addEventListener('submit', function (event) {
    var tokenInput = form.querySelector('textarea[name="cf-turnstile-response"], input[name="cf-turnstile-response"]');
    var token = tokenInput && tokenInput.value ? tokenInput.value.trim() : '';
    if (!token) {
      event.preventDefault();
      setReady(false);
      if (statusEl) {
        statusEl.hidden = false;
        statusEl.textContent = 'Complete the check above, then download.';
      }
    }
  });
})();
