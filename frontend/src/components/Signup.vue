<template>
  <div class="signup">
    <form @submit.prevent="signup">
      <input placeholder="username" v-model="username">
      <input type="email" placeholder="email" v-model="email">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Signup</button>
    </form>
    <h2>{{ response }}</h2>
  </div>
</template>

<script>
export default {
    name: 'Signup',
    data() {
        return {
            username: '',
            email: '',
            password: '',
            response: '',
        };
    },
    methods: {
        signup() {
            const user = {
                username: this.username,
                email: this.email,
                password: this.password,
            };
            this.$http.post('/api/user', user)
                .then(() => {
                    this.$emit('signup');
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
  margin-right: 0;
  margin-left: 0;
}

form {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}
</style>
