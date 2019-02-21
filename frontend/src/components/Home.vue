<template>
  <div class="home">
    <div class="postWrapper">
      <div v-if="posts[1]" ref="nextPost">
        <FullscreenPost :post="posts[1]"/>
      </div>
      <div v-if="posts[0]" ref="curPost" v-on:mousedown="startDrag">
        <FullscreenPost :post="posts[0]"/>
      </div>
    </div>
  </div>
</template>

<script>
import FullscreenPost from '@/components/FullscreenPost.vue';

export default {
    name: 'home',
    data() {
        return {
            id: '',
            username: '',
            email: '',
            posts: [],
            lastX: 0,
            lastY: 0,
            translation: 0,
            acceptX: 200,
        };
    },
    props: { authenticated: Boolean },
    components: {
        FullscreenPost,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
        getPosts() {
            this.$http.get('/api/followings/posts')
                .then((response) => {
                    this.posts = response.data;
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
        startDrag(e) {
            e = e || window.event;
            e.preventDefault();
            this.lastX = e.clientX;
            this.lastY = e.clientY;
            this.translation = 0;

            const cur = this.$refs.curPost;
            cur.onmousemove = this.doDrag;
            cur.onmouseup = this.endDrag;
            cur.onmouseleave = this.endDrag;
            cur.style.transition = "";

            const next = this.$refs.nextPost;
            next.style.opacity = 0;
        },
        doDrag(e) {
            e = e || window.event;
            e.preventDefault();
            const dx = this.lastX - e.clientX;
            const dy = this.lastY - e.clientY;
            this.lastX = e.clientX;
            this.lastY = e.clientY;
            this.translation -= dx;

            const cur = this.$refs.curPost;
            cur.style.transform = "translate(" + this.translation + "px, 0)";

            const next = this.$refs.nextPost;
            next.style.opacity = Math.abs(this.translation / this.acceptX);
        },
        endDrag(e) {
            e.preventDefault();
            const post = this.$refs.curPost;
            post.onmousemove = undefined;
            post.style.transform = "translate(0,0)";
            post.style.transition = "100ms ease-in-out";
        },
    },
    created() {
        this.$http.getProtected('/api/user')
            .then((response) => {
                this.id = response.data.id;
                this.username = response.data.username;
                this.email = response.data.email;
                this.getPosts();
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.home {
  overflow: hidden;
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  margin-top: $nav-height;
  z-index: -100;
}

.topPost {
  z-index: 100;
}

.bottomPost {
  z-index: -300;
}
</style>
