<template>
  <div class="home">
    <FullscreenPost v-if="post" :post="post"/>
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
            post: undefined,
            posts: [],
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
                    this.post = this.posts[0];
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
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
  margin: $nav-height auto 0 auto;
  max-width: $page-width;
  z-index: -100;
}
</style>
