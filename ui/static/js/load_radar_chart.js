function loadRadarData(elementId, labels, data) {
    let radarChartData = {
        labels: labels,
        datasets: [{
            label: 'Keep Defence Scores',
            data: data,
            color: '#ffffff',
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            borderColor: 'rgba(54, 162, 235, 1)',
            pointBackgroundColor: 'rgba(54, 162, 235, 1)',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 2,
            pointHoverRadius: 10,
        }]
    };

    let radarChartConfig = {
        type: 'radar',
        data: radarChartData,
        options: {
            scales: {
                r: {
                    beginAtZero: true,
                    max: 2,
                    min: 0,
                    ticks: {
                        stepSize: 1,
                        color: '#ffffff', // Light color for tick labels
                        backdropColor: 'rgba(0, 0, 0, 0)', // Transparent backdrop
                        callback: function (value) {
                            return value.toFixed(1);
                        }
                    },
                    grid: {
                        color: 'rgba(255, 255, 255, 0.2)' // Light color for grid lines
                    },
                    angleLines: {
                        color: 'rgba(255, 255, 255, 0.2)' // Light color for angle lines
                    },
                    pointLabels: {
                        color: '#ffffff', // Light color for point labels
                        font: {
                            size: 12 // Adjust the font size here
                        }
                    }
                }
            },
            plugins: {
                tooltip: {
                    callbacks: {
                        label: function (context) {
                            const label = context.label || '';
                            const value = context.raw;
                            return `${label}\n${value}`;
                        },
                    },
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    titleFont: {
                        size: 16,
                        weight: 'bold',
                        color: '#ffffff' // Light color for tooltip title
                    },
                    bodyFont: {
                        size: 14,
                        color: '#ffffff' // Light color for tooltip body
                    },
                    footerFont: {
                        size: 12,
                        style: 'italic',
                        color: '#ffffff' // Light color for tooltip footer
                    },
                    padding: 10,
                    displayColors: false
                }
            },
            onClick: (event, elements) => {
                if (elements.length > 0) {
                    // Get the index of the clicked element
                    const index = elements[0].index;
                    const defenceName = radarChartData.labels[index];
                    const score = radarChartData.datasets[0].data[index];
                    // Customize the tooltip text
                    const tooltipText = `Defence: ${defenceName}\nScore: ${score}`;
                    alert(tooltipText);
                }
            }
        }
    };

    let radarChartCtx = document.getElementById(elementId).getContext('2d');
    new Chart(radarChartCtx, radarChartConfig);
}
