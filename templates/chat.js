{{define "chat.js"}}
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

// Функция для добавления нового сообщения в чат
function addMessageToChat(message) {
    console.log('Добавление сообщения в чат:', message);
    if (!message || !message.id) {
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
            <i class="bi bi-person-circle"></i> ID: ${message.author_id || ''} • ${formatDate(message.created_at)}
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
    
    fetch('/api/chat', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content: content })
    })
    .then(response => {
        console.log('Получен ответ:', response.status);
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Ошибка при отправке сообщения');
        }
    })
    .then(message => {
        console.log('Получено сообщение:', message);
        addMessageToChat(message);
        document.getElementById('chatMessage').value = '';
    })
    .catch(error => {
        console.error('Error:', error);
        alert(error.message);
    });
}

// Функция для прокрутки чата вниз
function scrollChatToBottom() {
    const chatMessages = document.getElementById('chatMessages');
    if (chatMessages) {
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
}

// Прокручиваем чат вниз при загрузке страницы
document.addEventListener('DOMContentLoaded', scrollChatToBottom);
{{end}} 