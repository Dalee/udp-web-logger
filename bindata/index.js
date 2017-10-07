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
    const entriesContainer = w.document.getElementById('entries');
    const lastUpdateContainer = w.document.getElementById('last_update');
    const startButton = w.document.getElementById('start');
    const stopButton = w.document.getElementById('stop');
    const clearButton = w.document.getElementById('clear');

    /**
     * Stops interval, disables the stop button
     * and enables the start.
     */
    function onStopClick() {
        clearInterval(timer);

        stopButton.disabled = true;
        startButton.disabled = false;
    }

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
    function onStartClick() {
        fetchLogEntries();
        timer = setInterval(fetchLogEntries, 2000);

        stopButton.disabled = false;
        startButton.disabled = true;
    }

    stopButton.addEventListener('click', onStopClick);
    startButton.addEventListener('click', onStartClick);
    clearButton.addEventListener('click', onClearClick);

    startButton.dispatchEvent(new Event('click'));

    /**
     * Generates markup for an entry.
     *
     * @param {LogEntry} entry
     * @returns {String}
     */
    function getEntryContentMarkup(entry) {
        return `
            <div>
                <div class="payload">${entry.payload}</div>
                <div class="meta">${entry.time} ${entry.ip}</div>
            </div>
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
            const el = w.document.createElement('div');
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
