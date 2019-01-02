<template>
  <div class="login">
    <form @submit.prevent="login">
      <input placeholder="username" v-model="username">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script>
export default {
    name: 'Login',
    data() {
        return {
            username: '',
            password: '',
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
@import "../styles/settings.scss";

input {
  margin-top: 0;
  margin-bottom: 0;
  width: 120px;
}

button {
  margin-left: $margin0;
}

form {
  display: flex;
  align-items: center;
}
</style>
