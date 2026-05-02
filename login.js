document.addEventListener('DOMContentLoaded', () => {
    const tabBtns = document.querySelectorAll('.tab-btn');
    const authForms = document.querySelectorAll('.auth-form');

    tabBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            // Remove active class from all tabs
            tabBtns.forEach(b => b.classList.remove('active'));
            // Add active class to clicked tab
            btn.classList.add('active');

            // Hide all forms
            authForms.forEach(form => form.classList.remove('active-form'));
            
            // Show target form
            const targetId = btn.getAttribute('data-target');
            document.getElementById(targetId).classList.add('active-form');
        });
    });

    // Handle form submissions (mock)
    document.getElementById('user-form').addEventListener('submit', (e) => {
        e.preventDefault();
        // Mock login delay
        const btn = e.target.querySelector('button[type="submit"]');
        const originalText = btn.innerHTML;
        btn.innerHTML = '<i class="ph ph-spinner ph-spin"></i> Logging in...';
        btn.disabled = true;

        setTimeout(() => {
            window.location.href = 'index.html'; // Redirect to dashboard
        }, 800);
    });

    document.getElementById('admin-form').addEventListener('submit', (e) => {
        e.preventDefault();
        // Mock login delay
        const btn = e.target.querySelector('button[type="submit"]');
        const originalText = btn.innerHTML;
        btn.innerHTML = '<i class="ph ph-spinner ph-spin"></i> Verifying Credentials...';
        btn.disabled = true;

        setTimeout(() => {
            btn.innerHTML = originalText;
            btn.disabled = false;
            alert('Admin dashboard will be built later!');
        }, 1200);
    });
});
