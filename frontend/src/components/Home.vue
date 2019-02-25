<template>
  <div class="home">
    <div class="postWrapper">
      <div v-if="posts[1]"
           ref="nextPost">
        <FullscreenPost :post="posts[1]"
           :translation="nextTranslation"
           :alpha="nextAlpha"
           :transition="nextTransition"/>
      </div>
      <div v-if="posts[0]"
           ref="curPost"
           v-on:mousedown="startDrag">
        <FullscreenPost
           :post="posts[0]"
           :translation="curTranslation"
           :alpha="curAlpha"
           :color="curColor"
           :transition="curTransition"/>
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
            acceptX: 200,
            curTranslation: 0,
            curAlpha: 1,
            curTransition: '',
            curColor: 'white',
            nextTranslation: 0,
            nextAlpha: 1,
            nextTransition: '',
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

            const cur = this.$refs.curPost;
            cur.onmousemove = this.doDrag;
            cur.onmouseup = this.endDrag;
            cur.onmouseleave = this.endDrag;

            this.curTranslation = 0;
            this.nextAlpha = 0;
        },
        doDrag(e) {
            e.preventDefault();

            this.curTransition = '';
            this.nextTransition = '';

            const dx = this.lastX - e.clientX;
            this.lastX = e.clientX;

            this.curTranslation -= dx;
            this.nextAlpha = Math.abs(this.curTranslation / (this.acceptX + 200));

            if (Math.abs(this.curTranslation) > this.acceptX) {
                this.nextPost();
            }
        },
        endDrag(e) {
            e.preventDefault();

            this.resetHandlers();

            this.curTransition = '100ms ease-in-out';
            this.nextTransition = '100ms ease-in-out';

            this.curTranslation = 0;
            this.nextAlpha = 0;
        },
        nextPost() {
            this.resetHandlers();

            this.curTransition = '500ms ease-out';
            this.nextTransition = '500ms ease-out';

            let target = 0;
            if (this.curTranslation > 0) {
                target = this.acceptX + 300;
                this.curColor = '#40E040';
            } else {
                target = -this.acceptX - 300;
                this.curColor = '#FF4040';
            }
            this.curTranslation = target;
            this.curAlpha = 0;
            this.nextAlpha = 1;

            setTimeout(() => {
                this.curTransition = '';
                this.nextTransition = '';
                this.posts.shift();
                this.curTranslation = 0;
                this.curAlpha = 1;
                this.curColor = '';
                this.nextAlpha = 0;
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
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  margin-top: $nav-height;
  z-index: -100;
  overflow-x: hidden;
}

.topPost {
  z-index: 100;
}

.bottomPost {
  z-index: -300;
}
</style>
