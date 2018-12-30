<template>
  <div class="signout">
    <form @submit.prevent="signout">
      <button type="submit">Signout</button>
    </form>
  </div>
</template>

<script>
export default {
    name: 'Signout',
    methods: {
        signout() {
            this.$cookie.delete(this.$usernameCookie);
            this.$cookie.delete(this.$sessionIdCookie);
            this.$emit('signout');
        },
    },
};
</script>
