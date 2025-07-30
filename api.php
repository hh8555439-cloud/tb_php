<?php
require_once 'config.php';
require_once 'CommentController.php';

header('Content-Type: application/json');

$controller = new CommentController($db);

try {
  $action = $_GET['action'] ?? '';
  $goodsId = $_GET['goods_id'] ?? 0;

  switch ($action) {
    case 'get_comments':
      $comments = $controller->getGoodsComments($goodsId);
      echo json_encode(['code' => 0, 'data' => $comments]);
      break;

    case 'add_comment':
      $data = [
        'content' => $_POST['content'],
        'user_id' => $_POST['user_id'],
        'goods_id' => $_POST['goods_id'],
        'to_user_id' => $_POST['to_user_id'] !== '' ? $_POST['to_user_id'] : null,
        'root_id' => $_POST['root_id'] ?? null,
        'to_answer_id' => $_POST['to_answer_id'] ?? null,
        'type' => $_POST['type'] ?? 'root'
      ];

      $result = $controller->addComment($data);
      echo json_encode(['code' => 0, 'data' => $result]);
      break;

    case 'get_user':
      session_start();
      if (isset($_SESSION['user_id'])) {
        echo json_encode([
          'code' => 0,
          'data' => [
            'id' => $_SESSION['user_id'],
            'username' => $_SESSION['username'],
            'role' => $_SESSION['role']
          ]
        ]);
      } else {
        echo json_encode(['code' => 1, 'data' => null]);
      }
      break;

    case 'logout':
      session_destroy();
      echo json_encode(['code' => 0]);
      break;

    case 'get_messages':
      // 查询所有留言
      $stmt = $db->query("SELECT m.*, u.username FROM messages m LEFT JOIN users u ON m.user_id = u.id ORDER BY m.created_at DESC");
      $messages = [];
      while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
        $messages[] = [
          'id' => $row['id'],
          'user' => [
            'id' => $row['user_id'],
            'name' => $row['username']
          ],
          'content' => $row['content'],
          'created_at' => $row['created_at']
        ];
      }
      echo json_encode(['code' => 0, 'data' => $messages]);
      break;

    case 'add_message':
      $userId = $_POST['user_id'] ?? 0;
      $content = $_POST['content'] ?? '';
      if (!$userId || !$content) {
        echo json_encode(['code' => 1, 'message' => '参数错误']);
        exit;
      }
      $stmt = $db->prepare("INSERT INTO messages (user_id, content, created_at) VALUES (?, ?, NOW())");
      $result = $stmt->execute([$userId, $content]);
      if ($result) {
        echo json_encode(['code' => 0, 'message' => '留言成功']);
      } else {
        echo json_encode(['code' => 1, 'message' => '留言失败']);
      }
      break;

    case 'delete_message':
      session_start();
      if ($_SESSION['role'] !== 'admin') {
        echo json_encode(['code' => 1, 'message' => '无权限']);
        exit;
      }
      $id = $_POST['id'] ?? 0;
      $stmt = $db->prepare("DELETE FROM messages WHERE id = ?");
      $result = $stmt->execute([$id]);
      echo json_encode(['code' => $result ? 0 : 1, 'message' => $result ? '删除成功' : '删除失败']);
      break;

    case 'delete_comment':
      session_start();
      if ($_SESSION['role'] !== 'admin') {
        echo json_encode(['code' => 1, 'message' => '无权限']);
        exit;
      }
      $id = $_POST['id'] ?? 0;
      $stmt = $db->prepare("DELETE FROM comments WHERE id = ?");
      $result = $stmt->execute([$id]);
      echo json_encode(['code' => $result ? 0 : 1, 'message' => $result ? '删除成功' : '删除失败']);
      break;

    default:
      echo json_encode(['code' => 1, 'message' => 'Invalid action']);
  }
} catch (Exception $e) {
  echo json_encode(['code' => 1, 'message' => $e->getMessage()]);
}
?>