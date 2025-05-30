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
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
    console.log('Отформатированная дата:', formatted);
    return formatted;
}

// Функция для редактирования треда
function editThread(id, title) {
    document.getElementById('editThreadId').value = id;
    document.getElementById('editThreadTitle').value = title;
    new bootstrap.Modal(document.getElementById('editThreadModal')).show();
}

// Функция для удаления треда
function deleteThread(id) {
    if (confirm('Вы уверены, что хотите удалить этот тред?')) {
        fetch(`/api/threads/${id}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                loadThreads(); // Перезагружаем список тредов после удаления
            } else {
                alert('Ошибка при удалении треда');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Ошибка при удалении треда');
        });
    }
}

// Функция для добавления нового треда в список
function addThreadToList(thread) {
    console.log('Добавление треда в список:', thread);
    console.log('Debug - Thread Author ID:', thread.author_id);
    console.log('Debug - Window User ID:', window.userId);
    console.log('Debug - Window User Role:', window.userRole);
    console.log('Debug - Is Author:', thread.author_id === window.userId);
    console.log('Debug - Is Admin:', window.userRole === "admin");
    console.log('Debug - Should Show Buttons:', thread.author_id === window.userId || window.userRole === "admin");

    if (!thread || !thread.id) {
        console.error('Некорректный тред:', thread);
        return;
    }

    const threadList = document.querySelector('.thread-list');
    if (!threadList) {
        console.error('Элемент .thread-list не найден');
        return;
    }

    const noThreadsMessage = threadList.querySelector('.text-center');
    if (noThreadsMessage) {
        noThreadsMessage.remove();
    }
    
    const threadCard = document.createElement('div');
    threadCard.className = 'card mb-3';
    const html = `
        <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
                <h5 class="card-title mb-0">
                    <a href="/threads/${thread.id}" class="text-decoration-none">${thread.title || ''}</a>
                </h5>
                ${(thread.author_id === window.userId || window.userRole === "admin") ? `
                <div class="btn-group">
                    <button type="button" class="btn btn-outline-primary btn-sm" onclick="editThread(${thread.id}, '${thread.title || ''}')">
                        <i class="bi bi-pencil"></i>
                    </button>
                    <button type="button" class="btn btn-outline-danger btn-sm" onclick="deleteThread(${thread.id})">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
                ` : ''}
            </div>
            <p class="card-text text-muted mt-2">
                <small>
                    <i class="bi bi-clock"></i> ${formatDate(thread.created_at)}
                </small>
            </p>
        </div>
    `;
    console.log('HTML для треда:', html);
    threadCard.innerHTML = html;
    threadList.appendChild(threadCard);
}

// Обработчик формы редактирования
document.getElementById('editThreadForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const id = document.getElementById('editThreadId').value;
    const title = document.getElementById('editThreadTitle').value;
    
    fetch(`/api/threads/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: title })
    })
    .then(response => {
        if (response.ok) {
            loadThreads(); // Перезагружаем список тредов после редактирования
            bootstrap.Modal.getInstance(document.getElementById('editThreadModal')).hide();
        } else {
            alert('Ошибка при обновлении треда');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Ошибка при обновлении треда');
    });
});

// Обработчик формы создания треда
document.getElementById('createThreadForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const title = document.getElementById('threadTitle').value;
    console.log('Отправка запроса на создание треда:', title);
    
    fetch('/api/threads', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: title })
    })
    .then(response => {
        console.log('Получен ответ:', response.status);
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Ошибка при создании треда');
        }
    })
    .then(thread => {
        console.log('Получен тред:', thread);
        addThreadToList(thread);
        document.getElementById('threadTitle').value = '';
        bootstrap.Modal.getInstance(document.getElementById('createThreadModal')).hide();
    })
    .catch(error => {
        console.error('Error:', error);
        alert(error.message);
    });
});

// Функция для загрузки тредов с сервера
function loadThreads() {
    console.log('Загрузка тредов...');
    fetch('/api/threads')
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка при загрузке тредов');
            }
            return response.json();
        })
        .then(data => {
            console.log('Получены треды:', data);
            const threadList = document.querySelector('.thread-list');
            if (!threadList) {
                console.error('Элемент .thread-list не найден');
                return;
            }
            
            threadList.innerHTML = '';
            if (!data || data.length === 0) {
                threadList.innerHTML = `
                    <div class="text-center text-muted py-5">
                        <i class="bi bi-chat-square-text display-1"></i>
                        <p class="mt-3">Пока нет тредов. Создайте первый!</p>
                    </div>
                `;
            } else {
                data.forEach(thread => addThreadToList(thread));
            }
        })
        .catch(error => {
            console.error('Ошибка при загрузке тредов:', error);
            const threadList = document.querySelector('.thread-list');
            if (threadList) {
                threadList.innerHTML = `
                    <div class="text-center text-danger py-5">
                        <i class="bi bi-exclamation-triangle display-1"></i>
                        <p class="mt-3">Ошибка при загрузке тредов. Попробуйте обновить страницу.</p>
                    </div>
                `;
            }
        });
}

// Функция для периодической подгрузки тредов
function startThreadsPolling() {
    loadThreads(); // Загружаем треды сразу
    setInterval(loadThreads, 5000); // Обновляем каждые 5 секунд
}

// Запускаем подгрузку тредов при загрузке страницы
document.addEventListener('DOMContentLoaded', startThreadsPolling); 