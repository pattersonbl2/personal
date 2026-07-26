(function () {
  var form = document.getElementById('resume-download-form');
  var btn = document.getElementById('resume-submit');
<<<<<<< HEAD
  var label = btn ? btn.querySelector('.resume-btn-label') : null;
  var statusEl = document.getElementById('resume-status');
  var widget = document.getElementById('resume-turnstile');
  if (!form || !btn || !widget) return;

  var pendingSubmit = false;
  var defaultLabel = label ? label.textContent : 'Download PDF';

  function setBusy(busy) {
    btn.disabled = busy;
    if (label) label.textContent = busy ? 'Preparing…' : defaultLabel;
  }

  function showError(msg) {
    if (!statusEl) return;
    statusEl.hidden = !msg;
    statusEl.textContent = msg || '';
  }

  function getToken() {
    var tokenInput = form.querySelector(
      'textarea[name="cf-turnstile-response"], input[name="cf-turnstile-response"]'
    );
    return tokenInput && tokenInput.value ? tokenInput.value.trim() : '';
  }

  window.onResumeTurnstileSuccess = function () {
    if (!pendingSubmit) return;
    pendingSubmit = false;
    setBusy(false);
    showError('');
    // Native submit skips our listener's preventDefault path once token exists.
    HTMLFormElement.prototype.submit.call(form);
  };

  window.onResumeTurnstileError = function () {
    pendingSubmit = false;
    setBusy(false);
    showError('Verification failed. Please try again.');
    if (window.turnstile && typeof window.turnstile.reset === 'function') {
      try { window.turnstile.reset(widget); } catch (e) { /* ignore */ }
    }
  };

  form.addEventListener('submit', function (event) {
    if (getToken()) return; // already verified — allow POST

    event.preventDefault();
    showError('');

    if (!window.turnstile || typeof window.turnstile.execute !== 'function') {
      showError('Verification is still loading. Try again in a moment.');
      return;
    }

    pendingSubmit = true;
    setBusy(true);
    try {
      window.turnstile.execute(widget);
    } catch (e) {
      pendingSubmit = false;
      setBusy(false);
      showError('Verification failed. Please try again.');
=======
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
>>>>>>> origin/main
    }
  });
})();
