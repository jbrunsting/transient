<template>
  <div class="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </div>
</template>

<script>
export default {
    name: 'Nav',
};
</script>

<style scoped lang="scss">
.nav {
  a {
    padding: 30px;
    font-weight: bold;

    &.router-link-exact-active {
      color: red;
    }
  }
}
</style>
