<template>
  <div class="home">
    <form @submit.prevent="search">
      <input placeholder="username" v-model="username">
      <button type="submit">Search</button>
    </form>
    <ul>
      <li v-for="user in users" :key="user.id">
        <h3>{{ user.username }}</h3>
        <form @submit.prevent="follow">
          <button type="submit">Search</button>
        </form>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
    name: 'home',
    data() {
        return {
            username: '',
            users: [],
        };
    },
    props: { authenticated: Boolean },
    methods: {
        search() {
            this.$http.get(`/api/users/search?username=${this.username}`)
                .then((response) => {
                    this.users = response.data;
                    console.log(this.users);
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
    },
};
</script>
