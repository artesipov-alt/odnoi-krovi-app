// app.js - Main application logic for Odnoi Krovi API testing page

// Clear results function
function clearResults() {
    const resultsContainer = document.getElementById('results');
    resultsContainer.innerHTML = `
        <div class="text-center py-12 fade-in">
            <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                </svg>
            </div>
            <p class="text-gray-500">Здесь будут отображаться результаты запросов</p>
            <p class="text-gray-400 text-sm mt-2">Нажмите на любую кнопку выше для тестирования API</p>
        </div>
    `;
}

// Format timestamp
function getTimestamp() {
    const now = new Date();
    return now.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
}

// Format JSON with syntax highlighting
function formatJSON(obj) {
    const json = JSON.stringify(obj, null, 2);
    return json
        .replace(/"([^"]+)":/g, '<span class="json-key">"$1"</span>:')
        .replace(/: "([^"]*)"/g, ': <span class="json-string">"$1"</span>')
        .replace(/: (\d+)/g, ': <span class="json-number">$1</span>')
        .replace(/: (true|false)/g, ': <span class="json-boolean">$1</span>')
        .replace(/: null/g, ': <span class="json-null">null</span>');
}

// Create result item element
function createResultItem(type, title, data, isJSON = true) {
    const resultDiv = document.createElement('div');
    resultDiv.className = `result-item ${type}`;

    const header = document.createElement('div');
    header.className = 'result-header';

    const titleSpan = document.createElement('span');
    titleSpan.textContent = title;

    const badge = document.createElement('span');
    badge.className = `badge ${type}`;
    badge.textContent = type === 'success' ? 'Success' : 'Error';

    header.appendChild(titleSpan);
    header.appendChild(badge);

    const dataDiv = document.createElement('div');
    dataDiv.className = 'result-data';

    if (isJSON) {
        try {
            const jsonData = typeof data === 'string' ? JSON.parse(data) : data;
            dataDiv.innerHTML = formatJSON(jsonData);
        } catch (e) {
            dataDiv.textContent = data;
        }
    } else {
        dataDiv.textContent = data;
    }

    const timestamp = document.createElement('div');
    timestamp.className = 'result-timestamp';
    timestamp.textContent = `Время запроса: ${getTimestamp()}`;

    resultDiv.appendChild(header);
    resultDiv.appendChild(dataDiv);
    resultDiv.appendChild(timestamp);

    return resultDiv;
}

// Update status indicator
function updateStatusIndicator(isOnline) {
    const indicator = document.getElementById('status-indicator');
    if (indicator) {
        const dot = indicator.querySelector('div');
        const text = indicator.querySelector('span');

        if (isOnline) {
            dot.className = 'w-3 h-3 rounded-full status-online';
            text.textContent = 'Онлайн';
            text.className = 'text-sm text-green-600 font-medium';
        } else {
            dot.className = 'w-3 h-3 rounded-full status-offline';
            text.textContent = 'Оффлайн';
            text.className = 'text-sm text-red-600 font-medium';
        }
    }
}

// HTMX Event Handlers

// Handle successful requests
document.addEventListener('htmx:afterRequest', function(event) {
    const target = event.detail.target;
    const path = event.detail.requestConfig.path;

    // Special handling for status check
    if (target.id === 'status-result') {
        if (event.detail.successful) {
            updateStatusIndicator(true);
            try {
                const responseData = JSON.parse(event.detail.xhr.responseText);
                const resultItem = createResultItem(
                    'success',
                    `Статус сервера: ${path}`,
                    responseData
                );
                target.innerHTML = '';
                target.appendChild(resultItem);
            } catch (e) {
                target.innerHTML = `
                    <div class="result-item success">
                        <div class="result-header">
                            <span>Сервер доступен</span>
                            <span class="badge success">Success</span>
                        </div>
                        <div class="result-data">${event.detail.xhr.responseText}</div>
                    </div>
                `;
            }
        } else {
            updateStatusIndicator(false);
        }
        return;
    }

    // Handle regular API requests
    if (target && event.detail.successful) {
        try {
            const responseData = JSON.parse(event.detail.xhr.responseText);
            const resultItem = createResultItem(
                'success',
                `Ответ от ${path}`,
                responseData
            );

            // Clear placeholder if exists
            const placeholder = target.querySelector('.text-center.py-12');
            if (placeholder) {
                target.innerHTML = '';
            }

            target.insertBefore(resultItem, target.firstChild);
        } catch (e) {
            // If not JSON, display as plain text
            const resultItem = createResultItem(
                'success',
                `Ответ от ${path}`,
                event.detail.xhr.responseText,
                false
            );

            const placeholder = target.querySelector('.text-center.py-12');
            if (placeholder) {
                target.innerHTML = '';
            }

            target.insertBefore(resultItem, target.firstChild);
        }
    }
});

// Handle request errors
document.addEventListener('htmx:responseError', function(event) {
    const target = event.detail.target;
    const path = event.detail.requestConfig.path;

    // Update status indicator if it's a status check
    if (target.id === 'status-result') {
        updateStatusIndicator(false);
    }

    if (target) {
        const errorMessage = `Статус: ${event.detail.xhr.status} - ${event.detail.xhr.statusText}`;
        const resultItem = createResultItem(
            'error',
            `Ошибка при запросе к ${path}`,
            errorMessage,
            false
        );

        // Clear placeholder if exists
        const placeholder = target.querySelector('.text-center.py-12');
        if (placeholder) {
            target.innerHTML = '';
        }

        target.insertBefore(resultItem, target.firstChild);
    }
});

// Handle network errors
document.addEventListener('htmx:sendError', function(event) {
    const target = event.detail.target;
    const path = event.detail.requestConfig.path;

    if (target.id === 'status-result') {
        updateStatusIndicator(false);
    }

    if (target) {
        const resultItem = createResultItem(
            'error',
            `Ошибка сети при запросе к ${path}`,
            'Не удалось выполнить запрос. Проверьте подключение к серверу.',
            false
        );

        const placeholder = target.querySelector('.text-center.py-12');
        if (placeholder) {
            target.innerHTML = '';
        }

        target.insertBefore(resultItem, target.firstChild);
    }
});

// Add loading state to buttons
document.addEventListener('htmx:beforeRequest', function(event) {
    if (event.detail.elt.tagName === 'BUTTON') {
        event.detail.elt.classList.add('loading');
    }
});

document.addEventListener('htmx:afterRequest', function(event) {
    if (event.detail.elt.tagName === 'BUTTON') {
        event.detail.elt.classList.remove('loading');
    }
});

// Console welcome message
console.log('%c🩸 Одной Крови - API Testing Page', 'font-size: 16px; font-weight: bold; color: #ef4444;');
console.log('%cДобро пожаловать! Эта страница для тестирования API бэкенда.', 'font-size: 12px; color: #6b7280;');
