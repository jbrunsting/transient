<template>
  <div class="profile">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <h1>{{ username }}</h1>
    <ul>
      <li v-for="post in posts" :key="post.postId">
        <Post :post="post"/>
      </li>
    </ul>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Post from '@/components/Post.vue';

export default {
    name: 'profile',
    data() {
        return {
            username: this.$route.params.username,
            user: {},
            posts: [],
        };
    },
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Post,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
    created() {
        this.$http.get(`/api/users/exact/${this.username}`)
            .then((response) => {
                console.log('Setting user');
                this.user = response.data;
            })
            .then(() => {
                this.$http.get(`/api/posts/${this.user.id}`)
                    .then((response) => {
                        this.posts = response.data;
                    });
            })
            .catch((e) => {
                console.log(JSON.stringify(e));
                if (e.response.status === 404) {
                    this.$router.replace('/404');
                } else {
                    // TODO: Set up proper error page for unexpected errors
                    console.log(JSON.stringify(e));
                }
            });
    },
};
</script>
