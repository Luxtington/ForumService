{{define "threads.js"}}
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
                window.location.reload();
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
            window.location.reload();
        } else {
            alert('Ошибка при обновлении треда');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Ошибка при обновлении треда');
    });
});

// Обработчик формы создания
document.getElementById('createThreadForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const title = document.getElementById('threadTitle').value;
    
    fetch('/api/threads', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: title })
    })
    .then(response => {
        if (response.ok) {
            window.location.reload();
        } else {
            alert('Ошибка при создании треда');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Ошибка при создании треда');
    });
});
{{end}} 