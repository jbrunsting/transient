<template>
  <div ref="post" class="post" v-on:mousedown="startDrag">
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
    },
    data() {
        return {
            date: '',
            lastX: 0,
            lastY: 0,
            translation: 0,
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
        startDrag(e) {
            e = e || window.event;
            e.preventDefault();
            this.lastX = e.clientX;
            this.lastY = e.clientY;
            this.translation = 0;
            const post = this.$refs.post
            post.onmousemove = this.doDrag
            post.onmouseup = function() {
                post.onmousemove = undefined
                post.style.transform = "translate(0,0)"
            }
        },
        doDrag(e) {
            e = e || window.event;
            e.preventDefault();
            const dx = this.lastX - e.clientX;
            const dy = this.lastY - e.clientY;
            this.lastX = e.clientX;
            this.lastY = e.clientY;
            this.translation -= dx;
            const post = this.$refs.post
            post.style.transform = "translate(" + this.translation + "px, 0)"
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

.post {
  background-color: $base0;
  margin: $margin2;
  padding: $margin2;
  border-radius: $margin2;
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
