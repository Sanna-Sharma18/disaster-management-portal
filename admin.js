const API = 'http://localhost:8080';

// ── Auth guard ──────────────────────────────────────────────────────────────

function logout() {
    localStorage.removeItem('relief_user');
    localStorage.removeItem('relief_user_role');
    window.location.href = 'login.html';
}

const _stored = localStorage.getItem('relief_user');
const _role   = localStorage.getItem('relief_user_role');
if (!_stored || _role !== 'admin') {
    window.location.href = 'login.html';
}
const admin = JSON.parse(_stored || '{}');

// ── Utilities ───────────────────────────────────────────────────────────────

function showMsg(text, type = 'error') {
    const el = document.getElementById('admin-msg');
    el.textContent = text;
    el.className = `admin-msg ${type}`;
    clearTimeout(el._timer);
    el._timer = setTimeout(() => el.className = 'admin-msg hidden', 4000);
}

async function post(endpoint, body) {
    const res  = await fetch(`${API}${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Request failed');
    return data;
}

function notifyDataUpdated() {
    localStorage.setItem('relief_data_updated', Date.now());
}

// ── Disaster form ────────────────────────────────────────────────────────────

document.getElementById('disaster-form').addEventListener('submit', async e => {
    e.preventDefault();
    try {
        const dateVal = document.getElementById('d-date').value;
        const data = await post('/api/disasters', {
            disaster_name: document.getElementById('d-name').value.trim(),
            disaster_type: document.getElementById('d-type').value.trim(),
            start_date:    dateVal ? new Date(dateVal + 'T00:00:00Z').toISOString() : new Date().toISOString(),
            status:        document.getElementById('d-status').value
        });
        document.getElementById('d-id-display').value = `DSR-${String(data.disaster_id).padStart(3, '0')}`;
        showMsg(`Disaster added — ID ${data.disaster_id}`, 'success');
        notifyDataUpdated();
        document.getElementById('d-name').value = '';
        document.getElementById('d-type').value = '';
        document.getElementById('d-date').value = '';
    } catch (err) {
        showMsg(err.message);
    }
});

// ── Affected Area form ───────────────────────────────────────────────────────

document.getElementById('area-form').addEventListener('submit', async e => {
    e.preventDefault();
    try {
        const data = await post('/api/areas', {
            area_name:   document.getElementById('a-name').value.trim(),
            severity:    document.getElementById('a-severity').value,
            population:  parseInt(document.getElementById('a-population').value) || 0,
            disaster_id: parseInt(document.getElementById('a-disaster-id').value)
        });
        document.getElementById('a-id-display').value = data.area_id;
        showMsg(`Area added — ID ${data.area_id}`, 'success');
        notifyDataUpdated();
        document.getElementById('a-name').value        = '';
        document.getElementById('a-population').value  = '';
        document.getElementById('a-disaster-id').value = '';
    } catch (err) {
        showMsg(err.message);
    }
});

// ── Shelter form ─────────────────────────────────────────────────────────────

document.getElementById('shelter-form').addEventListener('submit', async e => {
    e.preventDefault();
    try {
        const data = await post('/api/shelters', {
            shelter_name:    document.getElementById('s-name').value.trim(),
            location:        document.getElementById('s-location').value.trim(),
            capacity:        parseInt(document.getElementById('s-capacity').value) || 0,
            occupied_number: parseInt(document.getElementById('s-occupied').value) || 0,
            contact_number:  document.getElementById('s-contact').value.trim(),
            area_id:         parseInt(document.getElementById('s-area-id').value)
        });
        document.getElementById('s-id-display').value = data.shelter_id;
        showMsg(`Shelter added — ID ${data.shelter_id}`, 'success');
        notifyDataUpdated();
        document.getElementById('s-name').value    = '';
        document.getElementById('s-location').value = '';
        document.getElementById('s-capacity').value = '';
        document.getElementById('s-occupied').value = '0';
        document.getElementById('s-contact').value  = '';
        document.getElementById('s-area-id').value  = '';
    } catch (err) {
        showMsg(err.message);
    }
});

// ── Distribution form ────────────────────────────────────────────────────────

document.getElementById('dist-form').addEventListener('submit', async e => {
    e.preventDefault();
    try {
        const adminId = admin.admin_id || null;
        const data = await post('/api/distributions', {
            material_name: document.getElementById('dist-material').value.trim(),
            quantity:      parseInt(document.getElementById('dist-qty').value) || 0,
            area_id:       parseInt(document.getElementById('dist-area-id').value),
            admin_id:      adminId
        });
        document.getElementById('dist-id-display').value = data.distribution_id;
        showMsg(`Distribution added — ID ${data.distribution_id}`, 'success');
        notifyDataUpdated();
        document.getElementById('dist-material').value = '';
        document.getElementById('dist-qty').value      = '';
        document.getElementById('dist-area-id').value  = '';
    } catch (err) {
        showMsg(err.message);
    }
});

// ── Donations table ──────────────────────────────────────────────────────────

async function loadDonations() {
    const tbody = document.getElementById('donations-tbody');
    try {
        const [donations, users] = await Promise.all([
            fetch(`${API}/api/donations`).then(r => r.json()),
            fetch(`${API}/api/users`).then(r => r.json())
        ]);

        const userMap = {};
        (users || []).forEach(u => { userMap[u.user_id] = u.user_name; });

        if (!donations || !donations.length) {
            tbody.innerHTML = '<tr><td colspan="5" class="table-empty">No donations yet</td></tr>';
            return;
        }

        tbody.innerHTML = donations.map(d => {
            const name = d.user_id ? (userMap[d.user_id] || `User #${d.user_id}`) : 'Anonymous';
            const date = d.donation_date ? d.donation_date.slice(0, 10) : '—';
            const amt  = Number(d.amount).toLocaleString('en-US', { minimumFractionDigits: 2 });
            return `<tr>
                <td>${d.donation_id}</td>
                <td>${name}</td>
                <td>MONEY</td>
                <td class="amount">$${amt}</td>
                <td>${date}</td>
            </tr>`;
        }).join('');
    } catch {
        tbody.innerHTML = '<tr><td colspan="5" class="table-empty">Failed to load</td></tr>';
    }
}

// ── Populate dropdowns ───────────────────────────────────────────────────────

async function loadDropdowns() {
    const [disasters, areas] = await Promise.all([
        fetch(`${API}/api/disasters`).then(r => r.json()).catch(() => []),
        fetch(`${API}/api/areas`).then(r => r.json()).catch(() => [])
    ]);

    const disasterSel = document.getElementById('a-disaster-id');
    (disasters || []).forEach(d => {
        const opt = document.createElement('option');
        opt.value = d.disaster_id;
        opt.textContent = `#${d.disaster_id} — ${d.disaster_name}`;
        disasterSel.appendChild(opt);
    });

    ['s-area-id', 'dist-area-id'].forEach(selId => {
        const sel = document.getElementById(selId);
        (areas || []).forEach(a => {
            const opt = document.createElement('option');
            opt.value = a.area_id;
            opt.textContent = `#${a.area_id} — ${a.area_name}`;
            sel.appendChild(opt);
        });
    });
}

loadDropdowns();
loadDonations();
