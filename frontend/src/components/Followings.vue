<template>
  <div class="home">
    <ul>
      <li class="follow" v-for="user in followings" :key="user.id">
        <p>{{ user.username }}</p>
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

<style scoped lang="scss">
@import "../styles/settings.scss";

.follow {
  display: flex;
}

.follow p {
  display: inline-block;
  margin: $margin0;
}

.follow form {
  display: inline-block;
  margin: $margin0;
  margin-left: auto;
}
</style>
