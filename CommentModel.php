<?php
class CommentModel
{
  private $db;

  public function __construct($db)
  {
    $this->db = $db;
  }

  // 添加评论
  public function addComment($data)
  {
    $stmt = $this->db->prepare("INSERT INTO comments (content, user_id, goods_id, to_user_id, root_id, to_answer_id, type) 
                                   VALUES (:content, :user_id, :goods_id, :to_user_id, :root_id, :to_answer_id, :type)");
    return $stmt->execute($data);
  }

  // 获取商品所有评论
  public function getCommentsByGoods($goodsId)
  {
    $stmt = $this->db->prepare("SELECT * FROM comments WHERE goods_id = :goods_id ORDER BY create_time ASC");
    $stmt->execute([':goods_id' => $goodsId]);
    return $stmt->fetchAll(PDO::FETCH_ASSOC);
  }

  // 获取用户信息（模拟，实际应从用户服务获取）
  public function getUserInfo($userId)
  {
    $stmt = $this->db->prepare("SELECT id, username FROM users WHERE id = ?");
    $stmt->execute([$userId]);
    $row = $stmt->fetch(PDO::FETCH_ASSOC);
    if ($row) {
      return [
        'id' => $row['id'],
        'name' => $row['username'], // 这里用 name 作为前端显示
      ];
    }
    return [
      'id' => $userId,
      'name' => '未知用户'
    ];
  }
}
?>