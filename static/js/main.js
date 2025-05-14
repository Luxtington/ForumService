document.addEventListener('DOMContentLoaded', () => {
    // Обработка формы создания поста
    document.querySelector('.new-post')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const content = e.target.content.value;

        try {
            const response = await fetch('/api/posts', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
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
            const textarea = e.target.querySelector('textarea');
            const postId = e.target.dataset.postId;

            try {
                const response = await fetch('/api/comments', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    },
                    body: JSON.stringify({
                        post_id: postId,
                        content: textarea.value
                    })
                });

                if (response.ok) {
                    const comment = await response.json();
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
                }
            } catch (err) {
                console.error(err);
                alert('Ошибка при отправке комментария');
            }
        });
    });

    // Обработка удаления комментариев
    document.querySelectorAll('.delete-comment').forEach(button => {
        button.addEventListener('click', async (e) => {
            e.preventDefault();
            if (!confirm('Вы уверены, что хотите удалить этот комментарий?')) {
                return;
            }

            const commentId = button.dataset.commentId;
            try {
                const response = await fetch(`/api/comments/${commentId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
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