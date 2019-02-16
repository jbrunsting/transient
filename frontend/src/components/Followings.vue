<template>
  <div class="home">
    <h1>Welcome {{ username }}</h1>
    <ul>
      <li v-for="user in followings" :key="user.id">
        <h3>{{ user.username }}</h3>
        <form @submit.prevent="() => unfollow(user.id)">
          <button type="submit">Unfollow</button>
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
            id: '',
            username: '',
            email: '',
        };
    },
    props: { authenticated: Boolean, followings: Array },
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
    },
};
</script>
