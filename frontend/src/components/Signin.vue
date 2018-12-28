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
                    const expiry = (new Date(response.expiry)).getTime();
                    const current = (new Date()).getTime();
                    const daysToExpiry = Math.floor((expiry - current) / (24 * 60 * 60 * 1000));
                    this.$cookie.set('username', this.username, daysToExpiry);
                    this.$cookie.set('session', response.data.session,
                        daysToExpiry);
                    this.$emit('signin');
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
