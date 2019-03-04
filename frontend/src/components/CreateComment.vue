<template>
  <div class="createComment">
    <Error direction="none" class="unknown error">
      Could not create comment, please try again later
    </Error>
    <form @submit.prevent="createComment">
      <textarea v-model="content" placeholder="content"></textarea>
      <button type="submit">Comment</button>
    </form>
  </div>
</template>

<script>
import Error from '@/components/Error.vue';

export default {
    name: 'CreateComment',
    data() {
        return {
            content: '',
        };
    },
    props: {
        postId: String,
    },
    components: {
        Error,
    },
    methods: {
        createComment() {
            /* eslint-disable no-param-reassign */
            this.$el.querySelectorAll('.error').forEach((c) => {
                c.style.display = 'none';
            });

            const comment = {
                content: this.content,
            };

            this.$http.post('/api/post/id/comment', comment)
                .then(() => {
                    this.content = '';
                    this.$emit('createComment');
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

.createComment {
  position: relative;
}

textarea {
  resize: vertical;
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
