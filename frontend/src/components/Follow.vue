<template>
  <div class="home">
    <form @submit.prevent="search">
      <input placeholder="username" v-model="username">
      <button type="submit">Search</button>
    </form>
    <ul>
      <li v-for="user in users" :key="user.id">
        <h3>{{ user.username }}</h3>
        <template v-if="followings.includes(user.id)">
          <form @submit.prevent="() => unfollow(user.id)">
            <button type="submit">Unfollow</button>
          </form>
        </template>
        <template v-else>
          <form @submit.prevent="() => follow(user.id)">
            <button type="submit">Follow</button>
          </form>
        </template>
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
    props: {
        followings: Array,
        authenticated: Boolean,
    },
    methods: {
        unfollow(id) {
            this.$http.delete(`/api/following/${id}`)
                .then(() => {
                    this.$emit('follow');
                })
                .catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
        follow(id) {
            this.$http.post(`/api/following/${id}`)
                .then(() => {
                    this.$emit('follow');
                })
                .catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
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
