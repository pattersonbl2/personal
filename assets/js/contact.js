(function () {
  var form = document.getElementById('contact-form');
  if (!form) return;

  var statusEl = document.getElementById('contact-status');
  var submitBtn = document.getElementById('contact-submit');
  var tsInput = document.getElementById('form_ts');
  if (tsInput) {
    tsInput.value = String(Math.floor(Date.now() / 1000));
  }

  function setStatus(message, kind) {
    if (!statusEl) return;
    statusEl.hidden = !message;
    statusEl.textContent = message || '';
    statusEl.className = 'contact-status' + (kind ? ' contact-status--' + kind : '');
  }

  function setBusy(busy) {
    if (!submitBtn) return;
    submitBtn.disabled = busy;
    submitBtn.setAttribute('aria-busy', busy ? 'true' : 'false');
    submitBtn.textContent = busy ? 'Sending…' : 'Send';
  }

  form.addEventListener('submit', function (event) {
    event.preventDefault();
    setStatus('', '');
    setBusy(true);

    if (tsInput && !tsInput.value) {
      tsInput.value = String(Math.floor(Date.now() / 1000));
    }

    // urlencoded (not FormData/multipart) so Go's ParseForm and CORS preflight stay simple
    var data = new URLSearchParams(new FormData(form));

    fetch(form.action, {
      method: 'POST',
      body: data,
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        'X-Requested-With': 'XMLHttpRequest'
      },
      credentials: 'omit'
    })
      .then(function (res) {
        return res.json().then(function (payload) {
          return { res: res, payload: payload };
        }).catch(function () {
          return { res: res, payload: null };
        });
      })
      .then(function (result) {
        var payload = result.payload || {};
        if (result.res.ok && payload.ok) {
          setStatus(payload.message || 'Message received. Thanks!', 'ok');
          form.reset();
          if (tsInput) tsInput.value = String(Math.floor(Date.now() / 1000));
          if (window.turnstile && typeof window.turnstile.reset === 'function') {
            try { window.turnstile.reset(); } catch (e) { /* ignore */ }
          }
          return;
        }
        var err = (payload && payload.error) || 'Something went wrong. Please try again.';
        setStatus(err, 'error');
        if (window.turnstile && typeof window.turnstile.reset === 'function') {
          try { window.turnstile.reset(); } catch (e) { /* ignore */ }
        }
      })
      .catch(function (err) {
        var detail = err && err.message ? String(err.message) : '';
        setStatus(
          detail && detail !== 'Failed to fetch'
            ? 'Network error: ' + detail
            : 'Network error. Please try again.',
          'error'
        );
      })
      .finally(function () {
        setBusy(false);
      });
  });
})();
