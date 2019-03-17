<template>
  <div class="comments">
    <CreateComment :postId="postId" />
    <div class="content">
      <ul>
        <li v-for="comment in comments" :key="comment.commentId">
          <Comment :comment="comment" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import CreateComment from '@/components/CreateComment.vue';
import Comment from '@/components/Comment.vue';

export default {
    name: 'comments',
    props: {
        postId: String,
        comments: Array
    },
    components: {
        CreateComment,
        Comment,
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.date {
  padding: 0;
  margin: 0 $margin1 0 auto;
  font-size: $fontsize1;
}

.header {
  display: flex;
  margin-bottom: $margin1;
}

.content {
  white-space: pre-line;
  padding: 0;
  margin: 0;
}
</style>
