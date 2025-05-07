{{define "chat.js"}}
// Функция для отправки сообщения в чат
function sendChatMessage(event) {
    event.preventDefault();
    
    const content = document.getElementById('chatMessage').value.trim();
    if (!content) return;
    
    fetch('/api/chat', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content: content })
    })
    .then(response => {
        if (response.ok) {
            document.getElementById('chatMessage').value = '';
            window.location.reload();
        } else {
            alert('Ошибка при отправке сообщения');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Ошибка при отправке сообщения');
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