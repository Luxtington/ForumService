// Функция для форматирования даты
function formatDate(dateString) {
    console.log('Форматирование даты:', dateString);
    if (!dateString) return '';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        console.log('Некорректная дата:', dateString);
        return '';
    }
    const formatted = date.toLocaleString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit'
    });
    console.log('Отформатированная дата:', formatted);
    return formatted;
}

// Инициализация WebSocket соединения
let ws = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 5000; // 5 секунд
let isConnecting = false;

function initWebSocket() {
    if (isConnecting) {
        console.log('Уже идет попытка подключения...');
        return;
    }

    isConnecting = true;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    console.log('Попытка подключения к WebSocket:', wsUrl);
    
    if (ws) {
        console.log('Закрытие существующего соединения');
        ws.close();
    }
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = function() {
        console.log('WebSocket соединение установлено');
        console.log('WebSocket протокол:', ws.protocol);
        console.log('WebSocket readyState:', ws.readyState);
        reconnectAttempts = 0;
        isConnecting = false;
    };
    
    ws.onmessage = function(event) {
        console.log('Получено сообщение:', event.data);
        try {
            const message = JSON.parse(event.data);
            if (message.type === 'error') {
                console.error('Ошибка от сервера:', message.content);
                alert(message.content);
                return;
            }
            addMessageToChat(message);
        } catch (error) {
            console.error('Ошибка при разборе сообщения:', error);
        }
    };
    
    ws.onclose = function(event) {
        console.log('WebSocket соединение закрыто:', event.code);
        console.log('WebSocket wasClean:', event.wasClean);
        isConnecting = false;
        
        if (reconnectAttempts < maxReconnectAttempts) {
            console.log(`Попытка переподключения ${reconnectAttempts + 1}/${maxReconnectAttempts}...`);
            reconnectAttempts++;
            setTimeout(initWebSocket, reconnectDelay);
        } else {
            console.log('Достигнуто максимальное количество попыток переподключения');
        }
    };
    
    ws.onerror = function(error) {
        console.log('WebSocket ошибка:', error);
        console.log('WebSocket readyState при ошибке:', ws.readyState);
    };
}

// Функция для добавления нового сообщения в чат
function addMessageToChat(message) {
    console.log('Добавление сообщения в чат:', message);
    if (!message || !message.content) {
        console.error('Некорректное сообщение:', message);
        return;
    }

    const chatMessages = document.getElementById('chatMessages');
    if (!chatMessages) {
        console.error('Элемент #chatMessages не найден');
        return;
    }

    const noMessagesMessage = chatMessages.querySelector('.text-center');
    if (noMessagesMessage) {
        noMessagesMessage.remove();
    }
    
    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message';
    const html = `
        <div class="chat-message-meta">
            <i class="bi bi-person-circle"></i> ${message.author_name || 'Аноним'} • ${formatDate(message.created_at)}
        </div>
        <div class="chat-message-content">
            ${message.content || ''}
        </div>
    `;
    console.log('HTML для сообщения:', html);
    messageDiv.innerHTML = html;
    
    chatMessages.appendChild(messageDiv);
    console.log('Сообщение добавлено в чат');
    scrollChatToBottom();
}

// Функция для отправки сообщения в чат
function sendChatMessage(event) {
    event.preventDefault();
    
    const content = document.getElementById('chatMessage').value.trim();
    if (!content) return;
    
    console.log('Отправка сообщения:', content);
    
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.log('WebSocket не подключен, пробуем переподключиться...');
        initWebSocket();
        setTimeout(() => {
            if (ws && ws.readyState === WebSocket.OPEN) {
                sendMessage(content);
            } else {
                alert('Не удалось отправить сообщение. Пожалуйста, обновите страницу.');
            }
        }, 1000);
        return;
    }
    
    sendMessage(content);
}

function sendMessage(content) {
    const message = {
        type: 'message',
        content: content,
        created_at: new Date().toISOString()
    };
    
    try {
        ws.send(JSON.stringify(message));
        document.getElementById('chatMessage').value = '';
    } catch (error) {
        console.error('Ошибка при отправке сообщения:', error);
        alert('Не удалось отправить сообщение. Пожалуйста, попробуйте еще раз.');
    }
}

// Функция для прокрутки чата вниз
function scrollChatToBottom() {
    const chatMessages = document.getElementById('chatMessages');
    if (chatMessages) {
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
}

// Функция для загрузки сообщений чата
function loadChatMessages() {
    fetch('/api/chat')
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка при загрузке сообщений');
            }
            return response.json();
        })
        .then(data => {
            const chatMessages = document.getElementById('chatMessages');
            if (!chatMessages) {
                console.error('Элемент #chatMessages не найден');
                return;
            }
            
            chatMessages.innerHTML = '';
            if (!data || data.length === 0) {
                chatMessages.innerHTML = `
                    <div class="text-center text-muted py-5">
                        <i class="bi bi-chat-square-text display-1"></i>
                        <p class="mt-3">Чат пуст. Напишите первое сообщение!</p>
                    </div>
                `;
            } else {
                data.forEach(message => addMessageToChat(message));
            }
            scrollChatToBottom();
        })
        .catch(error => {
            console.error('Ошибка при загрузке сообщений чата:', error);
            const chatMessages = document.getElementById('chatMessages');
            if (chatMessages) {
                chatMessages.innerHTML = `
                    <div class="text-center text-danger py-5">
                        <i class="bi bi-exclamation-triangle display-1"></i>
                        <p class="mt-3">Ошибка при загрузке сообщений. Попробуйте обновить страницу.</p>
                    </div>
                `;
            }
        });
}

// Загружаем сообщения при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    loadChatMessages();
    initWebSocket();
    scrollChatToBottom();
    
    // Обновляем сообщения каждые 30 секунд
    setInterval(loadChatMessages, 30000);
}); 