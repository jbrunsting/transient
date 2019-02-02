<template>
  <div class="home">
    <h1>Welcome {{ username }}</h1>
    <ul>
      <li v-for="id in followings" :key="id">
        <h3>{{ id }}</h3>
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
