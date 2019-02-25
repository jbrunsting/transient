<template>
  <div class="post">
    <div class="header">
      <h3 class="title" v-if="post.postUrl" >
        <a :href="post.postUrl">{{ post.title }}</a>
      </h3>
      <h3 class="title" v-else>{{ post.title }}</h3>
      <p class="date">{{ date }}</p>
    </div>
    <p class="body">{{ post.content }}</p>
    <form class="delete" v-if="profileView" @submit.prevent="deletePost">
      <button type="submit">Delete</button>
    </form>
  </div>
</template>

<script>
export default {
    name: 'post',
    props: {
        post: Object,
        profileView: { type: Boolean, default: false },
    },
    data() {
        return {
            date: '',
        };
    },
    methods: {
        deletePost() {
            this.$http.delete(`/api/post/${this.post.postId}`, {})
                .then(() => {
                    this.$router.go();
                })
                .catch((e) => {
                    console.log(`${JSON.stringify(e)}`);
                });
        },
    },
    created() {
        try {
            this.date = new Date(this.post.time).toLocaleString();
        } catch (e) {
            console.log(e);
        }
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.title {
  display: inline-block;
  padding: 0;
  margin: 0 $margin1 0 0;
  font-size: $fontsize2;
}

.date {
  padding: 0;
  margin: 0 $margin1 0 auto;
  font-size: $fontsize1;
}

.header {
  display: flex;
  margin-bottom: $margin1;
}

.delete {
  margin: $margin1 0 0 0;
}

.body {
  white-space: pre-line;
  padding: 0;
  margin: 0;
}
</style>
