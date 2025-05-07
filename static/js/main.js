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
                    window.location.reload();
                }
            } catch (err) {
                console.error(err);
            }
        });
    });
});