<template>
  <div class="post">
    <a v-if="post.postUrl" :href="post.postUrl"><h3>{{ post.title }}</h3></a>
    <h3 v-else>{{ post.title }}</h3>
    <p>{{ post.content }}</p>
    <p>{{ new Date(post.time).toLocaleString() }}</p>
    <form v-if="showControls" @submit.prevent="deletePost">
      <button type="submit">Delete</button>
    </form>
  </div>
</template>

<script>
export default {
    name: 'post',
    props: {
        post: Object,
        showControls: { type: Boolean, default: false },
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
};
</script>
