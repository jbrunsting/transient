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
            e.preventDefault();

            this.lastX = e.clientX;
            this.lastY = e.clientY;
            this.translation = 0;

            const cur = this.$refs.curPost;
            cur.onmousemove = this.doDrag;
            cur.onmouseup = this.endDrag;
            cur.onmouseleave = this.endDrag;
            cur.style.transition = '';

            const next = this.$refs.nextPost;
            if (next) {
                next.style.opacity = 0;
            }
        },
        doDrag(e) {
            e.preventDefault();

            const dx = this.lastX - e.clientX;
            this.lastX = e.clientX;
            this.lastY = e.clientY;
            this.translation -= dx;

            const cur = this.$refs.curPost;
            cur.style.transition = '';
            cur.style.transform = `translate(${this.translation}px, 0)`;

            const next = this.$refs.nextPost;
            if (next) {
                next.style.transition = '';
                next.style.opacity = Math.abs(this.translation / (this.acceptX + 200));
            }

            if (Math.abs(this.translation) > this.acceptX) {
                this.nextPost();
            }
        },
        endDrag(e) {
            e.preventDefault();

            this.resetHandlers();

            const cur = this.$refs.curPost;
            cur.style.transition = '100ms ease-in-out';
            cur.style.transform = 'translate(0,0)';

            const next = this.$refs.nextPost;
            if (next) {
                next.style.transition = '100ms ease-in-out';
                next.style.opacity = 0;
            }
        },
        nextPost() {
            this.resetHandlers();

            const cur = this.$refs.curPost;
            cur.style.transition = '500ms ease-out';
            if (this.translation > 0) {
                cur.style.transform = `translate(${this.acceptX + 200}px,0)`;
            } else {
                cur.style.transform = `translate(-${this.acceptX + 200}px,0)`;
            }
            cur.style.opacity = 0;

            const next = this.$refs.nextPost;
            if (next) {
                next.style.transition = '500ms ease-out';
                next.style.opacity = 1;
            }

            setTimeout(() => {
                this.posts.shift();
                cur.style.transform = '';
                cur.style.transition = '';
                cur.style.opacity = 1;
            }, 500);
        },
        resetHandlers() {
            const cur = this.$refs.curPost;
            cur.onmousemove = undefined;
            cur.onmouseup = undefined;
            cur.onmouseleave = undefined;
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
