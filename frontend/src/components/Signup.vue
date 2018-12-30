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
                .then((response) => {
                    this.response = JSON.stringify(response.data);
                    const expiry = (new Date(response.expiry)).getTime();
                    const current = (new Date()).getTime();
                    const daysToExpiry = Math.floor((expiry - current) / (24 * 60 * 60 * 1000));
                    this.$cookie.set(this.$usernameCookie, this.username, daysToExpiry);
                    this.$cookie.set(this.$sessionIdCookie, response.data.sessionId,
                        daysToExpiry);
                    this.$emit('signup');
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
