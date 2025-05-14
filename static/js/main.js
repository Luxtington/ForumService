document.addEventListener('DOMContentLoaded', () => {
    // Функция для получения токена
    const getToken = () => {
        // Получаем все куки
        const cookies = document.cookie.split(';');
        for (let cookie of cookies) {
            const [name, value] = cookie.trim().split('=');
            if (name === 'auth_token') {
                return value;
            }
        }
        return null;
    };

    // Функция для проверки авторизации
    const checkAuth = () => {
        const token = getToken();
        if (!token) {
            console.log('Пользователь не авторизован');
            return false;
        }
        return true;
    };

    // Обработка формы создания поста
    document.querySelector('.new-post')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        if (!checkAuth()) {
            alert('Пожалуйста, войдите в систему');
            return;
        }

        const content = e.target.content.value;
        try {
            const response = await fetch('/api/posts', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${getToken()}`
                },
                body: JSON.stringify({
                    thread_id: window.location.pathname.split('/')[2],
                    content
                })
            });

            if (response.ok) {
                window.location.reload();
            }
        } catch (err) {
            console.error(err);
        }
    });

    // Обработка форм комментариев
    document.querySelectorAll('.comment-form').forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            if (!checkAuth()) {
                alert('Пожалуйста, войдите в систему');
                return;
            }

            const textarea = e.target.querySelector('textarea[name="content"]');
            const postIdInput = e.target.querySelector('input[name="post_id"]');
            
            if (!postIdInput || !textarea) {
                console.error('Не найдены необходимые поля формы');
                return;
            }

            const postId = postIdInput.value;
            const content = textarea.value;

            console.log('Отправка комментария:', {
                post_id: postId,
                content: content,
                post_id_type: typeof postId
            });

            if (!postId) {
                alert('Ошибка: ID поста не найден');
                return;
            }

            try {
                const response = await fetch('/api/comments', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${getToken()}`
                    },
                    body: JSON.stringify({
                        post_id: String(postId),
                        content: content
                    })
                });

                const responseData = await response.json();
                console.log('Ответ сервера:', responseData);

                if (!response.ok) {
                    throw new Error(responseData.error || 'Ошибка при отправке комментария');
                }

                const comment = responseData;
                const commentsContainer = document.getElementById('comments-container');
                
                // Удаляем сообщение "нет комментариев", если оно есть
                const noComments = commentsContainer.querySelector('.text-center.text-muted');
                if (noComments) {
                    noComments.remove();
                }

                // Создаем новый элемент комментария
                const commentElement = document.createElement('div');
                commentElement.className = 'card mb-3 fade-in';
                commentElement.innerHTML = `
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-center mb-2">
                            <div class="d-flex align-items-center">
                                <i class="bi bi-person-circle me-2"></i>
                                <strong>${comment.author_name}</strong>
                            </div>
                            <small class="text-muted">
                                <i class="bi bi-clock"></i> ${new Date(comment.created_at).toLocaleString()}
                            </small>
                        </div>
                        <p class="card-text mb-0">${comment.content}</p>
                        <div class="mt-2">
                            <button class="btn btn-outline-danger btn-sm delete-comment" data-comment-id="${comment.id}" onclick="return confirm('Вы уверены, что хотите удалить этот комментарий?')">
                                <i class="bi bi-trash"></i> Удалить
                            </button>
                        </div>
                    </div>
                `;

                // Добавляем комментарий в начало списка
                commentsContainer.insertBefore(commentElement, commentsContainer.firstChild);
                
                // Очищаем форму и закрываем модальное окно
                textarea.value = '';
                const modal = bootstrap.Modal.getInstance(document.getElementById('commentModal'));
                if (modal) {
                    modal.hide();
                }
            } catch (err) {
                console.error('Ошибка при отправке комментария:', err);
                alert(err.message);
            }
        });
    });

    // Обработка удаления комментариев
    document.querySelectorAll('.delete-comment').forEach(button => {
        button.addEventListener('click', async (e) => {
            e.preventDefault();
            if (!checkAuth()) {
                alert('Пожалуйста, войдите в систему');
                return;
            }

            if (!confirm('Вы уверены, что хотите удалить этот комментарий?')) {
                return;
            }

            const commentId = button.dataset.commentId;
            try {
                const response = await fetch(`/api/comments/${commentId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${getToken()}`
                    }
                });

                if (response.ok) {
                    // Удаляем элемент комментария из DOM
                    const commentElement = button.closest('.card');
                    commentElement.remove();

                    // Если больше нет комментариев, показываем сообщение
                    const commentsContainer = document.getElementById('comments-container');
                    if (!commentsContainer.querySelector('.card')) {
                        commentsContainer.innerHTML = `
                            <div class="text-center text-muted">
                                <i class="bi bi-chat-square-text display-4"></i>
                                <p class="mt-3">Пока нет комментариев. Будьте первым!</p>
                            </div>
                        `;
                    }
                }
            } catch (err) {
                console.error(err);
                alert('Ошибка при удалении комментария');
            }
        });
    });
});