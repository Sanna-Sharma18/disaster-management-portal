const API = 'http://localhost:8080';

// ── Utilities ──────────────────────────────────────────────────────────────

function showMessage(text, type = 'error') {
    const el = document.getElementById('auth-message');
    el.textContent = text;
    el.className = `auth-message ${type}`;
}

function clearMessage() {
    document.getElementById('auth-message').className = 'auth-message hidden';
}

function setLoading(btn, on, original) {
    btn.disabled = on;
    btn.innerHTML = on ? '<i class="ph ph-spinner ph-spin"></i> Please wait…' : original;
}

function switchView(panel, view) {
    document.querySelectorAll(`#panel-${panel} .auth-view`).forEach(v =>
        v.classList.remove('active-view')
    );
    document.getElementById(`${panel}-${view}`).classList.add('active-view');
    document.querySelectorAll(`#panel-${panel} .auth-tab-btn`).forEach(b =>
        b.classList.toggle('active', b.dataset.view === view)
    );
    clearMessage();
}

// ── Role tab switching (User ↔ Admin) ──────────────────────────────────────

document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        clearMessage();
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        document.querySelectorAll('.role-panel').forEach(p => p.classList.add('hidden'));
        document.getElementById(`panel-${btn.dataset.role}`).classList.remove('hidden');
    });
});

// ── Inner Sign In / Sign Up tab buttons (User panel only) ─────────────────

document.querySelectorAll('.auth-tab-btn').forEach(btn => {
    btn.addEventListener('click', () => switchView(btn.dataset.panel, btn.dataset.view));
});

document.querySelectorAll('.auth-tab-link').forEach(link => {
    link.addEventListener('click', e => {
        e.preventDefault();
        switchView(link.dataset.panel, link.dataset.view);
    });
});

// ── USER SIGN UP ───────────────────────────────────────────────────────────

document.getElementById('user-signup-form').addEventListener('submit', async e => {
    e.preventDefault();
    clearMessage();

    const name     = document.getElementById('up-name').value.trim();
    const email    = document.getElementById('up-email').value.trim();
    const phone    = document.getElementById('up-phone').value.trim();
    const password = document.getElementById('up-password').value;

    if (!name || !email || !password) {
        showMessage('Name, email and password are required.');
        return;
    }

    const btn = document.getElementById('user-signup-btn');
    const orig = btn.innerHTML;
    setLoading(btn, true);

    try {
        const res  = await fetch(`${API}/api/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_name: name, user_email: email, user_phoneno: phone, password })
        });
        const data = await res.json();

        if (!res.ok) {
            showMessage(data.error || 'Sign-up failed. That email may already be registered.');
            return;
        }

        // Clear form, pre-fill sign-in email, switch to sign-in view
        document.getElementById('up-name').value     = '';
        document.getElementById('up-email').value    = '';
        document.getElementById('up-phone').value    = '';
        document.getElementById('up-password').value = '';
        document.getElementById('us-email').value    = email;

        showMessage(`Account created! Welcome, ${data.user_name}. Please sign in.`, 'success');
        switchView('user', 'signin');

    } catch {
        showMessage('Cannot reach the server. Make sure the backend is running.');
    } finally {
        setLoading(btn, false, orig);
    }
});

// ── USER SIGN IN ───────────────────────────────────────────────────────────

document.getElementById('user-signin-form').addEventListener('submit', async e => {
    e.preventDefault();
    clearMessage();

    const email    = document.getElementById('us-email').value.trim();
    const password = document.getElementById('us-password').value;

    if (!email || !password) {
        showMessage('Email and password are required.');
        return;
    }

    const btn = document.getElementById('user-signin-btn');
    const orig = btn.innerHTML;
    setLoading(btn, true);

    try {
        const res  = await fetch(`${API}/api/users/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        const data = await res.json();

        if (!res.ok) {
            showMessage('Incorrect email or password.');
            return;
        }

        localStorage.setItem('relief_user',     JSON.stringify(data));
        localStorage.setItem('relief_user_role', 'user');
        window.location.href = 'index.html';

    } catch {
        showMessage('Cannot reach the server. Make sure the backend is running.');
    } finally {
        setLoading(btn, false, orig);
    }
});

// ── ADMIN SIGN IN ──────────────────────────────────────────────────────────

document.getElementById('admin-signin-form').addEventListener('submit', async e => {
    e.preventDefault();
    clearMessage();

    const email    = document.getElementById('ad-email').value.trim();
    const password = document.getElementById('ad-password').value;

    if (!email || !password) {
        showMessage('Email and password are required.');
        return;
    }

    const btn = document.getElementById('admin-signin-btn');
    const orig = btn.innerHTML;
    setLoading(btn, true);

    try {
        const res  = await fetch(`${API}/api/admins/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        const data = await res.json();

        if (!res.ok) {
            showMessage('Invalid admin credentials.');
            return;
        }

        localStorage.setItem('relief_user',     JSON.stringify(data));
        localStorage.setItem('relief_user_role', 'admin');
        window.location.href = 'admin.html';

    } catch {
        showMessage('Cannot reach the server. Make sure the backend is running.');
    } finally {
        setLoading(btn, false, orig);
    }
});
