<template>
  <div class="followingRecommends">
    <ul>
      <li class="follow" v-for="user in recommends" :key="user.id">
        <p>{{ user.username }}</p>
        <form @submit.prevent="() => follow(user.id)">
          <button type="submit">Follow</button>
        </form>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
    name: 'followingRecommends',
    data() {
        return {
            recommends: [],
        };
    },
    props: { authenticated: Boolean, followings: Array },
    methods: {
        getRecommends() {
            this.$http.get('/api/recommends/followings')
                .then((response) => {
                    console.log(response);
                    this.recommends = response.data;
                }).catch((e) => {
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
    },
    created() {
        this.getRecommends();
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
