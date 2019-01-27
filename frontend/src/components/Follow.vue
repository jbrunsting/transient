<template>
  <div class="home">
    <form @submit.prevent="follow">
      <input placeholder="username" v-model="username">
      <button type="submit">Follow</button>
    </form>
  </div>
</template>

<script>
export default {
    name: 'home',
    data() {
        return {
            id: '',
            username: '',
            email: '',
        };
    },
    props: { authenticated: Boolean },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
    created() {
        this.$http.getProtected('/api/followings')
            .then((response) => {
                console.log(response.data);
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
    },
};
</script>
