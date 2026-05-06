const API = 'http://localhost:8080';

// ── Auth guard (runs before DOM) ────────────────────────────────────────────
function logout() {
    localStorage.removeItem('relief_user');
    localStorage.removeItem('relief_user_role');
    window.location.href = 'login.html';
}

const _stored = localStorage.getItem('relief_user');
if (!_stored) window.location.href = 'login.html';

const currentUser = JSON.parse(_stored || '{}');
const userRole    = localStorage.getItem('relief_user_role');

// ── Helpers ─────────────────────────────────────────────────────────────────

const TYPE_CLASS = {
    flood: 'tag-flood', wildfire: 'tag-fire', fire: 'tag-fire',
    earthquake: 'tag-earthquake', cyclone: 'tag-cyclone',
    hurricane: 'tag-cyclone', typhoon: 'tag-cyclone',
    landslide: 'tag-landslide', drought: 'tag-drought'
};

function typeClass(t)   { return TYPE_CLASS[(t || '').toLowerCase()] || 'tag-flood'; }
function sevClass(s)    { const m = {critical:'tag-critical',high:'tag-high',moderate:'tag-moderate',low:'tag-low'}; return m[(s||'').toLowerCase()] || 'tag-low'; }
function statusClass(s) { const m = {active:'status-active',contained:'status-contained',resolved:'status-resolved',monitoring:'status-monitoring'}; return m[(s||'').toLowerCase()] || 'status-active'; }

function maxSeverity(areas) {
    for (const sev of ['Critical','High','Moderate','Low']) {
        if (areas.some(a => a.severity === sev)) return sev;
    }
    return areas.length ? (areas[0].severity || 'Low') : 'Low';
}

function fmtDate(d) {
    if (!d) return '—';
    return new Date(d).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' });
}

// ── Dashboard data ───────────────────────────────────────────────────────────

let allDisasters  = [];
let disasterAreas = {}; // disaster_id → [areas]

async function loadDashboard() {
    try {
        const [disasters, areas, shelters] = await Promise.all([
            fetch(`${API}/api/disasters`).then(r => r.json()),
            fetch(`${API}/api/areas`).then(r => r.json()),
            fetch(`${API}/api/shelters`).then(r => r.json())
        ]);

        // Build area lookup
        disasterAreas = {};
        (areas || []).forEach(a => {
            if (!disasterAreas[a.disaster_id]) disasterAreas[a.disaster_id] = [];
            disasterAreas[a.disaster_id].push(a);
        });

        allDisasters = disasters || [];

        // Stats
        const active   = allDisasters.filter(d => d.status === 'Active').length;
        const totalPop = (areas || []).reduce((s, a) => s + (a.population || 0), 0);
        const totalCap = (shelters || []).reduce((s, sh) => s + (sh.capacity || 0), 0);
        const totalOcc = (shelters || []).reduce((s, sh) => s + (sh.occupied_number || 0), 0);
        const occPct   = totalCap ? Math.round((totalOcc / totalCap) * 100) : 0;

        setText('stat-disasters',  active);
        setText('stat-population', totalPop >= 1000 ? `${(totalPop / 1000).toFixed(0)}k` : totalPop);
        setText('stat-occupancy',  `${occPct}%`);
        setText('stat-shelters',   (shelters || []).length);

        renderDisasters(allDisasters);
        renderShelters(shelters || []);
    } catch {
        setHTML('disaster-grid', '<p class="grid-empty">Could not load data — is the backend running?</p>');
        setHTML('shelter-grid',  '<p class="grid-empty">Could not load data.</p>');
    }
}

function setText(id, val) { const el = document.getElementById(id); if (el) el.textContent = val; }
function setHTML(id, html) { const el = document.getElementById(id); if (el) el.innerHTML = html; }

// ── Disaster rendering ───────────────────────────────────────────────────────

function renderDisasters(disasters) {
    const grid    = document.getElementById('disaster-grid');
    const countEl = document.getElementById('events-count');

    if (!disasters.length) {
        grid.innerHTML = '<p class="grid-empty">No disasters recorded yet.</p>';
        if (countEl) countEl.textContent = '0';
        return;
    }

    if (countEl) countEl.textContent = disasters.length;

    grid.innerHTML = disasters.map(d => {
        const areas    = disasterAreas[d.disaster_id] || [];
        const sev      = maxSeverity(areas);
        const totalPop = areas.reduce((s, a) => s + (a.population || 0), 0);
        const location = areas.length ? areas[0].area_name : '—';

        return `<div class="card disaster-card" data-type="${(d.disaster_type || '').toLowerCase()}">
            <div class="card-header">
                <div class="tags-left">
                    <span class="tag tag-id">DSR-${String(d.disaster_id).padStart(3, '0')}</span>
                    <span class="tag tag-type ${typeClass(d.disaster_type)}">${(d.disaster_type || 'UNKNOWN').toUpperCase()}</span>
                </div>
                <div class="tags-right">
                    <span class="tag tag-severity ${sevClass(sev)}"><span class="dot ${sev.toLowerCase()}"></span> ${sev.toUpperCase()}</span>
                    <span class="tag tag-status ${statusClass(d.status)}">${d.status || 'Unknown'}</span>
                </div>
            </div>
            <h3>${d.disaster_name}</h3>
            <p>${areas.length} affected area${areas.length !== 1 ? 's' : ''} registered.</p>
            <div class="card-meta">
                <span><i class="ph ph-map-pin"></i> ${location}</span>
                <span><i class="ph ph-users"></i> ${totalPop.toLocaleString()} affected</span>
                <span><i class="ph ph-calendar-blank"></i> ${fmtDate(d.start_date)}</span>
            </div>
        </div>`;
    }).join('');
}

// ── Shelter rendering ────────────────────────────────────────────────────────

function renderShelters(shelters) {
    const grid = document.getElementById('shelter-grid');

    if (!shelters.length) {
        grid.innerHTML = '<p class="grid-empty">No shelters registered yet.</p>';
        return;
    }

    grid.innerHTML = shelters.map(s => {
        const pct   = s.capacity ? Math.round((s.occupied_number / s.capacity) * 100) : 0;
        const fill  = pct >= 90 ? 'fill-red' : pct >= 70 ? 'fill-yellow' : 'fill-green';
        const stTag = pct >= 95
            ? '<span class="tag status-full">Full</span>'
            : '<span class="tag status-open">Open</span>';

        return `<div class="card shelter-card">
            <div class="card-header">
                <div class="icon-circle"><i class="ph ph-house"></i></div>
                ${stTag}
            </div>
            <h3>${s.shelter_name}</h3>
            <p class="location">${s.location || '—'}</p>
            <div class="occupancy">
                <div class="occ-text">
                    <span><i class="ph ph-users"></i> ${(s.occupied_number || 0).toLocaleString()} / ${(s.capacity || 0).toLocaleString()}</span>
                    <span class="percentage">${pct}%</span>
                </div>
                <div class="progress-bar"><div class="progress ${fill}" style="width:${pct}%;"></div></div>
            </div>
            <div class="resources">
                <span class="res-tag"><i class="ph ph-phone"></i> ${s.contact_number || 'N/A'}</span>
            </div>
        </div>`;
    }).join('');
}

// ── Type filter ──────────────────────────────────────────────────────────────

function initFilters() {
    document.querySelectorAll('.filter-group').forEach((group, idx) => {
        const btns = group.querySelectorAll('.filter-btn');
        btns.forEach(btn => {
            btn.addEventListener('click', () => {
                btns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                if (idx === 0) {
                    const type = btn.textContent.trim().toLowerCase();
                    renderDisasters(type === 'all'
                        ? allDisasters
                        : allDisasters.filter(d => (d.disaster_type || '').toLowerCase().includes(type))
                    );
                }
            });
        });
    });
}

// ── Donation modal ───────────────────────────────────────────────────────────

function setAmount(val) {
    document.getElementById('donate-amount').value = val;
}

function openDonateModal() {
    document.getElementById('donate-modal').classList.add('open');
    document.getElementById('donate-msg').className = 'donate-msg hidden';
}

function closeDonateModal() {
    document.getElementById('donate-modal').classList.remove('open');
}

function initDonationModal() {
    // Close on backdrop click
    document.getElementById('donate-modal').addEventListener('click', e => {
        if (e.target.id === 'donate-modal') closeDonateModal();
    });

    document.getElementById('donate-form').addEventListener('submit', async e => {
        e.preventDefault();
        const amount = parseFloat(document.getElementById('donate-amount').value);
        if (!amount || amount <= 0) return;

        const msgEl = document.getElementById('donate-msg');
        const btn   = document.getElementById('donate-btn');
        btn.disabled    = true;
        btn.textContent = 'Processing…';

        try {
            const userId = userRole === 'user' ? (currentUser.user_id || null) : null;
            const res    = await fetch(`${API}/api/donations`, {
                method:  'POST',
                headers: { 'Content-Type': 'application/json' },
                body:    JSON.stringify({ amount, user_id: userId })
            });
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || 'Failed');

            msgEl.textContent = `Thank you! $${amount.toLocaleString()} donation recorded.`;
            msgEl.className   = 'donate-msg success';
            document.getElementById('donate-amount').value = '';
            setTimeout(closeDonateModal, 2200);
        } catch (err) {
            msgEl.textContent = err.message;
            msgEl.className   = 'donate-msg error';
        } finally {
            btn.disabled    = false;
            btn.textContent = 'Donate';
        }
    });
}

// ── Chart ────────────────────────────────────────────────────────────────────

function initChart() {
    const canvas = document.getElementById('activityChart');
    if (!canvas) return;
    const ctx = canvas.getContext('2d');

    const gradOrange = ctx.createLinearGradient(0, 0, 0, 400);
    gradOrange.addColorStop(0, 'rgba(249,115,22,0.2)');
    gradOrange.addColorStop(1, 'rgba(249,115,22,0)');

    const gradGreen = ctx.createLinearGradient(0, 0, 0, 400);
    gradGreen.addColorStop(0, 'rgba(34,197,94,0.2)');
    gradGreen.addColorStop(1, 'rgba(34,197,94,0)');

    new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Apr 17','Apr 18','Apr 19','Apr 20','Apr 21','Apr 22','Apr 23'],
            datasets: [
                { label:'Donations', data:[38,45,55,65,80,95,110], borderColor:'#f97316', backgroundColor:gradOrange, borderWidth:2, fill:true, tension:0.4, pointRadius:0 },
                { label:'Relief',    data:[22,30,40,50,60,72,80],  borderColor:'#22c55e', backgroundColor:gradGreen,  borderWidth:2, fill:true, tension:0.4, pointRadius:0 },
                { label:'Disasters', data:[5,6,5,7,6,8,8],         borderColor:'#1c5253', backgroundColor:'transparent', borderWidth:2, tension:0.4, pointRadius:0 }
            ]
        },
        options: {
            responsive: true, maintainAspectRatio: false,
            plugins: {
                legend: { display: false },
                tooltip: { mode:'index', intersect:false, backgroundColor:'#1e293b', titleFont:{family:'Inter',size:13}, bodyFont:{family:'Inter',size:13}, padding:12, cornerRadius:8 }
            },
            scales: {
                x: { grid:{display:false,drawBorder:false}, ticks:{color:'#94a3b8',font:{family:'Inter',size:12}} },
                y: { min:0, max:120, ticks:{stepSize:30,color:'#94a3b8',font:{family:'Inter',size:12}}, grid:{color:'#e2e8f0',borderDash:[5,5],drawBorder:false} }
            },
            interaction: { mode:'nearest', axis:'x', intersect:false }
        }
    });
}

// ── Boot ─────────────────────────────────────────────────────────────────────

document.addEventListener('DOMContentLoaded', () => {
    // Show logged-in user name
    const nameEl = document.getElementById('nav-username');
    if (nameEl) {
        const displayName = userRole === 'admin'
            ? (currentUser.admin_name || 'Admin')
            : (currentUser.user_name || 'User');
        nameEl.textContent = displayName + (userRole === 'admin' ? ' (Admin)' : '');
    }

    loadDashboard();
    initChart();
    initFilters();
    initDonationModal();
});
