<template>
  <div class="about">
    <h1>Welcome {{ username }}</h1>
    <p>{{ email }}</p>
    <CreatePost v-on:createPost="getPosts"/>
    <ul>
      <li v-for="post in posts" :key="post.postId">
        <ul>
          <a v-if="post.postUrl" :href="post.postUrl"><h3>{{ post.title }}</h3></a>
          <h3 v-else>{{ post.title }}</h3>
          <p>{{ post.content }}</p>
        </ul>
      </li>
    </ul>
  </div>
</template>

<script>
import CreatePost from '@/components/CreatePost.vue';

export default {
    name: 'about',
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
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>
