<template>
  <div class="signup">
    <form @submit.prevent="signup">
      <input placeholder="username" v-model="username">
      <input type="email" placeholder="email" v-model="email">
      <input type="password" placeholder="password" v-model="password">
      <button type="submit">Signup</button>
    </form>
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
        };
    },
    methods: {
        signup() {
            console.log("Signup: " + this.username + ", " + this.email + ", " + this.password);
        },
    },
};
</script>

<style scoped lang="scss">
input {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
}
</style>
