// Light/dark theme toggle (the same file is used by alertas).
//
// Loaded as a classic script in <head>, not a module and not deferred, so a
// saved choice is applied before the first paint (no flash of the other
// theme); the CSP forbids inline scripts, hence a file. Without a saved
// choice there is no data-theme attribute and Pico follows the device.
(() => {
  const KEY = 'theme';
  const root = document.documentElement;
  const media = window.matchMedia('(prefers-color-scheme: dark)');

  // localStorage can throw (private mode, blocked storage): best effort.
  let saved = null;
  try { saved = localStorage.getItem(KEY); } catch { /* ignore */ }
  if (saved === 'light' || saved === 'dark') root.dataset.theme = saved;

  const current = () => root.dataset.theme || (media.matches ? 'dark' : 'light');

  // The button shows the mode it switches to.
  function sync(btn) {
    const dark = current() === 'dark';
    const label = dark ? 'Mudar para modo claro' : 'Mudar para modo escuro';
    btn.textContent = dark ? '☀️' : '🌙';
    btn.setAttribute('aria-label', label);
    btn.title = label;
  }

  document.addEventListener('DOMContentLoaded', () => {
    const btn = document.getElementById('theme-toggle');
    if (!btn) return;
    sync(btn);
    btn.hidden = false; // hidden in the HTML so it never shows without JS
    btn.addEventListener('click', () => {
      root.dataset.theme = current() === 'dark' ? 'light' : 'dark';
      try { localStorage.setItem(KEY, root.dataset.theme); } catch { /* ignore */ }
      sync(btn);
    });
    // While the user hasn't chosen, follow the device as it changes.
    media.addEventListener('change', () => sync(btn));
  });
})();
