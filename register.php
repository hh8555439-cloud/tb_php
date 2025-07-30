<?php
header('Content-Type: application/json');
session_start();
require_once 'db_connection.php';

$response = ['success' => false, 'message' => '', 'errors' => []];

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
  $username = trim($_POST['username'] ?? '');
  $email = trim($_POST['email'] ?? '');
  $password = $_POST['password'] ?? '';
  $confirm_password = $_POST['confirm_password'] ?? '';

  // 输入验证
  if (empty($username))
    $response['errors']['username'] = '用户名不能为空';
  if (empty($email) || !filter_var($email, FILTER_VALIDATE_EMAIL))
    $response['errors']['email'] = '邮箱格式无效';
  if (empty($password))
    $response['errors']['password'] = '密码不能为空';
  if ($password !== $confirm_password)
    $response['errors']['confirm_password'] = '两次密码不一致';

  if (empty($response['errors'])) {
    try {
      $db = Database::getInstance();

      // 检查用户名和邮箱唯一性
      $stmt = $db->prepare("SELECT id FROM users WHERE username = :username OR email = :email");
      $stmt->execute([':username' => $username, ':email' => $email]);

      if ($stmt->rowCount() > 0) {
        $response['errors']['global'] = '用户名或邮箱已存在';
      } else {
        //$password_hash = password_hash($password, PASSWORD_DEFAULT);
        $password_hash = $password;
        $stmt = $db->prepare("INSERT INTO users (username, email, password) VALUES (:username, :email, :password)");
        $stmt->execute([':username' => $username, ':email' => $email, ':password' => $password_hash]);

        $response['success'] = true;
        $response['message'] = '注册成功';
        //$_SESSION['user_id'] = $db->lastInsertId();
      }
    } catch (PDOException $e) {
      error_log('Database error: ' . $e->getMessage());
      $response['message'] = '系统错误，请稍后再试';
    }
  }
}

echo json_encode($response);
?>