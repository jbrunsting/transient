<template>
  <div class="login">
    <form @submit.prevent="login">
      <input placeholder="username" v-model="username">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Login</button>
    </form>
    <h2>{{ response }}</h2>
  </div>
</template>

<script>
export default {
    name: 'Login',
    data() {
        return {
            username: '',
            password: '',
            response: '',
        };
    },
    methods: {
        login() {
            const identification = {
                username: this.username,
                password: this.password,
            };
            this.$http.post('/api/user/login', identification)
                .then(() => {
                    this.$emit('login');
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
    },
};
</script>

<style scoped lang="scss">
input {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
}
</style>
