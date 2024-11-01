function getGaugeColor(score) {
    if (score >= 0 && score <= 1) return 'linear-gradient(to bottom, #8B0000, #B22222)'; // Dark to slightly lighter red
    if (score === 2) return 'linear-gradient(to bottom, #B22222, #FF6347)'; // Slightly lighter red to tomato red
    if (score >= 3 && score <= 4) return 'linear-gradient(to bottom, #FF8C00, #FFD700)'; // Dark yellow to light yellow
    if (score === 5) return 'linear-gradient(to bottom, #FFA500, #FFD700)'; // Dark yellow to light yellow
    if (score === 6) return 'linear-gradient(to bottom, #FFD700, #ADFF2F)'; // Light yellow to green yellow
    if (score === 7) return 'linear-gradient(to bottom, #98FB98, #32CD32)'; // Light green to lime green
    if (score >= 8 && score <= 9) return 'linear-gradient(to bottom, #32CD32, #228B22)'; // Lime green to forest green
    if (score === 10) return 'linear-gradient(to bottom, #006400, #228B22)'; // Dark green to forest green
    return '#44c767';
}