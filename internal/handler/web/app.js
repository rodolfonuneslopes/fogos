'use strict';

const DEFAULT_CONCELHO = 'Castelo Branco';
const REFRESH_INTERVAL = 60_000;

const select = document.getElementById('concelho');
const statusEl = document.getElementById('status');
const incidentsEl = document.getElementById('incidents');

let refreshTimer = null;

async function loadConcelhos() {
  const res = await fetch('/api/concelhos');
  if (!res.ok) throw new Error('Failed to load concelhos');
  return res.json();
}

async function loadIncidents(concelho) {
  const res = await fetch(`/api/incidents?concelho=${encodeURIComponent(concelho)}`);
  if (!res.ok) throw new Error(`Server error: ${res.status}`);
  return res.json();
}

function renderIncidents(incidents) {
  statusEl.hidden = true;
  statusEl.removeAttribute('aria-busy');
  incidentsEl.innerHTML = '';

  if (!incidents || incidents.length === 0) {
    statusEl.textContent = 'Sem fogos ativos neste concelho.';
    statusEl.hidden = false;
    return;
  }

  for (const inc of incidents) {
    const card = document.createElement('article');
    const meiosMeta = (inc.man || inc.terrain || inc.aerial)
      ? `🧑‍🚒 ${inc.man} · 🚒 ${inc.terrain} · 🚁 ${inc.aerial}`
      : '';
    card.innerHTML = `
      <header>
        <h3>${esc(inc.freguesia || '—')}</h3>
        <span class="status-badge">${esc(inc.status || '—')}</span>
      </header>
      <p class="meta meios">
        <span>${meiosMeta}</span>
        ${inc.date ? `<span class="date">${esc(inc.date)} ${esc(inc.hour)}</span>` : ''}
      </p>
      ${inc.extra ? `<p class="meta extra">${esc(inc.extra)}</p>` : ''}
    `;
    incidentsEl.appendChild(card);
  }
}

function showError(msg) {
  incidentsEl.innerHTML = '';
  statusEl.removeAttribute('aria-busy');
  statusEl.textContent = msg;
  statusEl.hidden = false;
}

function showLoading() {
  incidentsEl.innerHTML = '';
  statusEl.setAttribute('aria-busy', 'true');
  statusEl.textContent = 'A carregar...';
  statusEl.hidden = false;
}

async function refresh() {
  const concelho = select.value;
  if (!concelho) return;
  try {
    const incidents = await loadIncidents(concelho);
    renderIncidents(incidents);
  } catch (err) {
    showError('Erro ao obter dados. Tente novamente.');
    console.error(err);
  }
}

function scheduleRefresh() {
  clearInterval(refreshTimer);
  if (!document.hidden) {
    refreshTimer = setInterval(refresh, REFRESH_INTERVAL);
  }
}

document.addEventListener('visibilitychange', () => {
  if (document.hidden) {
    clearInterval(refreshTimer);
  } else {
    refresh();
    scheduleRefresh();
  }
});

async function onConcelhoChange() {
  showLoading();
  await refresh();
  scheduleRefresh();
}

async function init() {
  try {
    const concelhos = await loadConcelhos();
    for (const c of concelhos) {
      const opt = document.createElement('option');
      opt.value = c;
      opt.textContent = c;
      if (c === DEFAULT_CONCELHO) opt.selected = true;
      select.appendChild(opt);
    }
  } catch (err) {
    showError('Erro ao carregar lista de concelhos.');
    console.error(err);
    return;
  }

  select.addEventListener('change', onConcelhoChange);
  await onConcelhoChange();
}

function esc(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

init();
