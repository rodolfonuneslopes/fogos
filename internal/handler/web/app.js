'use strict';

import { sortIncidents, filterIncidents, esc } from './logic.js';

// Maps to CSS classes in styles.css rather than inline styles, since the
// Content-Security-Policy (default-src 'self') blocks inline style attributes.
const STATUS_CLASSES = {
  4: 'status-alert',     // Despacho de 1º Alerta
  5: 'status-active',    // Em Curso
  7: 'status-resolving', // Em Resolução
  8: 'status-concluded', // Conclusão
  9: 'status-watch',     // Vigilância
};

const distritоEl = document.getElementById('distrito');
const concelhoEl = document.getElementById('concelho');
const statusEl  = document.getElementById('status');
const incidentsEl = document.getElementById('incidents');

let allIncidents = [];
let refreshTimer = null;
const REFRESH_INTERVAL = 60_000;

async function fetchIncidents() {
  const res = await fetch('/api/incidents');
  if (!res.ok) throw new Error(`Server error: ${res.status}`);
  return res.json();
}

function filteredIncidents() {
  return filterIncidents(allIncidents, distritоEl.value, concelhoEl.value);
}

function populateDistrito() {
  const current = distritоEl.value;
  const districts = [...new Set(allIncidents.map(i => i.district).filter(Boolean))].sort();
  distritоEl.innerHTML = '<option value="">Todos os distritos</option>';
  for (const d of districts) {
    const opt = document.createElement('option');
    opt.value = d;
    opt.textContent = d;
    if (d === current) opt.selected = true;
    distritоEl.appendChild(opt);
  }
}

function populateConcelho() {
  const distrito = distritоEl.value;
  const current = concelhoEl.value;
  const source = distrito
    ? allIncidents.filter(i => i.district === distrito)
    : allIncidents;
  const concelhos = [...new Set(source.map(i => i.concelho).filter(Boolean))].sort();
  concelhoEl.innerHTML = '<option value="">Todos os concelhos</option>';
  for (const c of concelhos) {
    const opt = document.createElement('option');
    opt.value = c;
    opt.textContent = c;
    if (c === current) opt.selected = true;
    concelhoEl.appendChild(opt);
  }
  concelhoEl.disabled = concelhos.length === 0;
}

function render() {
  const sorted = sortIncidents(filteredIncidents());
  statusEl.hidden = true;
  statusEl.removeAttribute('aria-busy');
  incidentsEl.innerHTML = '';

  if (sorted.length === 0) {
    statusEl.textContent = 'Sem fogos ativos.';
    statusEl.hidden = false;
    return;
  }

  for (const inc of sorted) {
    const card = document.createElement('article');
    const meiosMeta = (inc.man || inc.terrain || inc.aerial)
      ? `🧑‍🚒 ${inc.man} · 🚒 ${inc.terrain} · 🚁 ${inc.aerial}`
      : '';
    const title = concelhoEl.value
      ? esc(inc.freguesia || '—')
      : esc((inc.concelho ? inc.concelho + ' — ' : '') + (inc.freguesia || '—'));
    card.innerHTML = `
      <header>
        <h3>${title}</h3>
      </header>
      <p class="meta meios">
        <span class="meios-left">
          <span class="status-badge ${STATUS_CLASSES[inc.statusCode] || 'status-unknown'}">${esc(inc.status || '—')}</span>
          <span>${meiosMeta}</span>
        </span>
      </p>
      ${inc.id ? `<a class="source-link" href="https://fogos.pt/pt/fogo/${encodeURIComponent(inc.id)}/detalhe" target="_blank" rel="noopener noreferrer">Ver mais informações ↗</a>` : ''}
    `;
    incidentsEl.appendChild(card);
  }
}

function showLoading() {
  incidentsEl.innerHTML = '';
  statusEl.setAttribute('aria-busy', 'true');
  statusEl.textContent = 'A carregar...';
  statusEl.hidden = false;
}

function showError(msg) {
  incidentsEl.innerHTML = '';
  statusEl.removeAttribute('aria-busy');
  statusEl.textContent = msg;
  statusEl.hidden = false;
}

async function refresh() {
  try {
    allIncidents = await fetchIncidents();
    populateDistrito();
    populateConcelho();
    render();
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

async function init() {
  showLoading();
  distritоEl.addEventListener('change', () => {
    concelhoEl.value = '';
    populateConcelho();
    render();
  });
  concelhoEl.addEventListener('change', render);
  await refresh();
  scheduleRefresh();
}

init();
