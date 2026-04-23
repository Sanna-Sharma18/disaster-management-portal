const tabs = document.querySelectorAll(".tab");
const forms = document.querySelectorAll(".form");

const registerBtn = document.getElementById("registerBtn");
const backBtn = document.getElementById("backToLogin");
const registerForm = document.getElementById("registerForm");

// Switch Tabs
tabs.forEach(tab => {
    tab.addEventListener("click", () => {
        tabs.forEach(t => t.classList.remove("active"));
        forms.forEach(f => f.classList.remove("active"));

        tab.classList.add("active");
        document.getElementById(tab.dataset.role + "Form").classList.add("active");
    });
});

// Toggle Password
document.querySelectorAll(".toggle").forEach(icon => {
    icon.addEventListener("click", () => {
        const input = document.getElementById(icon.dataset.target);
        input.type = input.type === "password" ? "text" : "password";
    });
});

// Login Forms
document.querySelectorAll("form").forEach(form => {
    if (form.id !== "registerForm") {
        form.addEventListener("submit", e => {
            e.preventDefault();

            const email = form.querySelector("input[type='email']").value;
            const password = form.querySelector("input[type='password']").value;

            if (!email || !password) {
                alert("Please fill all fields");
                return;
            }

            if (!email.includes("@")) {
                alert("Enter valid email");
                return;
            }

            if (form.id === "adminForm") {
                alert("Redirecting to Admin Dashboard...");
            } 
            else if (form.id === "userForm") {
                alert("Redirecting to User Home...");
            } 
            else {
                alert("Redirecting to Volunteer Dashboard...");
            }
        });
    }
});

// Open Registration
registerBtn.addEventListener("click", (e) => {
    e.preventDefault();

    tabs.forEach(t => t.classList.remove("active"));
    forms.forEach(f => f.classList.remove("active"));

    registerForm.classList.add("active");
});

// Back to Login
backBtn.addEventListener("click", (e) => {
    e.preventDefault();

    forms.forEach(f => f.classList.remove("active"));
    tabs[2].classList.add("active");

    document.getElementById("volunteerForm").classList.add("active");
});

// Register Submit
registerForm.addEventListener("submit", (e) => {
    e.preventDefault();

    const inputs = registerForm.querySelectorAll("input, select");

    for (let input of inputs) {
        if (!input.value) {
            alert("Please fill all fields");
            return;
        }
    }

    alert("Registration Successful!");

    registerForm.classList.remove("active");
    document.getElementById("volunteerForm").classList.add("active");
    tabs[2].classList.add("active");
});