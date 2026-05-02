document.addEventListener('DOMContentLoaded', () => {
    // 1. Initialize Chart.js for Activity Trend
    const ctx = document.getElementById('activityChart').getContext('2d');

    // Gradient for Donations (Orange)
    const gradientOrange = ctx.createLinearGradient(0, 0, 0, 400);
    gradientOrange.addColorStop(0, 'rgba(249, 115, 22, 0.2)');
    gradientOrange.addColorStop(1, 'rgba(249, 115, 22, 0)');

    // Gradient for Relief (Green)
    const gradientGreen = ctx.createLinearGradient(0, 0, 0, 400);
    gradientGreen.addColorStop(0, 'rgba(34, 197, 94, 0.2)');
    gradientGreen.addColorStop(1, 'rgba(34, 197, 94, 0)');

    const activityChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: ['Apr 17', 'Apr 18', 'Apr 19', 'Apr 20', 'Apr 21', 'Apr 22', 'Apr 23'],
            datasets: [
                {
                    label: 'Donations',
                    data: [38, 45, 55, 65, 80, 95, 110],
                    borderColor: '#f97316', // orange
                    backgroundColor: gradientOrange,
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4,
                    pointRadius: 0
                },
                {
                    label: 'Relief',
                    data: [22, 30, 40, 50, 60, 72, 80],
                    borderColor: '#22c55e', // green
                    backgroundColor: gradientGreen,
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4,
                    pointRadius: 0
                },
                {
                    label: 'Disasters',
                    data: [5, 6, 5, 7, 6, 8, 8],
                    borderColor: '#1c5253', // teal
                    backgroundColor: 'transparent',
                    borderWidth: 2,
                    tension: 0.4,
                    pointRadius: 0
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false // We built a custom HTML legend
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    backgroundColor: '#1e293b',
                    titleFont: { family: 'Inter', size: 13 },
                    bodyFont: { family: 'Inter', size: 13 },
                    padding: 12,
                    cornerRadius: 8,
                }
            },
            scales: {
                x: {
                    grid: {
                        display: false,
                        drawBorder: false
                    },
                    ticks: {
                        color: '#94a3b8',
                        font: { family: 'Inter', size: 12 }
                    }
                },
                y: {
                    min: 0,
                    max: 120,
                    ticks: {
                        stepSize: 30,
                        color: '#94a3b8',
                        font: { family: 'Inter', size: 12 }
                    },
                    grid: {
                        color: '#e2e8f0',
                        borderDash: [5, 5],
                        drawBorder: false
                    }
                }
            },
            interaction: {
                mode: 'nearest',
                axis: 'x',
                intersect: false
            }
        }
    });

    // 2. Simple Filter Button Toggle (Visual Only)
    const filterGroups = document.querySelectorAll('.filter-group');
    
    filterGroups.forEach(group => {
        const buttons = group.querySelectorAll('.filter-btn');
        buttons.forEach(btn => {
            btn.addEventListener('click', () => {
                // Remove active class from all buttons in this group
                buttons.forEach(b => b.classList.remove('active'));
                // Add active class to clicked button
                btn.classList.add('active');
            });
        });
    });
});
