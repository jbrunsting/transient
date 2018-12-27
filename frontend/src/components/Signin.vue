<template>
  <div class="signin">
    <form @submit.prevent="signin">
      <input placeholder="username" v-model="username">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Signin</button>
    </form>
    <h2>{{ response }}</h2>
  </div>
</template>

<script>
export default {
    name: 'Signin',
    data() {
        return {
            username: '',
            password: '',
            response: '',
        };
    },
    methods: {
        signin() {
            const user = {
                username: this.username,
                password: this.password,
            };
            this.$http.post('/api/user/login', user)
                .then((response) => {
                    this.response = JSON.stringify(response.data);
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
