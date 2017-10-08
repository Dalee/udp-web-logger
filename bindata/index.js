(function (w) {
    /**
     * @typedef {Object} LogEntry
     * @property {String} time
     * @property {String} ip
     * @property {String} payload
     * @property {Number} unix
     */

    let timer = null;
    let lastUnix = 0;
    let isFetching = false;

    const entriesContainer = w.document.getElementById('entries');
    const lastUpdateContainer = w.document.getElementById('last_update');
    const toggleButton = w.document.getElementById('toggle');
    const clearButton = w.document.getElementById('clear');

    /**
     * Clears all the messages.
     */
    function onClearClick() {
        while (entriesContainer.firstChild) {
            entriesContainer.removeChild(entriesContainer.firstChild);
        }
    }

    /**
     * Starts interval, disables the start button
     * and enables the stop.
     */
    function onToggleClick() {
        if (!isFetching) {
            fetchLogEntries();
            timer = setInterval(fetchLogEntries, 2000);

            toggleButton.innerText = toggleButton.dataset.stop;
            toggleButton.classList.remove(toggleButton.dataset.start);
            toggleButton.classList.add(toggleButton.dataset.stop);
        } else {
            clearInterval(timer);

            toggleButton.innerText = toggleButton.dataset.start;
            toggleButton.classList.remove(toggleButton.dataset.stop);
            toggleButton.classList.add(toggleButton.dataset.start);
        }

        isFetching = !isFetching;
    }

    toggleButton.addEventListener('click', onToggleClick);
    clearButton.addEventListener('click', onClearClick);

    toggleButton.dispatchEvent(new Event('click'));

    /**
     * Generates markup for an entry.
     *
     * @param {LogEntry} entry
     * @returns {String}
     */
    function getEntryContentMarkup(entry) {
        return `
            <td class="log_col_left_pad log_col_pad log_col_nowrap">${entry.time}</td>
            <td class="log_col_pad log_col_nowrap">${entry.ip}</td>
            <td>${entry.payload}</td>
        `;
    }

    /**
     * Renders log entries to the DOM.
     *
     * @param {Array.<LogEntry>} entries
     */
    function renderLogEntries(entries) {
        if (entries.length === 0) {
            return;
        }

        lastUnix = entries[entries.length - 1].unix;

        entries.forEach(entry => {
            const el = w.document.createElement('tr');
            el.className = 'entry';
            el.innerHTML = getEntryContentMarkup(entry);
            entriesContainer.prepend(el);
        });
    }

    /**
     * Fetches log entries from backend.
     */
    function fetchLogEntries() {
        fetch(`/api/log?ts=${lastUnix}`)
            .then(response => response.json())
            .then(renderLogEntries)
            .then(() => {
                const lastUpdateTime = (new Date).toLocaleString();
                lastUpdateContainer.innerHTML = `Last update: ${lastUpdateTime}`;
            });
    }
})(window);
