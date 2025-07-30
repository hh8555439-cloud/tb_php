<?php
// 开启会话
session_start();
session_unset(); // 清除所有 session 变量
session_destroy(); // 销毁 session
session_start(); // 重新开启 session

// 如果用户已登录，重定向到主页
if (isset($_SESSION['user_id'])) {
  header("Location: index.html");
  exit();
}

// 引入数据库连接
require_once 'config.php';

// 初始化错误信息
$error = '';
// 处理登录表单提交
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
  // 获取并清理输入数据
  $username = trim($_POST['username'] ?? '');
  $password = $_POST['password'] ?? '';

  // 验证输入
  if (empty($username) || empty($password)) {
    $error = '用户名和密码不能为空';
  } else {
    try {
      // 使用 config.php 中的 $db
      $stmt = $db->prepare("SELECT id, username, password, role FROM users WHERE username = :username");
      $stmt->bindParam(':username', $username, PDO::PARAM_STR);
      $stmt->execute();

      // 检查用户是否存在
      if ($stmt->rowCount() === 1) {
        $user = $stmt->fetch(PDO::FETCH_ASSOC);

        // 如果数据库存的是 hash，用 password_verify
        // if (password_verify($password, $user['password'])) {
        // 如果存的是明文密码，用下面这行
        if ($password === $user['password']) {
          // 设置会话变量
          $_SESSION['user_id'] = $user['id'];
          $_SESSION['username'] = $user['username'];
          $_SESSION['role'] = $user['role'];

          // 重定向到主页
          header("Location: index.html");
          exit();
        } else {
          $error = '用户名或密码错误';
        }
      } else {
        $error = '用户名或密码错误';
      }
    } catch (PDOException $e) {
      error_log('数据库错误: ' . $e->getMessage());
      $error = '系统错误，请稍后再试';
    }
  }
}
?>
<!DOCTYPE html>
<html lang="zh-CN">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>用户登录</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f5f5f5;
      margin: 0;
      padding: 0;
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
    }

    .login-container {
      width: 400px;
      background-color: #fff;
      border-radius: 8px;
      box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
      padding: 40px;
    }

    .login-header {
      text-align: center;
      margin-bottom: 30px;
    }

    .login-header h1 {
      color: #333;
      font-size: 24px;
      margin-bottom: 10px;
    }

    .login-form .form-group {
      margin-bottom: 20px;
    }

    .login-form label {
      display: block;
      margin-bottom: 8px;
      color: #555;
      font-weight: bold;
    }

    .login-form input {
      width: 100%;
      padding: 12px;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 16px;
      box-sizing: border-box;
    }

    .login-form input:focus {
      border-color: #007BFF;
      outline: none;
    }

    .login-form button {
      width: 100%;
      padding: 12px;
      background-color: #007BFF;
      color: #fff;
      border: none;
      border-radius: 4px;
      font-size: 16px;
      cursor: pointer;
      transition: background-color 0.3s;
    }

    .login-form button:hover {
      background-color: #0056b3;
    }

    .error-message {
      color: #dc3545;
      margin-bottom: 20px;
      text-align: center;
    }

    .register-link {
      text-align: center;
      margin-top: 20px;
    }

    .register-link a {
      color: #007BFF;
      text-decoration: none;
    }

    .register-link a:hover {
      text-decoration: underline;
    }
  </style>
</head>

<body>
  <div class="login-container">
    <div class="login-header">
      <h1>用户登录</h1>
    </div>

    <?php if (!empty($error)): ?>
      <div class="error-message"><?php echo htmlspecialchars($error); ?></div>
    <?php endif; ?>

    <form class="login-form" action="login.php" method="POST">
      <div class="form-group">
        <label for="username">用户名</label>
        <input type="text" id="username" name="username" required>
      </div>

      <div class="form-group">
        <label for="password">密码</label>
        <input type="password" id="password" name="password" required>
      </div>

      <div class="form-group">
        <button type="submit">登录</button>
      </div>
    </form>

    <div class="register-link">
      还没有账号？<a href="register.html">立即注册</a>
    </div>
  </div>
</body>

</html>