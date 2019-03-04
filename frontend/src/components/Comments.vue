<template>
  <div class="comments">
    <CreateComment :postId="postId" />
  </div>
</template>

<script>
import CreateComment from '@/components/CreateComment.vue';

export default {
    name: 'comments',
    props: {
        postId: String,
        comments: Object,
        profileView: { type: Boolean, default: false },
    },
    components: {
        CreateComment,
    },
    data() {
        return {
            date: '',
        };
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
