<template>
  <div class="home">
    <h1>Welcome {{ username }}</h1>
    <ul>
      <li v-for="user in followings" :key="user.id">
        <h3>{{ user.username }}</h3>
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
};
</script>
