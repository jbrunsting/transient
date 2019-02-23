<template>
  <div ref="post" class="post">
    <div class="header">
      <h3 class="title" v-if="post.postUrl" >
        <a :href="post.postUrl">{{ post.title }}</a>
      </h3>
      <h3 class="title" v-else>{{ post.title }}</h3>
      <p class="date">{{ date }}</p>
      <p class="username">
        <a :href="'/profile/' + post.username">{{ post.username }}</a>
      </p>
    </div>
    <p class="body">{{ post.content }}</p>
  </div>
</template>

<script>
export default {
    name: 'post',
    props: {
        post: Object,
        translation: Number,
        alpha: Number,
        color: String,
        transition: String,
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
    watch: {
        translation() {
            this.$refs.post.style.transform = `translate(${this.translation}px, 0)`;
        },
        alpha() {
            this.$refs.post.style.opacity = this.alpha;
        },
        color() {
            this.$refs.post.style.backgroundColor = this.color;
        },
        transition() {
            this.$refs.post.style.transition = this.transition;
        },
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.post {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  background-color: $base0;
  margin: $margin2 auto;
  padding: $margin2;
  border-radius: $margin2;
  max-width: 600px;
}

.username {
  float: right;
  padding: 0;
  margin: 0;
  margin: 0 0 0 $margin1;
  font-size: $fontsize1;
}

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

.body {
  padding: 0;
  margin: $margin1 0 0 0;
}
</style>
