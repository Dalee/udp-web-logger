(function (w) {
    const container = w.document.getElementById('container');

    function fetchLogs() {
        return fetch('/api/log')
            .then(response => response.json())
    }

    setInterval(() => {
        fetchLogs().then(response => {
            if (response.length === 0) {
                return;
            }

            while (container.firstChild) {
                container.removeChild(container.firstChild);
            }

            response.forEach(entry => {
                const el = w.document.createElement('div');
                el.className = 'entry';
                el.innerHTML = `<span class="time">${entry.time}</span> <span class="ip">${entry.ip}</span>: <span class="payload">${entry.payload}</span>`;
                container.appendChild(el);
            });
        });
    }, 3000);
})(window);
