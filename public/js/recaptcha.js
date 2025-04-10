document.body.addEventListener('htmx:afterSwap', function(event) {
    if (event.detail.target.id === 'validation-errors') {
        // Check if validation-errors contains an alert-danger div
        if (event.detail.target.querySelector('.alert-danger')) {
            // Reset all reCAPTCHA widgets
            if (typeof grecaptcha !== 'undefined' && grecaptcha) {
                try {
                    grecaptcha.reset();
                } catch (e) {
                    console.error('Failed to reset reCAPTCHA:', e);
                }
            }
        }
    }
});
