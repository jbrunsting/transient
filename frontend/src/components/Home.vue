<template>
  <div class="home">
    <h1>Welcome to your home feed {{ username }}</h1>
    <ul>
      <li v-for="post in posts" :key="post.postId">
        <p>{{ post.username }}</p>
        <Post :post="post"/>
      </li>
    </ul>
  </div>
</template>

<script>
import Post from '@/components/Post.vue';

export default {
    name: 'home',
    data() {
        return {
            id: '',
            username: '',
            email: '',
            posts: [],
        };
    },
    props: { authenticated: Boolean },
    components: {
        Post,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
        getPosts() {
            this.$http.get('/api/followings/posts')
                .then((response) => {
                    console.log(this.posts);
                    this.posts = response.data;
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
