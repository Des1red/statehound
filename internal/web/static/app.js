// ── STATE ──────────────────────────────────────────────
const state = {
  offset:        0,
  activeTab:     'all',
  counts:        { critical: 0, normal: 0 },
};

// ── DOM REFS ───────────────────────────────────────────
const eventsContainer = document.getElementById('events-container');
const eventsEmpty     = document.getElementById('events-empty');
const criticalCount   = document.getElementById('critical-count');
const normalCount     = document.getElementById('normal-count');
const snapshotEl      = document.getElementById('snapshot-container');
const refreshBtn      = document.getElementById('refresh-snapshot');

// ── EVENTS ─────────────────────────────────────────────
function parseEventLine(line) {
  // Format: [timestamp] [urgency] TYPE message
  // or:     [timestamp] TYPE message  (no urgency)
  const m = line.match(/^\[([^\]]+)\]\s+(?:\[(critical|normal)\]\s+)?(\S+)\s*(.*)/);
  if (!m) return null;
  return {
    time:    m[1],
    urgency: m[2] || null,
    type:    m[3],
    msg:     m[4] || '',
    raw:     line,
  };
}

function renderEvent(ev, isNew) {
  const el = document.createElement('div');
  el.className = 'event-line' + (ev.urgency ? ' ' + ev.urgency : '') + (isNew ? ' new' : '');
  el.dataset.urgency = ev.urgency || '';

  const time = document.createElement('span');
  time.className = 'event-time';
  time.textContent = formatTime(ev.time);

  const type = document.createElement('span');
  type.className = 'event-type';
  type.textContent = ev.type;

  const msg = document.createElement('span');
  msg.className = 'event-msg';
  msg.textContent = ev.msg;

  el.appendChild(time);
  el.appendChild(type);
  el.appendChild(msg);

  el._ev = ev;
  return el;
}

function formatTime(iso) {
  try {
    const d = new Date(iso);
    return d.toTimeString().slice(0, 8);
  } catch {
    return iso;
  }
}

function applyTabFilter() {
  const lines = eventsContainer.querySelectorAll('.event-line');
  lines.forEach(el => {
    const urgency = el.dataset.urgency;
    if (state.activeTab === 'all') {
      el.classList.remove('hidden');
    } else {
      el.classList.toggle('hidden', urgency !== state.activeTab);
    }
  });
}

function appendEvents(lines, isNew) {
  if (lines.length === 0) return;

  eventsEmpty && eventsEmpty.remove();

  const atBottom = eventsContainer.scrollHeight - eventsContainer.scrollTop
    <= eventsContainer.clientHeight + 40;

  lines.forEach(raw => {
    const ev = parseEventLine(raw);
    if (!ev) return;

    if (ev.urgency === 'critical') state.counts.critical++;
    if (ev.urgency === 'normal')   state.counts.normal++;

    const el = renderEvent(ev, isNew);
    eventsContainer.appendChild(el);

    // apply current tab filter immediately
    if (state.activeTab !== 'all' && el.dataset.urgency !== state.activeTab) {
      el.classList.add('hidden');
    }
  });

  criticalCount.textContent = state.counts.critical;
  normalCount.textContent   = state.counts.normal;

  if (atBottom) {
    eventsContainer.scrollTop = eventsContainer.scrollHeight;
  }
}

async function loadEvents() {
  try {
    const res  = await fetch('/api/events');
    const data = await res.json();
    state.offset = data.offset;
    appendEvents(data.events || [], false);
  } catch (e) {
    console.error('failed to load events', e);
  }
}

async function pollEvents() {
  try {
    const res  = await fetch('/api/events?offset=' + state.offset);
    const data = await res.json();
    if (data.offset !== state.offset) {
      state.offset = data.offset;
      appendEvents(data.events || [], true);
    }
  } catch (e) {
    console.error('poll error', e);
  }
}

// ── SNAPSHOT ───────────────────────────────────────────
function renderSnapshot(text) {
  snapshotEl.innerHTML = '';

  const lines  = text.split('\n');
  let block    = null;
  let blockEl  = null;

  lines.forEach(line => {
    if (line === '') return;

    // section headers (e.g. "Active services:")
    if (!line.startsWith(' ') && !line.startsWith('\t') && line.endsWith(':')) {
      block = document.createElement('div');
      block.className = 'snapshot-block';

      const title = document.createElement('div');
      title.className = 'snapshot-block-title';
      title.textContent = line.slice(0, -1);
      block.appendChild(title);

      blockEl = block;
      snapshotEl.appendChild(block);
      return;
    }

    // summary lines (key=value)
    if (line.includes('=') && line.startsWith('  ')) {
      const el = document.createElement('div');
      el.className = 'snapshot-line summary-line';
      const [key, ...rest] = line.trim().split('=');
      el.innerHTML = key + '=<span>' + rest.join('=') + '</span>';
      if (blockEl) blockEl.appendChild(el);
      else snapshotEl.appendChild(el);
      return;
    }

    // regular lines
    const el = document.createElement('div');
    el.className = 'snapshot-line';
    el.textContent = line;
    if (blockEl) blockEl.appendChild(el);
    else snapshotEl.appendChild(el);
  });
}

async function loadSnapshot() {
  refreshBtn.classList.add('spinning');
  refreshBtn.textContent = '↻ Refreshing...';
  try {
    const res  = await fetch('/api/snapshot');
    const data = await res.json();
    renderSnapshot(data.snapshot || '');
  } catch (e) {
    snapshotEl.innerHTML = '<div class="empty-state">daemon not responding</div>';
  } finally {
    refreshBtn.classList.remove('spinning');
    refreshBtn.textContent = '↻ Refresh';
  }
}

// ── TABS ───────────────────────────────────────────────
document.querySelectorAll('.tab').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    btn.classList.add('active');
    state.activeTab = btn.dataset.tab;
    applyTabFilter();
  });
});

refreshBtn.addEventListener('click', loadSnapshot);

// ── INIT ───────────────────────────────────────────────
loadEvents();
loadSnapshot();
setInterval(pollEvents, 2000);

// ── HUNT ───────────────────────────────────────────────
const modalOverlay = document.getElementById('modal-overlay');
const modalTarget  = document.getElementById('modal-target');
const modalBody    = document.getElementById('modal-body');
const modalClose   = document.getElementById('modal-close');

function extractHuntTarget(ev) {
  const msg = ev.msg;
  // process= field
  const proc = msg.match(/process=(\S+)/);
  if (proc) return proc[1];
  // .service name
  const svc = msg.match(/(\S+\.service)/);
  if (svc) return svc[1];
  // port number from address:port
  const port = msg.match(/:(\d{2,5})(?:\s|$)/);
  if (port) return port[1];
  // exe= field
  const exe = msg.match(/exe=(\S+)/);
  if (exe) return exe[1].split('/').pop();
  // fall back to event type
  return ev.type;
}

function renderHuntResult(text) {
  modalBody.innerHTML = '';
  const lines = text.split('\n').filter(l => l.trim() !== '');

  let section = null;

  lines.forEach(line => {
    if (line.startsWith('hunt target:')) return;

    if (!line.startsWith(' ') && !line.startsWith('\t') && line.trim() !== '') {
      section = document.createElement('div');
      section.className = 'hunt-section';

      const title = document.createElement('div');
      title.className = 'hunt-section-title';
      title.textContent = line.trim();
      section.appendChild(title);
      modalBody.appendChild(section);
      return;
    }

    if (line.trim() === 'no live matches found') {
      const el = document.createElement('div');
      el.className = 'hunt-empty';
      el.textContent = 'no live matches found';
      modalBody.appendChild(el);
      return;
    }

    const el = document.createElement('div');
    el.className = 'hunt-line';
    el.textContent = line;
    if (section) section.appendChild(el);
    else modalBody.appendChild(el);
  });
}

async function hunt(target) {
  modalTarget.textContent = target;
  modalBody.innerHTML = '<div class="empty-state">searching...</div>';
  modalOverlay.classList.add('open');

  try {
    const res  = await fetch('/api/hunt?target=' + encodeURIComponent(target));
    const data = await res.json();
    renderHuntResult(data.result || 'no result');
  } catch (e) {
    modalBody.innerHTML = '<div class="hunt-empty">daemon not responding</div>';
  }
}

function closeModal() {
  modalOverlay.classList.remove('open');
}

modalClose.addEventListener('click', closeModal);
modalOverlay.addEventListener('click', e => {
  if (e.target === modalOverlay) closeModal();
});
document.addEventListener('keydown', e => {
  if (e.key === 'Escape') closeModal();
});

// attach click to event lines
eventsContainer.addEventListener('click', e => {
  const line = e.target.closest('.event-line');
  if (!line) return;
  const ev = line._ev;
  if (!ev) return;
  hunt(extractHuntTarget(ev));
});