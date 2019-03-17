<template>
  <div class="comment">
    <div class="header">
      <p class="date">{{ date }}</p>
    </div>
    <p class="body">{{ comment.content }}</p>
  </div>
</template>

<script>
export default {
    name: 'comment',
    props: {
        comment: Object,
    },
    data() {
        return {
            date: '',
        };
    },
    created() {
        try {
            this.date = new Date(this.comment.time).toLocaleString();
        } catch (e) {
            console.log(e);
        }
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

.body {
  white-space: pre-line;
  padding: 0;
  margin: 0;
}
</style>
