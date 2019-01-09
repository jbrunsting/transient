<template>
  <div class="createPost">
    <Error direction="none" class="unknown error">
      Could not create post, please try again later
    </Error>
    <form @submit.prevent="createPost">
      <input placeholder="title" v-model="title">
      <input type="url" placeholder="url" v-model="postUrl">
      <textarea v-model="content" placeholder="content"></textarea>
      <button type="submit">CreatePost</button>
    </form>
  </div>
</template>

<script>
import Error from '@/components/Error.vue';

export default {
    name: 'CreatePost',
    data() {
        return {
            title: '',
            content: '',
            postUrl: '',
        };
    },
    components: {
        Error,
    },
    methods: {
        createPost() {
            /* eslint-disable no-param-reassign */
            this.$el.querySelectorAll('.error').forEach((c) => {
                c.style.display = 'none';
            });

            const post = {
                title: this.title,
                content: this.content,
                postUrl: this.postUrl,
            };

            this.$http.post('/api/post', post)
                .then(() => {
                    this.$emit('createPost');
                }).catch((e) => {
                    console.log(`Error: ${e.response.data.message}`);
                    this.$el.querySelector('.unknown.error').style.display = 'inline-block';
                });
            /* eslint-enable no-param-reassign */
        },
    },
};
</script>

<style scoped lang="scss">
@import "../styles/settings.scss";

.createPost {
  position: relative;
}

input,
textarea {
  margin-right: 0;
  margin-left: 0;
  width: auto;
}

form {
  display: flex;
  flex-direction: column;
}

.error {
  display: none;
}
</style>
