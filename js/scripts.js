// ================= 敏感词过滤模块 =================

// 创建实例
const search = new WordsSearch();
// 从文件加载敏感词库
async function loadSensitiveWords() {
  try {
    const response = await fetch('./ToolGoodWords/sensitiveWords.txt');
    if (!response.ok) throw new Error('词库加载失败');
    const text = await response.text();
    return text;

  } catch (error) {
    console.error('加载敏感词库失败:', error);
    // 默认词库
    return [
      '脏话1', '脏话2', '政治敏感词', '色情词'
    ];
  }
}

// 使用实际存在的WordsSearch类
let wordSearch;
(async function initScanner() {
  try {
    const words = await loadSensitiveWords();
    wordSearch = new WordsSearch(); // 需确认导出方式
    wordSearch.SetKeywords(words.split('|'));
  } catch (e) {
    console.error('模块加载失败:', e);
  }
})();

// 修正版敏感词过滤
function filterSensitiveWords(content) {
  if (!content || !wordSearch) return {
    filtered: content,
    hits: [],
    safe: true
  };

  // 检测简体敏感词
  const result = wordSearch.FindFirst(content);
  return {
    filtered: wordSearch.Replace(content, '*'),
    hits: result,
    safe: result == null
  };
}

function getUserInfo() {
  fetch('api.php?action=get_user')
    .then(res => res.json())
    .then(data => {
      const userInfoDiv = document.getElementById('user-info');
      if (data.code === 0 && data.data) {
        userInfoDiv.innerHTML = `当前用户：${data.data.username} <button onclick="logout()">退出</button>`;
        document.getElementById('user-id').value = data.data.id;
        window.currentUserRole = data.data.role; // 这里设置角色
      } else {
        userInfoDiv.innerHTML = `
  <a href="login.php" id="login-btn" style="
    display: inline-block;
    padding: 6px 18px;
    background-color: #1890ff;
    color: #fff;
    border-radius: 4px;
    text-decoration: none;
    font-weight: bold;
    box-shadow: 0 2px 8px rgba(24,144,255,0.08);
    transition: background 0.2s;
    margin-left: 10px;
  ">请登录</a>
`;
        document.getElementById('user-id').value = '';
      }
    });
}

function logout() {
  window.location.href = 'logout.php';
}
getUserInfo();

// 获取所有留言
function loadMessages() {
  fetch('api.php?action=get_messages')
    .then(res => res.json())
    .then(data => {
      if (data.code === 0) {
        renderMessages(data.data);
      } else {
        console.error('获取留言失败:', data.message);
      }
    })
    .catch(error => console.error('Error:', error));
}

// 渲染留言和评论
function renderMessages(messages) {
  const container = document.getElementById('comments-container');
  container.innerHTML = '';

  const showCount = 3;
  const showMessages = messages.slice(0, showCount);

  showMessages.forEach(message => {
    let deleteBtn = '';
    if (window.currentUserRole === 'admin') {
      deleteBtn = `<button class="delete-btn" onclick="deleteMessage(${message.id})">删除</button>`;
    }
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message-block';
    messageDiv.innerHTML = `
  <div class="message-header">
    <span class="user-name">${message.user.name}</span>
    <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(message.created_at).toLocaleString()}</span>
    ${deleteBtn}
  </div>
  <div class="message-content">${escapeHtml(message.content)}</div>
  <div class="message-actions" style="display:flex;align-items:center;justify-content:space-between;margin-top:12px;">
    <span class="reply-btn" onclick="prepareReply(0, ${message.user.id}, -1, ${message.id}, '${message.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')" style="margin-left:0;">回复</span>
    <button class="toggle-comments-btn" data-id="${message.id}">展开评论</button>
  </div>
  <div class="comments-list" id="comments-list-${message.id}" style="display:none;"></div>
`;
    container.appendChild(messageDiv);
  });

  // 如果有更多留言，显示“查看更多留言”按钮
  if (messages.length > showCount) {
    const moreBtn = document.createElement('button');
    moreBtn.textContent = '查看更多留言';
    moreBtn.className = 'toggle-comments-btn';
    moreBtn.style.marginTop = '16px';
    moreBtn.onclick = function () {
      container.innerHTML = '';
      messages.forEach(message => {
        let deleteBtn = '';
        if (window.currentUserRole === 'admin') {
          deleteBtn = `<button class="delete-btn" onclick="deleteMessage(${message.id})">删除</button>`;
        }
        const messageDiv = document.createElement('div');
        messageDiv.className = 'message-block';
        messageDiv.innerHTML = `
      <div class="message-header">
        <span class="user-name">${message.user.name}</span>
        <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(message.created_at).toLocaleString()}</span>
        ${deleteBtn}
      </div>
      <div class="message-content">${escapeHtml(message.content)}</div>
      <div class="message-actions" style="display:flex;align-items:center;justify-content:space-between;margin-top:12px;">
        <span class="reply-btn" onclick="prepareReply(0, ${message.user.id}, -1, ${message.id}, '${message.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')" style="margin-left:0;">回复</span>
        <button class="toggle-comments-btn" data-id="${message.id}">展开评论</button>
      </div>
      <div class="comments-list" id="comments-list-${message.id}" style="display:none;"></div>
    `;
        container.appendChild(messageDiv);
      });
      bindToggleComments();
      this.remove();
    };
    container.appendChild(moreBtn);
  }

  bindToggleComments();
}

function bindToggleComments() {
  document.querySelectorAll('.toggle-comments-btn[data-id]').forEach(btn => {
    btn.addEventListener('click', function () {
      const messageId = this.getAttribute('data-id');
      const commentListDiv = document.getElementById('comments-list-' + messageId);
      if (commentListDiv.style.display === 'none') {
        loadComments(messageId, commentListDiv);
        commentListDiv.style.display = 'block';
        this.textContent = '收起评论';
      } else {
        commentListDiv.style.display = 'none';
        this.textContent = '展开评论';
      }
    });
  });
}

function renderComments(comments, container) {
  container.innerHTML = '';
  const showCount = 3;
  const showComments = comments.slice(0, showCount);

  showComments.forEach(comment => {
    let deleteBtn = '';
    if (window.currentUserRole === 'admin') {
      deleteBtn = `<button class="delete-btn" onclick="deleteComment(${comment.id}, ${comment.goods_id})">删除</button>`;
    }
    const commentElement = document.createElement('div');
    commentElement.className = 'comment';
    commentElement.innerHTML = `
  <div class="comment-header">
    <span class="user-name">${comment.user.name}</span>
    <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(comment.create_time).toLocaleString()}</span>
    ${deleteBtn}
  </div>
  <div class="comment-content">${escapeHtml(comment.content)}</div>
  <div class="comment-footer">
    <span class="reply-btn" onclick="prepareReply(${comment.id}, ${comment.user.id}, ${comment.id}, ${comment.goods_id}, '${comment.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')">回复</span>
  </div>
`;

    // 二级评论
    if (comment.answers && comment.answers.length > 0) {
      const answersContainer = document.createElement('div');
      answersContainer.className = 'answers-container';
      const showAnswers = comment.answers.slice(0, showCount);
      showAnswers.forEach(answer => {
        let replyText = '';
        if (answer.to_user) {
          replyText = ` 回复 @${answer.to_user.name}`;
        }
        let deleteBtn2 = '';
        if (window.currentUserRole === 'admin') {
          deleteBtn2 = `<button class="delete-btn" onclick="deleteComment(${answer.id}, ${comment.goods_id})">删除</button>`;
        }
        const answerElement = document.createElement('div');
        answerElement.className = 'answer';
        answerElement.innerHTML = `
      <div class="comment-header">
        <span class="user-name">${answer.user.name}</span>
        <span>${replyText}</span>
        <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(answer.create_time).toLocaleString()}</span>
        ${deleteBtn2}
      </div>
      <div class="comment-content">${escapeHtml(answer.content)}</div>
      <div class="comment-footer">
        <span class="reply-btn" onclick="prepareReply(${comment.id}, ${answer.user.id}, ${answer.id}, ${comment.goods_id}, '${answer.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')">回复</span>
      </div>
    `;
        answersContainer.appendChild(answerElement);
      });
      // 在 renderComments 的 showAnswers 部分后面加：
      if (comment.answers && comment.answers.length > showCount) {
        const moreAnswersBtn = document.createElement('button');
        moreAnswersBtn.textContent = '查看更多回复';
        moreAnswersBtn.className = 'toggle-comments-btn';
        moreAnswersBtn.onclick = function () {
          answersContainer.innerHTML = '';
          comment.answers.forEach(answer => {
            // ...渲染 answer 的代码同上...
            let replyText = '';
            if (answer.to_user) {
              replyText = ` 回复 @${answer.to_user.name}`;
            }
            let deleteBtn2 = '';
            if (window.currentUserRole === 'admin') {
              deleteBtn2 = `<button class="delete-btn" onclick="deleteComment(${answer.id}, ${comment.goods_id})">删除</button>`;
            }
            const answerElement = document.createElement('div');
            answerElement.className = 'answer';
            answerElement.innerHTML = `
  <div class="comment-header">
    <span class="user-name">${answer.user.name}</span>
    <span>${replyText}</span>
    <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(answer.create_time).toLocaleString()}</span>
    ${deleteBtn2}
  </div>
  <div class="comment-content">${escapeHtml(answer.content)}</div>
  <div class="comment-footer">
    <span class="reply-btn" onclick="prepareReply(${comment.id}, ${answer.user.id}, ${answer.id}, ${comment.goods_id}, '${answer.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')">回复</span>
  </div>
`;
            answersContainer.appendChild(answerElement);
          });
          this.remove();
        };
        answersContainer.appendChild(moreAnswersBtn);
      }

      commentElement.appendChild(answersContainer);
    }
    container.appendChild(commentElement);
  });

  // 如果有更多一级评论，显示“查看更多评论”按钮
  if (comments.length > showCount) {
    const moreBtn = document.createElement('button');
    moreBtn.textContent = '查看更多评论';
    moreBtn.className = 'toggle-comments-btn';
    moreBtn.onclick = function () {
      container.innerHTML = '';
      comments.forEach(comment => {
        let deleteBtn = '';
        if (window.currentUserRole === 'admin') {
          deleteBtn = `<button class="delete-btn" onclick="deleteComment(${comment.id}, ${comment.goods_id})">删除</button>`;
        }
        const commentElement = document.createElement('div');
        commentElement.className = 'comment';
        commentElement.innerHTML = `
      <div class="comment-header">
        <span class="user-name">${comment.user.name}</span>
        <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(comment.create_time).toLocaleString()}</span>
        ${deleteBtn}
      </div>
      <div class="comment-content">${escapeHtml(comment.content)}</div>
      <div class="comment-footer">
        <span class="reply-btn" onclick="prepareReply(${comment.id}, ${comment.user.id}, ${comment.id}, ${comment.goods_id}, '${comment.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')">回复</span>
      </div>
    `;
        // 二级评论
        if (comment.answers && comment.answers.length > 0) {
          const answersContainer = document.createElement('div');
          answersContainer.className = 'answers-container';
          comment.answers.forEach(answer => {
            let replyText = '';
            if (answer.to_user) {
              replyText = ` 回复 @${answer.to_user.name}`;
            }
            let deleteBtn2 = '';
            if (window.currentUserRole === 'admin') {
              deleteBtn2 = `<button class="delete-btn" onclick="deleteComment(${answer.id}, ${comment.goods_id})">删除</button>`;
            }
            const answerElement = document.createElement('div');
            answerElement.className = 'answer';
            answerElement.innerHTML = `
          <div class="comment-header">
            <span class="user-name">${answer.user.name}</span>
            <span>${replyText}</span>
            <span style="color:#999;font-size:12px;margin-left:10px;">${new Date(answer.create_time).toLocaleString()}</span>
            ${deleteBtn2}
          </div>
          <div class="comment-content">${escapeHtml(answer.content)}</div>
          <div class="comment-footer">
            <span class="reply-btn" onclick="prepareReply(${comment.id}, ${answer.user.id}, ${answer.id}, ${comment.goods_id}, '${answer.user.name.replace(/'/g, "\\'").replace(/"/g, '"')}')">回复</span>
          </div>
        `;
            answersContainer.appendChild(answerElement);
          });
          commentElement.appendChild(answersContainer);
        }
        container.appendChild(commentElement);
      });
      this.remove();
    };
    container.appendChild(moreBtn);
  }
}

// 获取某条留言下的评论
function loadComments(messageId, container) {
  fetch(`api.php?action=get_comments&goods_id=${messageId}`)
    .then(res => res.json())
    .then(data => {
      if (data.code === 0) {
        renderComments(data.data, container);
      } else {
        container.innerHTML = '<div style="color:red;">评论加载失败</div>';
      }
    })
    .catch(error => {
      container.innerHTML = '<div style="color:red;">评论加载失败</div>';
    });
}

// 页面初始化加载所有留言
loadMessages();

// 删除留言
function deleteMessage(id) {
  if (!confirm('确定要删除该留言吗？')) return;
  fetch('api.php?action=delete_message', {
    method: 'POST',
    body: new URLSearchParams({ id })
  })
    .then(res => res.json())
    .then(data => {
      if (data.code === 0) {
        loadMessages();
      } else {
        alert(data.message);
      }
    });
}

// 删除评论
function deleteComment(id, goodsId) {
  if (!confirm('确定要删除该评论吗？')) return;
  fetch('api.php?action=delete_comment', {
    method: 'POST',
    body: new URLSearchParams({ id })
  })
    .then(res => res.json())
    .then(data => {
      if (data.code === 0) {
        const commentListDiv = document.getElementById('comments-list-' + goodsId);
        loadComments(goodsId, commentListDiv); // 只刷新该留言下的评论
      } else {
        alert(data.message);
      }
    });
}

// 修改留言提交逻辑
document.getElementById('submit-comment').onclick = async function () {
  const rawContent = document.getElementById('comment-content').value;

  // 空内容检查
  if (!rawContent || rawContent.trim() === '') {
    alert('留言内容不能为空');
    return;
  }

  // 获取敏感词检测结果
  const { safe, filtered, hits } = filterSensitiveWords(rawContent);

  // 显示敏感词警告
  if (safe == false) {
    const confirmMsg = `内容包含敏感词: ${hits.Keyword}...\n\n无法提交！`;
    alert(confirmMsg);
    return;
  }

  // 获取用户ID
  const userId = document.getElementById('user-id').value;
  if (!userId) {
    alert('请先登录后再留言');
    window.location.href = 'login.html';
    return;
  }

  // 准备提交数据（使用过滤后的内容）
  const formData = new FormData();
  formData.append('user_id', userId);
  formData.append('content', safe ? rawContent : filtered);

  // 执行提交
  try {
    const response = await fetch('api.php?action=add_message', {
      method: 'POST',
      body: formData
    });
    const data = await response.json();

    if (data.code === 0) {
      document.getElementById('comment-content').value = '';
      loadMessages();
      alert('留言成功！');
    } else {
      alert('留言失败: ' + data.message);
    }
  } catch (error) {
    alert('留言失败: ' + error);
  }
};


let replyContext = {}; // 用于保存当前回复的上下文

// 将函数声明移至文件顶部
window.prepareReply = function (rootId, toUserId, answerId = -1, goodsId = '', toUserName = '') {
  const userId = document.getElementById('user-id').value;
  if (!userId) {
    alert('请先登录后再回复');
    window.location.href = 'login.html';
    return;
  }
  replyContext = { rootId, toUserId, answerId, goodsId };
  document.getElementById('reply-content').value = '';
  document.getElementById('reply-title').textContent = toUserName ? `回复 @${toUserName}` : '回复';
  document.getElementById('reply-modal').style.display = 'flex';
};

// 原函数定义位置改为引用
document.querySelectorAll('.reply-btn').forEach(btn => {
  btn.onclick = function () {
    const args = JSON.parse(this.dataset.args);
    window.prepareReply(...args);
  };
});

// 取消按钮
document.getElementById('reply-cancel').onclick = function () {
  document.getElementById('reply-modal').style.display = 'none';
};

// 提交回复
document.getElementById('reply-submit').addEventListener('click', function () {
  const rawContent = document.getElementById('reply-content').value;

  // 空内容检查
  if (!rawContent || rawContent.trim() === '') {
    alert('回复内容不能为空');
    return;
  }

  // 获取敏感词检测结果
  let { safe, filtered, hits } = filterSensitiveWords(rawContent);

  // 显示敏感词警告
  if (safe == false) {
    const confirmMsg = `内容包含敏感词: ${hits.Keyword}\n\n是否发送过滤后的内容？\n（原内容将被替换为：${filtered}）`;
    if (!confirm(confirmMsg)) return;

    // 用户确认后使用过滤内容
    document.getElementById('reply-content').value = filtered;
    safe = true; // 标记为安全以继续提交
  }

  // 前端过滤
  const filteredContent = filtered;
  const userId = document.getElementById('user-id').value;

  const formData = new FormData();
  formData.append('content', filteredContent);
  formData.append('user_id', userId);
  formData.append('goods_id', replyContext.goodsId);
  formData.append('to_user_id', replyContext.toUserId);
  formData.append('root_id', replyContext.rootId);
  formData.append('to_answer_id', replyContext.answerId);
  formData.append('type', 'answer');
  if (replyContext.rootId == 0) {
    formData.append('type', 'root');
  }

  fetch('api.php?action=add_comment', {
    method: 'POST',
    body: formData
  })
    .then(response => response.json())
    .then(data => {
      if (data.code === 0) {
        document.getElementById('reply-modal').style.display = 'none';
        // 重新加载评论
        const commentListDiv = document.getElementById('comments-list-' + replyContext.goodsId);
        loadComments(replyContext.goodsId, commentListDiv);
      } else {
        alert('回复失败: ' + data.message);
      }
    })
    .catch(error => alert('回复失败: ' + error));
});

// 在文件末尾添加初始化调用
// 页面加载时自动执行
window.addEventListener('DOMContentLoaded', function () {
  // 获取用户信息
  getUserInfo();

  // 加载留言
  loadMessages();
});

// 确保这些函数在文件末尾被调用
getUserInfo();
loadMessages();