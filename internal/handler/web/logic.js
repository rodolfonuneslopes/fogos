'use strict';

// Pure data-transform functions used by app.js, split out so they can be
// unit tested with Node's test runner without a DOM.

// Lower number = shown first
export const SEVERITY_ORDER = { 5: 1, 4: 2, 7: 3, 9: 4, 8: 5 };

export function sortIncidents(incidents) {
  return [...incidents].sort((a, b) => {
    const sa = SEVERITY_ORDER[a.statusCode] ?? 99;
    const sb = SEVERITY_ORDER[b.statusCode] ?? 99;
    if (sa !== sb) return sa - sb;
    return (b.man + b.terrain + b.aerial) - (a.man + a.terrain + a.aerial);
  });
}

export function filterIncidents(incidents, distrito, concelho) {
  return incidents.filter(inc =>
    (!distrito || inc.district === distrito) &&
    (!concelho || inc.concelho === concelho)
  );
}

export function esc(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
