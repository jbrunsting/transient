<template>
  <div class="homeprofile">
    <CreatePost v-on:createPost="getPosts"/>
    <ul>
      <li v-for="post in posts" :key="post.postId">
        <Post :post="post" :profileView="true"/>
      </li>
    </ul>
  </div>
</template>

<script>
import CreatePost from '@/components/CreatePost.vue';
import Post from '@/components/Post.vue';

export default {
    name: 'homeprofile',
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
        CreatePost,
        Post,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
        getPosts() {
            this.$http.get(`/api/posts/${this.id}`)
                .then((response) => {
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
                console.log(e);
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>
