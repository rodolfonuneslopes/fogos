import { test } from 'node:test';
import assert from 'node:assert/strict';
import { sortIncidents, groupByStatus, filterIncidents, esc } from '../../internal/handler/web/logic.js';

test('sortIncidents orders by severity first', () => {
  const incidents = [
    { statusCode: 8, man: 0, terrain: 0, aerial: 0 }, // Conclusão
    { statusCode: 5, man: 0, terrain: 0, aerial: 0 }, // Em Curso
    { statusCode: 9, man: 0, terrain: 0, aerial: 0 }, // Vigilância
  ];
  const sorted = sortIncidents(incidents);
  assert.deepEqual(sorted.map(i => i.statusCode), [5, 9, 8]);
});

test('sortIncidents breaks ties by total meios, descending', () => {
  const incidents = [
    { statusCode: 5, man: 1, terrain: 0, aerial: 0 },
    { statusCode: 5, man: 10, terrain: 0, aerial: 0 },
    { statusCode: 5, man: 5, terrain: 0, aerial: 0 },
  ];
  const sorted = sortIncidents(incidents);
  assert.deepEqual(sorted.map(i => i.man), [10, 5, 1]);
});

test('sortIncidents treats unknown status codes as lowest priority', () => {
  const incidents = [
    { statusCode: 999, man: 0, terrain: 0, aerial: 0 },
    { statusCode: 5, man: 0, terrain: 0, aerial: 0 },
  ];
  const sorted = sortIncidents(incidents);
  assert.deepEqual(sorted.map(i => i.statusCode), [5, 999]);
});

test('sortIncidents does not mutate the input array', () => {
  const incidents = [
    { statusCode: 8, man: 0, terrain: 0, aerial: 0 },
    { statusCode: 5, man: 0, terrain: 0, aerial: 0 },
  ];
  const original = [...incidents];
  sortIncidents(incidents);
  assert.deepEqual(incidents, original);
});

test('groupByStatus groups incidents sharing a statusCode together', () => {
  const incidents = [
    { statusCode: 5, status: 'Em Curso' },
    { statusCode: 9, status: 'Vigilância' },
    { statusCode: 5, status: 'Em Curso' },
  ];
  const groups = groupByStatus(incidents);
  assert.equal(groups.length, 2);
  assert.equal(groups[0].statusCode, 5);
  assert.equal(groups[0].incidents.length, 2);
  assert.equal(groups[1].statusCode, 9);
  assert.equal(groups[1].incidents.length, 1);
});

test('groupByStatus preserves first-seen order of statuses', () => {
  const incidents = [
    { statusCode: 8, status: 'Conclusão' },
    { statusCode: 5, status: 'Em Curso' },
    { statusCode: 9, status: 'Vigilância' },
  ];
  const groups = groupByStatus(incidents);
  assert.deepEqual(groups.map(g => g.statusCode), [8, 5, 9]);
});

test('filterIncidents returns all incidents when no filters set', () => {
  const incidents = [
    { district: 'Lisboa', concelho: 'Sintra' },
    { district: 'Porto', concelho: 'Gondomar' },
  ];
  assert.equal(filterIncidents(incidents, '', '').length, 2);
});

test('filterIncidents filters by district', () => {
  const incidents = [
    { district: 'Lisboa', concelho: 'Sintra' },
    { district: 'Porto', concelho: 'Gondomar' },
  ];
  const result = filterIncidents(incidents, 'Porto', '');
  assert.equal(result.length, 1);
  assert.equal(result[0].concelho, 'Gondomar');
});

test('filterIncidents filters by district and concelho together', () => {
  const incidents = [
    { district: 'Lisboa', concelho: 'Sintra' },
    { district: 'Lisboa', concelho: 'Cascais' },
  ];
  const result = filterIncidents(incidents, 'Lisboa', 'Cascais');
  assert.equal(result.length, 1);
  assert.equal(result[0].concelho, 'Cascais');
});

test('esc escapes HTML special characters', () => {
  assert.equal(esc(`<script>&"'</script>`), '&lt;script&gt;&amp;&quot;\'&lt;/script&gt;');
});

test('esc coerces non-string input', () => {
  assert.equal(esc(42), '42');
});
