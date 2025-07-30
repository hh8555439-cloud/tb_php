<?php
require_once 'CommentModel.php';

class CommentController
{
  private $model;

  public function __construct($db)
  {
    $this->model = new CommentModel($db);
  }

  // 获取商品评论
  public function getGoodsComments($goodsId)
  {
    $comments = $this->model->getCommentsByGoods($goodsId);
    $rootComments = [];
    $answerMap = [];
    $userCache = [];

    // 先处理所有评论
    foreach ($comments as $comment) {
      // 缓存用户信息
      if (!isset($userCache[$comment['user_id']])) {
        $userCache[$comment['user_id']] = $this->model->getUserInfo($comment['user_id']);
      }

      if ($comment['type'] == 'root') {
        // 根评论
        $rootComment = [
          'id' => $comment['id'],
          'content' => $comment['content'],
          'user' => $userCache[$comment['user_id']],
          'create_time' => $comment['create_time'],
          'answers' => [],
          'goods_id' => $comment['goods_id']
        ];

        $rootComments[$comment['id']] = $rootComment;
        $answerMap[$comment['id']] = &$rootComments[$comment['id']]['answers'];
      } else {
        // 子评论
        if (!isset($answerMap[$comment['root_id']])) {
          continue; // 根评论不存在，跳过
        }

        $toUser = null;
        if ($comment['to_user_id']) {
          if (!isset($userCache[$comment['to_user_id']])) {
            $userCache[$comment['to_user_id']] = $this->model->getUserInfo($comment['to_user_id']);
          }
          $toUser = $userCache[$comment['to_user_id']];
        }

        $answer = [
          'id' => $comment['id'],
          'content' => $comment['content'],
          'user' => $userCache[$comment['user_id']],
          'to_user' => $toUser,
          'create_time' => $comment['create_time'],
          'root_id' => $comment['root_id'],
          'to_answer_id' => $comment['to_answer_id'], // 新增
          'goods_id' => $comment['goods_id'] // 新增
        ];

        $answerMap[$comment['root_id']][] = $answer;
      }
    }

    // 按创建时间倒序排列根评论
    usort($rootComments, function ($a, $b) {
      return strtotime($b['create_time']) - strtotime($a['create_time']);
    });

    return $rootComments;
  }

  // 添加评论
  public function addComment($data)
  {
    $defaults = [
      'to_user_id' => null,
      'root_id' => null,
      'to_answer_id' => null,
      'type' => 'root'
    ];

    $data = array_merge($defaults, $data);

    if ($data['type'] == 'answer' && empty($data['root_id'])) {
      throw new Exception('回复评论必须指定root_id');
    }

    return $this->model->addComment($data);
  }
}
?>